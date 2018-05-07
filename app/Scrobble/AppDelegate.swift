//
//  AppDelegate.swift
//  Scrobble
//
//  Created by Nishanth Shanmugham on 5/5/18.
//  Copyright © 2018 Nishanth Shanmugham. All rights reserved.
//

// Some parts of this file are derived from
// https://www.raywenderlich.com/165853/menus-popovers-menu-bar-apps-macos.
// See license below.

/**
 * Copyright (c) 2017 Razeware LLC
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * Notwithstanding the foregoing, you may not use, copy, modify, merge, publish,
 * distribute, sublicense, create a derivative work, and/or sell copies of the
 * Software in any work that is designed, intended, or marketed for pedagogical or
 * instructional purposes related to programming, coding, application development,
 * or information technology.  Permission for such use, copying, modification,
 * merger, publication, distribution, sublicensing, creation of derivative works,
 * or sale is expressly withheld.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

import Cocoa

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate, LoginViewControllerDelegate,
    ScrobbleViewControllerDelegate, WatcherDelegate {

    static let kAPIToken = "API Token"
    static let baseUrl = "https://scrobble.allele.cc"

    let statusItem = NSStatusBar.system.statusItem(withLength:NSStatusItem.squareLength)
    let popover = NSPopover()
    let watcher = Watcher()
    var eventMonitor: EventMonitor?
    var reset = false

    lazy var loginViewController: LoginViewController = {
        let v = LoginViewController.create()
        v.delegate = self
        return v
    }()

    lazy var scrobbleViewController: ScrobbleViewController = {
        let v = ScrobbleViewController.create()
        v.delegate = self
        return v
    }()

    func applicationDidFinishLaunching(_ aNotification: Notification) {
        if let button = statusItem.button {
            button.image = NSImage(named: NSImage.Name(rawValue: "icon"))
            button.image?.isTemplate = true
            button.action = #selector(AppDelegate.togglePopover(_:))
        }

        popover.behavior = .transient
        popover.animates = true // gross hack: if animation is disabled, layout messes up b/w popover views
        popover.appearance = NSAppearance.init(named: NSAppearance.Name.vibrantLight)
        initializePopoverController()

        eventMonitor = EventMonitor(mask: [.leftMouseDown, .rightMouseDown]) { [weak self] event in
            if let strongSelf = self, strongSelf.popover.isShown {
                strongSelf.closePopover(sender: event)
            }
        }
    }

    func initializePopoverController() {
        if let token = UserDefaults.standard.string(forKey: AppDelegate.kAPIToken) {
            validateToken(token, completion: {account, status in
                guard let a = account else {
                    // validation failed
                    UserDefaults.standard.removeObject(forKey: AppDelegate.kAPIToken)
                    DispatchQueue.main.async {
                        self.loginViewController.invalid(status == .generalError ?
                            LoginViewController.tryAgainGeneral : LoginViewController.tryAgainClient)
                        self.loginMode()
                    }
                    return
                }
                // validation successful
                UserDefaults.standard.set(token, forKey: AppDelegate.kAPIToken)
                DispatchQueue.main.async {
                    self.scrobblingMode(token, a)
                }
            })
        }

        // no saved token or failed to validate
        UserDefaults.standard.removeObject(forKey: AppDelegate.kAPIToken)
        DispatchQueue.main.async {
            self.loginMode()
        }
    }

    func validateToken(_ t: String, completion: @escaping (Account?, ResponseStatus) -> Void) {
        let r = API.accountRequest(t)

        let task = URLSession.shared.dataTask(with: r, completionHandler: {(data, rsp, err) in
            if let h = rsp as? HTTPURLResponse {
                switch h.statusCode/100 {
                case 2:
                    // data is guaranteed to exist according to the docs
                    let account = try? JSONDecoder().decode(Account.self, from: data!)
                    completion(account, .success)
                case 3:
                    completion(nil, .authError)
                case 4:
                    completion(nil, .clientError)
                default:
                    completion(nil, .generalError)
                }
            }
            completion(nil, .generalError)
        })

        task.resume()
    }

    func loginMode() {
        watcher.stop()
        popover.contentViewController = loginViewController
    }

    func scrobblingMode(_ token: String, _ account: Account) {
        watcher.start(token)
        scrobbleViewController.username = account.username
        popover.contentViewController = scrobbleViewController
    }

    @objc private func togglePopover(_ sender: Any?) {
        if popover.isShown {
            closePopover(sender: sender)
        } else {
            showPopover(sender: sender)
        }
    }

    private func showPopover(sender: Any?) {
        // TODO: this is a gross hack
        if reset {
            if let v = popover.contentViewController as? LoginViewController {
                v.regular()
            }
        }

        if let button = statusItem.button {
            popover.show(relativeTo: button.bounds, of: button, preferredEdge: NSRectEdge.minY)
            eventMonitor?.start()
        }
    }

    private func closePopover(sender: Any?) {
        popover.performClose(sender)
        eventMonitor?.stop()
        reset = true
    }

    // MARK delegate functions for ScrobbleViewController

    func onAuthError() {
        UserDefaults.standard.removeObject(forKey: AppDelegate.kAPIToken)
        DispatchQueue.main.async {
            self.loginViewController.regular()
            self.loginMode()
        }
    }

    func onLogout() {
        UserDefaults.standard.removeObject(forKey: AppDelegate.kAPIToken)
        DispatchQueue.main.async {
            self.loginViewController.regular()
            self.loginMode()
        }
    }

    // MARK delegate functions for LoginViewController

    func onSubmitToken() {
        loginViewController.validating()
        let token = loginViewController.tokenInput.stringValue

        validateToken(token, completion: {account, status in
            // TODO: consider sharing this with the similar block above
            guard let a = account else {
                // validation failed
                UserDefaults.standard.removeObject(forKey: AppDelegate.kAPIToken)
                DispatchQueue.main.async {
                    self.loginViewController.invalid(status == .generalError ?
                        LoginViewController.tryAgainGeneral : LoginViewController.tryAgainClient)
                }
                return
            }
            // validation successful
            UserDefaults.standard.set(token, forKey: AppDelegate.kAPIToken)
            DispatchQueue.main.async {
                self.scrobblingMode(token, a)
            }
        })
    }
}

class EventMonitor {
    private var monitor: Any?
    private let mask: NSEvent.EventTypeMask
    private let handler: (NSEvent?) -> Void

    public init(mask: NSEvent.EventTypeMask, handler: @escaping (NSEvent?) -> Void) {
        self.mask = mask
        self.handler = handler
    }

    deinit {
        stop()
    }

    public func start() {
        monitor = NSEvent.addGlobalMonitorForEvents(matching: mask, handler: handler)
    }

    public func stop() {
        if monitor != nil {
            NSEvent.removeMonitor(monitor!)
            monitor = nil
        }
    }
}

