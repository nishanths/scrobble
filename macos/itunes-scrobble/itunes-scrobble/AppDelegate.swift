//
//  AppDelegate.swift
//  itunes-scrobble
//
//  Created by Nishanth Shanmugham on 8/26/18.
//  Copyright © 2018 Nishanth Shanmugham. All rights reserved.
//

import Cocoa
import iTunesLibrary

enum ErrorKind {
    case Auth
    case Other
}

// State is the state of the application.
//
// TODO: there's plenty of concurrent accesses of these vars,
// but for the most part the asynchronous functions performing the
// concurrent accesses run far apart in time from each other.
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
    // whether the latest request resulted in an error
    var error: ErrorKind?
    
    // NOTE: when adding a new field, you may also need to handle its reset
    // behavior in clearAPIKey()
    // TODO: this is gross
}

struct Constants {
    static let BaseUrl = "selective-scrobble.appspot.com"
    static let HelpLink = "https://scrobble.allele.cc"
}

func profileLink(_ username: String) -> String {
    return String(format: "https://scrobble.allele.cc/u/%@", username)
}

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate, NSTextFieldDelegate, NSAlertDelegate {
    private static let menuIconName = "itunes-scrobble-18x18" // size from https://stackoverflow.com/a/33708433
    private static let shortVersion = Bundle.main.object(forInfoDictionaryKey: "CFBundleShortVersionString") as! String
    
    // Keys for information saved to UserDefaults.
    private static let (keyAPIKey, keyRunning) = ("apiKey", "running")
    
    // Scrobble timer frequnecies.
    // The timer fires frequently, but scrobbling happens less often.
    private static let timerFreq: TimeInterval = 60 * 10;
    private static let scrobbleFreq: TimeInterval = 24 * 60 * 60 * 1;
    private var timer: Timer? = nil
    
    private let statusBarItem = NSStatusBar.system.statusItem(withLength:NSStatusItem.squareLength)
    private let pauseItem = NSMenuItem(title: "", action: nil, keyEquivalent: "")
    private let statusItem = NSMenuItem(title: "", action: nil, keyEquivalent: "")
    private let secondaryStatusItem = NSMenuItem(title: "", action: nil, keyEquivalent: "")
    private let multiItem = NSMenuItem(title: "", action: nil, keyEquivalent: "")
    private let profileLinkItem = NSMenuItem(title: "", action: nil, keyEquivalent: "")
    
    private var prevState: State? = nil
    private var state = State(running: false,
                              apiKey: nil,
                              lastScrobbled: nil,
                              latestPlayed: nil,
                              account: nil,
                              scrobbling: false,
                              error: nil)
    
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
        state.running = UserDefaults.standard.bool(forKey: AppDelegate.keyRunning)
        state.apiKey = UserDefaults.standard.string(forKey: AppDelegate.keyAPIKey)
        render()
        
        // initially fetch account info
        guard let key = state.apiKey else { return }
        let task = URLSession.shared.dataTask(with: API.accountRequest(key)) {(data, rsp, err) in
            // TODO: handling of failure scenarios
            guard err == nil else { return }
            guard let rr = rsp as! HTTPURLResponse? else { return }
            if (rr.statusCode == 200) {
                DispatchQueue.main.async {
                    let account = try? JSONDecoder().decode(API.Account.self, from: data!)
                    self.state.account = account
                    self.state.error = nil
                    self.render()
                }
            } else if rr.statusCode == 404 {
                DispatchQueue.main.async {
                    self.state.error = .Auth
                    self.render()
                }
                return
            }
        }
        task.resume()
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
        m.addItem(profileLinkItem)
        m.addItem(NSMenuItem.separator())
        m.addItem(statusItem)
        m.addItem(secondaryStatusItem)
        m.addItem(NSMenuItem.separator())
        m.addItem(v)
        m.addItem(NSMenuItem(title: "Quit", action: #selector(NSApplication.terminate(_:)), keyEquivalent: ""))
        
        return m
    }
    
    // TODO: there are several assignment to State's fields followed by
    // a render() call. Instead make it so that setting State's fields
    // automatically calls render().
    private func render() {
        defer { prevState = state }
        
        profileLinkItem.title = "Browse your scrobbles"
        profileLinkItem.action = #selector(openProfile(_:))
        
        // TODO: clean this up, gosh it's gnarly
        if state.apiKey == nil {
            assert(!state.running)
            multiItem.title = "Start scrobbling..."
            multiItem.action =  #selector(enterAPIKeyAction(_:))
            profileLinkItem.isHidden = true
            pauseItem.isHidden = true
        } else {
            if let a = state.account {
                multiItem.title = String(format: "Signed in as %@...", a.username)
                multiItem.action = #selector(scrobblingAsAction(_:))
                profileLinkItem.isHidden = false
            } else if let err = state.error {
                switch err {
                case .Auth:
                    multiItem.title = String(format: "Re-enter API Key")
                    multiItem.action = #selector(clearThenEnterAPIKeyAction(_:))
                case .Other:
                    multiItem.title = String(format: "Remove API Key and Sign out")
                    multiItem.action = #selector(clearAPIKeyAction(_:))
                }
                profileLinkItem.isHidden = true
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
        if let err = state.error {
            switch err {
            case .Auth:
                statusItem.title = "Error: API Key outdated?"
            case .Other:
                statusItem.title = "Error: Failed to scrobble"
            }
            statusItem.isHidden = false
        } else if state.scrobbling {
            statusItem.title = String(format: "Scrobbling now...")
            statusItem.isHidden = false
        } else {
            if let ls = state.lastScrobbled {
                statusItem.title = String(format: "Last scrobbled: %@", formatDate(ls)) // extra spaces for text alignment with secondary status item
                statusItem.isHidden = false
            } else {
                statusItem.isHidden = true
            }
        }
        
        // Secondary status item
        if state.error != nil || state.scrobbling {
            secondaryStatusItem.isHidden = true
        } else {
            if let lp = state.latestPlayed {
                secondaryStatusItem.title = String(format: "Latest song time: %@", formatDate(lp))
                secondaryStatusItem.isHidden = false
            } else {
                secondaryStatusItem.isHidden = true
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
        lib!.reloadData()
        render()
        
        var (items, latest) = scrobblableItems(from: lib!.allMediaItems)
        let loved = lovedItems(in: lib!)
        for idx in 0..<items.count {
            if loved.contains(items[idx]) {
                items[idx].loved = true
            }
        }

        guard let data = try? JSONEncoder().encode(items) else { return }
        
        let task = URLSession.shared.dataTask(with: API.scrobbleRequest(state.apiKey!, data)) {(data, rsp, err) in
            defer {
                DispatchQueue.main.async {
                    self.state.scrobbling = false
                    self.render()
                }
            }
            
            guard err == nil else { return }
            guard let r = rsp as! HTTPURLResponse? else { return }
            
            if r.statusCode == 404 {
                DispatchQueue.main.async {
                    self.state.error = .Auth
                    self.render()
                }
                return
            }
            
            if r.statusCode == 200 {
                DispatchQueue.main.async {
                    self.state.lastScrobbled = Date(timeIntervalSinceNow: 0)
                    self.state.latestPlayed = latest
                    self.state.error = nil
                    self.render()
                }
                self.handleArtwork(self.lib!)
                return
            }
            // any other status code
            DispatchQueue.main.async {
                self.state.error = .Other
                self.render()
            }
        }
        task.resume()
    }
    
    private func handleArtwork(_ lib: ITLibrary) {
        guard let key = state.apiKey else { return }
        let task = URLSession.shared.dataTask(with: API.missingArtworkRequest(key)) {(data, rsp, err) in
            if err != nil {
                return
            }
            guard let r = rsp as! HTTPURLResponse? else {
                return
            }
            if r.statusCode != 200 {
                return
            }
            
            guard let incomingHashes = try? JSONDecoder().decode([String: Bool].self, from: data!) else { return }
            
            DispatchQueue.global(qos: .userInitiated).async {
                var items = Dictionary<String, ITLibMediaItem>()
                for p in lib.allMediaItems {
                    guard let hash = API.MediaItem(fromITLibMediaItem: p).artworkHash else { continue }
                    items[hash] = p
                }
                
                // send artwork that the server is missing
                for (h, _) in incomingHashes {
                    guard self.state.running else { return }
                    guard let key = self.state.apiKey else { return }
                    guard let p = items[h] else { continue }
                    guard let d = p.artwork?.imageData else { continue }
                    guard let f = p.artwork?.imageDataFormat else { continue }
                    let upload = URLSession.shared.dataTask(with: API.artworkRequest(key, f, d)) {(_, _, _) in}
                    upload.resume()
                }
            }
        }
        task.resume()
    }
    
    private func clearAPIKey() {
        state.running = false
        state.apiKey = nil
        state.lastScrobbled = nil
        state.latestPlayed = nil
        state.account = nil
        state.scrobbling = false
        state.error = nil // TODO: need better discrimination errors; it doesn't feel right to unconditionally clear this
        UserDefaults.standard.set(state.running, forKey: AppDelegate.keyRunning)
        UserDefaults.standard.set(state.apiKey, forKey: AppDelegate.keyAPIKey)
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
        let leeway: TimeInterval = 60
        if state.lastScrobbled != nil && abs(state.lastScrobbled!.timeIntervalSinceNow) < AppDelegate.scrobbleFreq - leeway {
            // already scrobbled in near past
            return
        }
        scrobble()
    }
    
    @objc private func clearAPIKeyAction(_ sender: Any?) {
        assert(state.apiKey != nil)
        clearAPIKey()
    }
    
    @objc private func scrobblingAsAction(_ sender: Any?) {
        assert(state.apiKey != nil && state.account != nil)
        let a = NSAlert()
        a.alertStyle = .informational
        a.messageText = String(format: "Scrobbling as %@", state.account!.username)
        a.informativeText = String(format: "%@", profileLink(state.account!.username))
        a.showsSuppressionButton = false
        a.showsHelp = true
        a.delegate = self
        a.addButton(withTitle: "Close")
        a.addButton(withTitle: "Remove API Key and Sign out")
        
        let result = a.runModal()
        switch result {
        case NSApplication.ModalResponse.alertFirstButtonReturn:
            break // nothing to do
        case NSApplication.ModalResponse.alertSecondButtonReturn:
            clearAPIKey()
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
        alert!.messageText = "Enter API Key"
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
        textField!.font = NSFont(name: "Menlo", size: NSFont.systemFontSize)
        textField!.placeholderString = "D1A3903GB"
        alert!.accessoryView = textField
        alert!.window.initialFirstResponder = alert!.accessoryView
        
        alert!.runModal()
    }
    
    func alertShowHelp(_ alert: NSAlert) -> Bool {
        if let u = URL(string: Constants.HelpLink) {
            NSWorkspace.shared.open(u)
        }
        return true
    }
    
    @objc private func openProfile(_ sender: Any?) {
        if let u = URL(string: profileLink(state.account!.username)) {
            NSWorkspace.shared.open(u)
        }
    }
    
    @objc private func okButtonAction(_ sender: Any?) {
        let key = textField!.stringValue
        if key.isEmpty {
            return
        }
        
        let task = URLSession.shared.dataTask(with: API.accountRequest(key)) {(data, rsp, err) in
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
                            self.state.error = nil
                            UserDefaults.standard.set(self.state.running, forKey: AppDelegate.keyRunning)
                            UserDefaults.standard.set(self.state.apiKey, forKey: AppDelegate.keyAPIKey)
                            self.render()
                            NSApplication.shared.sendAction(self.oldOkButtonAction!, to: self.oldOkButtonTarget, from: sender)
                        }
                        return
                    } else if rr.statusCode == 404 {
                        DispatchQueue.main.async {
                            self.alert!.informativeText = String(format: "Invalid API Key.")
                        }
                        return
                    }
                    // any other status code
                    DispatchQueue.main.async {
                        self.alert!.informativeText = String(format: "Something went wrong but not on your end.")
                    }
                }
            }
            
            DispatchQueue.main.async {
                self.alert!.informativeText = String(format: "Something went wrong. Try again?")
            }
        }
        task.resume()
    }
    
    @objc private func clearThenEnterAPIKeyAction(_ sender: Any?) {
        clearAPIKey()
        enterAPIKeyAction(sender)
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
        // There is a bug (somewhere) that leads to some songs having last
        // played times far into the future, for instance, in the year 2040.
        // So ignore such songs.
        //
        // See related: 27833a2309d3a94a9ddbe37791174d7e00c7737d
        if p.lastPlayedDate!.timeIntervalSinceNow > 365 * 24 * 60 * 60 {
            continue
        }
        items.append(API.MediaItem(fromITLibMediaItem: p))
        if latestPlayed == nil || p.lastPlayedDate! > latestPlayed! {
            latestPlayed = p.lastPlayedDate!
        }
    }
    return (items, latestPlayed)
}

func lovedItems(in lib: ITLibrary) -> Set<API.MediaItem> {
    if let p = lib.allPlaylists.first(where: { $0.distinguishedKind == .kindLovedSongs }) {
        return Set(p.items.map({ API.MediaItem(fromITLibMediaItem: $0) }))
    }
    if let p = lib.allPlaylists.first(where: { $0.name == "Loved" }) {
        return Set(p.items.map({ API.MediaItem(fromITLibMediaItem: $0) }))
    }
    return Set()
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

