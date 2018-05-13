//
//  ScrobbleViewController.swift
//  Scrobble
//
//  Created by Nishanth Shanmugham on 5/5/18.
//  Copyright © 2018 Nishanth Shanmugham. All rights reserved.
//

import Cocoa
import Foundation

protocol LoginViewControllerDelegate {
    func onSubmitToken() -> Void
}

class LoginViewController: NSViewController {
    @IBOutlet weak var tokenInput: NSTextField!
    @IBOutlet weak var introText: NSTextField!
    @IBOutlet weak var startButton: NSButton!

    static let enterToken = "Enter API token to\nstart scrobbling";
    static let loggingIn = "\nSigning in..."
    static let tryAgainClient = "Invalid token.\nTry again?"
    static let tryAgainGeneral = "Something went wrong.\nTry again?"

    var delegate: LoginViewControllerDelegate?
    var enabled = true
    var introString = LoginViewController.enterToken

    func validating() {
        enabled = false
        introString = LoginViewController.loggingIn
        view.needsDisplay = true
    }

    func invalid(_ message: String) {
        enabled = true
        introString = message
        view.needsDisplay = true
    }

    func regular() {
        enabled = true
        introString = LoginViewController.enterToken
        view.needsDisplay = true
    }

    override func viewWillAppear() {
        redraw()
    }

    func redraw() {
        tokenInput.isEnabled = enabled
        startButton.isEnabled = enabled
        introText.stringValue = introString // TODO: do a Core Animation transition?
        tokenInput.becomeFirstResponder()
        tokenInput.selectText(nil)
    }

    @IBAction func docsClick(_ sender: Any) {
        NSWorkspace.shared.open(URL(string: AppDelegate.baseUrl)!)
    }

    @objc func onEnter() {
        delegate?.onSubmitToken()
    }

    @IBAction func start(_ sender: NSButton) {
        delegate?.onSubmitToken()
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        tokenInput.action = #selector(LoginViewController.onEnter)
        tokenInput.becomeFirstResponder()
    }

    @IBAction func quit(sender: NSButton) {
        NSApplication.shared.terminate(sender)
    }

    static func create() -> LoginViewController {
        let storyboard = NSStoryboard(name: NSStoryboard.Name(rawValue: "Main"), bundle: nil)
        let identifier = NSStoryboard.SceneIdentifier(rawValue: "LoginViewController")
        guard let viewcontroller = storyboard.instantiateController(withIdentifier: identifier) as? LoginViewController else {
            fatalError("failed to find LoginViewController")
        }
        return viewcontroller
    }
}

protocol ScrobbleViewControllerDelegate {
    func onLogout() -> Void
}

class ScrobbleViewController: NSViewController {
    @IBOutlet weak var usernameField: NSButton!
    @IBOutlet weak var imageField: NSImageView!
    var delegate: ScrobbleViewControllerDelegate?

    lazy var attrs: [NSAttributedStringKey : Any]? = {
        let pstyle = NSMutableParagraphStyle()
        pstyle.alignment = .center
        return [
            NSAttributedStringKey.foregroundColor: NSColor.black,
            NSAttributedStringKey.paragraphStyle: pstyle,
            NSAttributedStringKey.font: NSFont.init(name: "Menlo", size: 14)!
        ]
    }()

    var username: String? {
        didSet {
            view.needsDisplay = true
        }
    }

    static func create() -> ScrobbleViewController {
        let storyboard = NSStoryboard(name: NSStoryboard.Name(rawValue: "Main"), bundle: nil)
        let identifier = NSStoryboard.SceneIdentifier(rawValue: "ScrobbleViewController")
        guard let viewcontroller = storyboard.instantiateController(withIdentifier: identifier) as? ScrobbleViewController else {
            fatalError("failed to find SrobbleViewController")
        }
        return viewcontroller
    }

    @IBAction func usernameClick(_ sender: Any) {
        NSWorkspace.shared.open(URL(string: userUrl())!)
    }

    @IBAction func logout(_ sender: Any) {
        delegate?.onLogout()
    }

    override func viewWillAppear() {
       redraw()
    }

    private func redraw() {
        usernameField.attributedTitle = NSMutableAttributedString(string: (username ?? ""), attributes: attrs)
        usernameField.toolTip = userUrl()
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        imageField.image = NSImage(named: NSImage.Name(rawValue: "scrobbling"))
    }

    private func userUrl() -> String {
        return AppDelegate.baseUrl + "/" + (username ?? "")
    }

    @IBAction func quit(sender: NSButton) {
        NSApplication.shared.terminate(sender)
    }
}

