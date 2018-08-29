//
//  AppDelegate.swift
//  itunes-scrobble
//
//  Created by Nishanth Shanmugham on 8/26/18.
//  Copyright © 2018 Nishanth Shanmugham. All rights reserved.
//

import Cocoa
import iTunesLibrary

// Notes
// -----
//
// Using URLSession fails with "HTTP load failed (error code: 100 ...",
// presumably because the app is unsigned? No amount of Info.plist and
// itunes_scrobble.entitlements hacking fixes it. Hence the use of
// the deprecated NSURLConnection.sendAsynchronousRequest.

// State is the state of the application.
struct State {
    // whether scrobbling is running or paused
    var running: Bool
    // the user entered API key
    var apiKey: String?
    // last successful scrobble by this instance
    var lastScrobbled: Date?
    // latest "Last Played" time sent by this instance
    // TODO: this is not surfaced to the user
    var latestPlayed: Date?
    // the accounts response from the server
    var account: API.Account?
    // whether there is a scrobble request inflight
    var scrobbling: Bool
    // whether the latest scrobble request resulted in an auth error
    var authError: Bool
}

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate, NSTextFieldDelegate, NSAlertDelegate {
    private static let menuIconName = "itunes-scrobble-18x18" // size from https://stackoverflow.com/a/33708433
    static let baseUrl = "selective-scrobble.appspot.com"
    private static let helpLink = "https://scrobble.allele.cc"
    private static let shortVersion = Bundle.main.object(forInfoDictionaryKey: "CFBundleShortVersionString") as! String
    
    // Keys for information saved to UserDefaults.
    private static let (keyAPIKey, keyRunning, keyLastScrobbled, keyLatestPlayed) =
        ("apiKey", "running", "lastScrobbled", "latestPlayed")
    
    // Scrobble timer frequnecies.
    // The timer fires frequently, but scrobbling happens less often.
    private static let timerFreq: TimeInterval = 60 * 10;
    private static let scrobbleFreq: TimeInterval = 60 * 60 * 1;
    private var timer: Timer? = nil
    
    private let statusBarItem = NSStatusBar.system.statusItem(withLength:NSStatusItem.squareLength)
    private let pauseItem = NSMenuItem(title: "", action: nil, keyEquivalent: "")
    private let statusItem = NSMenuItem(title: "", action: nil, keyEquivalent: "")
    private let multiItem = NSMenuItem(title: "", action: nil, keyEquivalent: "")
    
    private var prevState: State? = nil
    private var state = State(running: false,
                              apiKey: nil,
                              lastScrobbled: nil,
                              latestPlayed: nil,
                              account: nil,
                              scrobbling: false,
                              authError: false)
    
    private var lib: ITLibrary? = nil
    
    func applicationDidFinishLaunching(_ aNotification: Notification) {
        if let l = try? ITLibrary.init(apiVersion: "1.0") {
            lib = l
        } else {
            os_log("failed to initialize ITLibrary")
            return
        }
        
        // make status bar item and menu
        let button = statusBarItem.button!
        button.image = NSImage(named:NSImage.Name(AppDelegate.menuIconName))
        button.image!.isTemplate = true
        statusBarItem.menu = makeMenu()
        
        // restore persisted information
        let ls = UserDefaults.standard.double(forKey: AppDelegate.keyLastScrobbled)
        let lp = UserDefaults.standard.double(forKey: AppDelegate.keyLatestPlayed)
        state.running = UserDefaults.standard.bool(forKey: AppDelegate.keyRunning)
        state.apiKey = UserDefaults.standard.string(forKey: AppDelegate.keyAPIKey)
        state.lastScrobbled = ls != 0 ? Date(timeIntervalSince1970: TimeInterval(ls)) : nil
        state.latestPlayed = lp != 0 ? Date(timeIntervalSince1970: TimeInterval(lp)) : nil
        render()
        
        // initially fetch account info
        guard let key = state.apiKey else { return }
        NSURLConnection.sendAsynchronousRequest(API.accountRequest(key), queue: OperationQueue.main) {(rsp, data, err) in
            guard err == nil else { return }
            guard let rr = rsp as! HTTPURLResponse? else { return }
            guard rr.statusCode == 200 else { return }
            DispatchQueue.main.async {
                let account = try? JSONDecoder().decode(API.Account.self, from: data!)
                self.state.account = account
                self.render()
            }
        }
    }
    
    func applicationWillTerminate(_ aNotification: Notification) {
    }
    
    private func makeMenu() -> NSMenu {
        let m = NSMenu()
        
        let v = NSMenuItem(title: String(format: "itunes-scrobble v%@", AppDelegate.shortVersion), action: nil, keyEquivalent: "")
        v.isEnabled = false
        
        m.addItem(pauseItem)
        m.addItem(NSMenuItem.separator())
        m.addItem(multiItem)
        m.addItem(statusItem)
        m.addItem(NSMenuItem.separator())
        m.addItem(v)
        m.addItem(NSMenuItem(title: "Quit", action: #selector(NSApplication.terminate(_:)), keyEquivalent: ""))
        
        return m
    }
    
    private func render() {
        defer { prevState = state }
        
        if state.apiKey == nil {
            assert(!state.running)
            multiItem.title = "Start scrobbling..."
            multiItem.action =  #selector(enterAPIKeyAction(_:))
            pauseItem.isHidden = true
        } else {
            if let a = state.account {
                multiItem.title = String(format: "Signed in as %@...", a.username)
                multiItem.action = #selector(scrobblingAsAction(_:))
            } else {
                multiItem.title = String(format: "Clear API Key")
                multiItem.action = #selector(clearAPIKeyAction(_:))
            }
            if (state.running) {
                pauseItem.title = "Pause scrobbling"
                pauseItem.action = #selector(pauseAction(_:))
            } else {
                pauseItem.title = "Continue scrobbling"
                pauseItem.action = #selector(startAction(_:))
            }
            pauseItem.isHidden = false
        }
        
        // Status item
        if state.authError {
            statusItem.title = "Failed to scrobble: outdated API Key?"
            statusItem.isHidden = false
        } else if state.scrobbling {
            statusItem.title = String(format: "Scrobbling now...")
            statusItem.isHidden = false
        } else {
            if let ls = state.lastScrobbled {
                statusItem.title = String(format: "Last scrobbled: %@", formatDate(ls))
                statusItem.isHidden = false
            } else {
                statusItem.isHidden = true
            }
        }
        
        // Timers
        let alreadyRunning = prevState?.running ?? false
        if alreadyRunning && !state.running {
            timer!.invalidate()
        } else if !alreadyRunning && state.running {
            DispatchQueue.global(qos: .userInitiated).async { self.scrobble() } // initial
            timer = Timer(timeInterval: AppDelegate.timerFreq, target: self, selector: #selector(timerFired(_:)), userInfo: nil, repeats: true)
            RunLoop.current.add(timer!, forMode: .commonModes)
        }
    }
    
    private func scrobble() {
        if state.scrobbling {
            // already scrobbling
            return
        }
        
        state.scrobbling = true
        render()
        
        lib!.reloadData()
        let (items, latest) = scrobblableItems(from: lib!.allMediaItems)
        guard let data = try? JSONEncoder().encode(items) else { return }
        
        NSURLConnection.sendAsynchronousRequest(API.scrobbleRequest(state.apiKey!, data), queue: OperationQueue.main) {(rsp, data, err) in
            defer {
                self.state.scrobbling = false
                self.render()
            }
            
            guard err == nil else { return }
            guard let r = rsp as! HTTPURLResponse? else { return }
            
            if r.statusCode == 401 {
                self.state.authError = true
                self.render()
                return
            }
            
            if r.statusCode == 200 {
                DispatchQueue.main.async {
                    self.state.lastScrobbled = Date(timeIntervalSinceNow: 0)
                    self.state.latestPlayed = latest
                    self.state.authError = false
                    UserDefaults.standard.set(self.state.lastScrobbled?.timeIntervalSince1970, forKey: AppDelegate.keyLastScrobbled)
                    UserDefaults.standard.set(self.state.latestPlayed?.timeIntervalSince1970, forKey: AppDelegate.keyLatestPlayed)
                    self.render()
                }
                self.handleArtwork()
                return
            }
        }
    }
    
    private func handleArtwork() {
        guard let key = state.apiKey else { return }
        NSURLConnection.sendAsynchronousRequest(API.missingArtworkRequest(key), queue: OperationQueue.main) {(rsp, data, err) in
            if err != nil {
                return
            }
            guard let r = rsp as! HTTPURLResponse? else {
                return
            }
            if r.statusCode != 200 {
                return
            }
            guard let d = data else {
                return
            }
            
            guard let incomingHashes = try? JSONDecoder().decode(Dictionary<String, Bool>.self, from: d) else { return }
            
            DispatchQueue.global(qos: .userInitiated).async {
                var items = Dictionary<String, ITLibMediaItem>()
                for p in self.lib!.allMediaItems {
                    guard let hash = API.MediaItem(fromITLibMediaItem: p).artworkHash else { continue }
                    items[hash] = p
                }
                
                // send artwork that the server asks for
                for (h, _) in incomingHashes {
                    guard self.state.running else { return }
                    guard let key = self.state.apiKey else { return }
                    guard let p = items[h] else { continue }
                    guard let d = p.artwork?.imageData else { continue }
                    guard let f = p.artwork?.imageDataFormat else { continue }
                    NSURLConnection.sendAsynchronousRequest(API.artworkRequest(key, f, d), queue: OperationQueue.main) {(_, _, _) in
                    }
                }
            }
        }
    }
    
    private func clearAPIKey() {
        state.running = false
        state.apiKey = nil
        state.lastScrobbled = nil
        state.latestPlayed = nil
        state.account = nil
        UserDefaults.standard.set(state.running, forKey: AppDelegate.keyRunning)
        UserDefaults.standard.set(state.apiKey, forKey: AppDelegate.keyAPIKey)
        UserDefaults.standard.set(state.lastScrobbled, forKey: AppDelegate.keyLastScrobbled)
        UserDefaults.standard.set(state.latestPlayed, forKey: AppDelegate.keyLatestPlayed)
        render()
    }
    
    @objc private func pauseAction(_ sender: Any?) {
        state.running = false
        UserDefaults.standard.set(state.running, forKey: AppDelegate.keyRunning)
        render()
    }
    
    @objc private func startAction(_ sender: Any?) {
        state.running = true
        UserDefaults.standard.set(state.running, forKey: AppDelegate.keyRunning)
        render()
    }
    
    @objc private func timerFired(_ sender: Any?) {
        if state.lastScrobbled != nil && abs(state.lastScrobbled!.timeIntervalSinceNow) < AppDelegate.scrobbleFreq {
            // already scrobbled in near past
            return
        }
        scrobble()
    }
    
    @objc private func clearAPIKeyAction(_ sender: Any?) {
        assert(state.apiKey != nil && state.account == nil)
        clearAPIKey()
    }
    
    @objc private func scrobblingAsAction(_ sender: Any?) {
        assert(state.apiKey != nil && state.account != nil)
        let a = NSAlert()
        a.alertStyle = .informational
        a.messageText = String(format: "Scrobbling as %@.", state.account!.username)
        a.showsSuppressionButton = false
        a.showsHelp = true
        a.delegate = self
        a.addButton(withTitle: "Clear API Key")
        a.addButton(withTitle: "Cancel")
        
        let result = a.runModal()
        switch result {
        case NSApplication.ModalResponse.alertFirstButtonReturn:
            clearAPIKey()
        case NSApplication.ModalResponse.alertSecondButtonReturn:
        break // nothing to do
        default:
            print("unhandled button", result)
        }
    }
    
    private var alert: NSAlert? = nil
    private var textField: NSTextField? = nil
    private var oldOkButtonTarget: AnyObject? = nil
    private var oldOkButtonAction: Selector? = nil
    
    @objc private func enterAPIKeyAction(_ sender: Any?) {
        alert = NSAlert()
        alert!.alertStyle = .informational
        alert!.messageText = "Enter API key."
        alert!.showsSuppressionButton = false
        alert!.showsHelp = true
        alert!.delegate = self
        
        let okButton = alert!.addButton(withTitle: "OK")
        oldOkButtonTarget = okButton.target
        oldOkButtonAction = okButton.action
        okButton.target = self
        okButton.action = #selector(okButtonAction(_:))
        alert!.addButton(withTitle: "Cancel")
        
        textField = NSTextField(frame: NSMakeRect(0, 0, 250, NSFont.systemFontSize * 1.8))
        textField!.usesSingleLineMode = true
        textField!.cell?.wraps = false
        textField!.cell?.isScrollable = false
        textField!.delegate = self
        textField!.font = NSFont.monospacedDigitSystemFont(ofSize: NSFont.systemFontSize, weight: .regular)
        textField!.placeholderString = "e.g., D1A3903GB"
        alert!.accessoryView = textField
        alert!.window.initialFirstResponder = alert!.accessoryView
        
        alert!.runModal()
    }
    
    func alertShowHelp(_ alert: NSAlert) -> Bool {
        if let u = URL(string: AppDelegate.helpLink) {
            NSWorkspace.shared.open(u)
        }
        return true
    }
    
    @objc private func okButtonAction(_ sender: Any?) {
        let key = textField!.stringValue
        if key.isEmpty {
            return
        }
        
        NSURLConnection.sendAsynchronousRequest(API.accountRequest(key), queue: OperationQueue.main) {(rsp, data, err) in
            if err == nil {
                if let rr = rsp as! HTTPURLResponse? {
                    if rr.statusCode == 200 {
                        // All good.
                        DispatchQueue.main.async {
                            let account = try? JSONDecoder().decode(API.Account.self, from: data!)
                            self.state.running = true
                            self.state.apiKey = key
                            self.state.lastScrobbled = nil
                            self.state.latestPlayed = nil
                            self.state.account = account
                            UserDefaults.standard.set(self.state.running, forKey: AppDelegate.keyRunning)
                            UserDefaults.standard.set(self.state.apiKey, forKey: AppDelegate.keyAPIKey)
                            UserDefaults.standard.set(self.state.lastScrobbled?.timeIntervalSince1970, forKey: AppDelegate.keyLastScrobbled)
                            UserDefaults.standard.set(self.state.latestPlayed?.timeIntervalSince1970, forKey: AppDelegate.keyLatestPlayed)
                            self.render()
                            NSApplication.shared.sendAction(self.oldOkButtonAction!, to: self.oldOkButtonTarget, from: sender)
                        }
                        return
                    } else if rr.statusCode == 401 {
                        DispatchQueue.main.async {
                            self.alert!.informativeText = String(format: "Invalid API key")
                        }
                        return
                    }
                }
            }
            
            DispatchQueue.main.async {
                self.alert!.informativeText = String(format: "Something went wrong. Try again?")
            }
        }
    }
    
    // uppercase API key input
    override func controlTextDidChange(_ obj: Notification) {
        if let text = obj.userInfo!["NSFieldEditor"] as? NSText {
            text.string = text.string.uppercased()
        }
    }
}

func scrobblableItems(from m: Array<ITLibMediaItem>) -> (Array<API.MediaItem>, Date?) {
    var items = Array<API.MediaItem>()
    var latestPlayed: Date? = nil
    for p in m {
        if p.mediaKind != ITLibMediaItemMediaKind.kindSong {
            continue
        }
        if (p.addedDate == nil) {
            // songs that are in playlists, but not explicity +ed
            // have nil addedDate?
            continue
        }
        if (p.lastPlayedDate == nil) {
            // ignore songs that don't have a last played date yet
            continue
        }
        if (p.mediaKind != .kindSong) {
            continue
        }
        items.append(API.MediaItem(fromITLibMediaItem: p))
        if latestPlayed == nil || p.lastPlayedDate! > latestPlayed! {
            latestPlayed = p.lastPlayedDate!
        }
    }
    return (items, latestPlayed)
}


func formatDate(_ t: Date) -> String {
    let y0 = Calendar.current.dateComponents(in: TimeZone.current, from: t).year
    let y1 = Calendar.current.dateComponents(in: TimeZone.current, from: Date(timeIntervalSinceNow: 0)).year
    
    let f = DateFormatter()
    f.dateFormat = y0 != y1 ? "MMM d YYYY, h:mm a" : "MMM d, h:mm a"
    f.timeZone = TimeZone.current
    return f.string(from: t)
}

func sha1(_ d: Data) -> Data {
    var digest = Data(count: Int(CC_SHA1_DIGEST_LENGTH))
    // TODO: understand this warning:
    // Result of call to 'withUnsafeMutableBytes' is unused
    digest.withUnsafeMutableBytes { bytes in
        d.withUnsafeBytes { b in
            CC_SHA1(b, CC_LONG(d.count), bytes)
        }
    }
    return digest
}

// Returns a string representation of the enum value.
// https://developer.apple.com/documentation/ituneslibrary/itlibartworkformat
func artworkFormatString(_ f: ITLibArtworkFormat) -> String {
    switch f {
    case .BMP:
        return "BMP"
    case .GIF:
        return "GIF"
    case .JPEG:
        return "JPEG"
    case .JPEG2000:
        return "JPEG2000"
    case .PICT:
        return "PICT"
    case .PNG:
        return "PNG"
    case .TIFF:
        return "TIFF"
    case .bitmap:
        return "bitmap"
    case .none:
        return "none"
    }
}

