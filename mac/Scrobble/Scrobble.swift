//
//  Scrobble.swift
//  Scrobble
//
//  Created by Nishanth Shanmugham on 5/6/18.
//  Copyright © 2018 Nishanth Shanmugham. All rights reserved.
//

import Foundation
import Cocoa

struct Account: Decodable {
    var username: String;
}

struct Song: Encodable {
    var duration: Int // milliseconds
    var genre, name, artist, album: String
    var year: Int
    var urlp, urli: String

    init(fromNotification info: Dictionary<AnyHashable, Any>) {
        if let v = info["Total Time"] as? Int {
            self.duration = v
        } else {
            self.duration = 0
        }

        self.genre = info["Genre"] as? String ?? ""
        self.name = info["Name"] as? String ?? ""
        self.artist = info["Artist"] as? String ?? ""
        self.album = info["Album"] as? String ?? ""
        self.year = info["Year"] as? Int ?? 0

        if let v = info["Store URL"] as? String {
            let c = URLComponents(string: v)
            self.urlp = c?.queryItems?.first(where: {$0.name == "p"})?.value ?? ""
            self.urli = c?.queryItems?.first(where: {$0.name == "i"})?.value ?? ""
        } else {
            self.urlp = ""
            self.urli = ""
        }
    }

    static func ==(lhs: Song, rhs: Song) -> Bool {
        // if urli and urlp are present in both, use
        // that as the determining factor
        if lhs.urlp != "" && lhs.urli != "" && rhs.urlp != "" && rhs.urli != "" {
            return lhs.urlp == rhs.urlp
                && lhs.urli == rhs.urli
        }
        // otherwise do a ghetto comparison of the attributes
        return lhs.duration == rhs.duration
            && lhs.genre == rhs.genre
            && lhs.name == rhs.name
            && lhs.artist == rhs.artist
            && lhs.album == rhs.album
            && lhs.year == rhs.year
    }
}

struct Playback: Encodable {
    var song: Song;
    var startTime: TimeInterval; // Unix timestamp
}

class API {
    static let baseUrl = AppDelegate.baseUrl + "/api"

    static func scrobbleRequest(_ token: String, _ body: Data) -> URLRequest {
        let url = URL(string: API.baseUrl + "/scrobble")!
        var r = URLRequest(url: url)
        r.httpMethod = "POST"
        r.setValue("application/json", forHTTPHeaderField: "Content-Type")
        r.setValue("Token " + token, forHTTPHeaderField: "Authentication")
        r.httpBody = body
        return r
    }

    static func accountRequest(_ token: String) -> URLRequest {
        let url = URL(string: API.baseUrl + "/account")!
        var r = URLRequest(url: url)
        r.httpMethod = "GET"
        r.setValue("Token " + token, forHTTPHeaderField: "Authentication")
        return r
    }
}

enum ResponseStatus {
    case success
    case clientError
    case authError
    case generalError
}

class Sender {
    var enc = JSONEncoder()

    private func req(_ p: [Playback], _ token: String) -> URLRequest {
        return API.scrobbleRequest(token, try! self.enc.encode(p))
    }

    func send(_ p: [Playback], _ token: String) -> ResponseStatus {
        // TODO: the scoping here is gross
        var status: ResponseStatus = .generalError

        let semaphore = DispatchSemaphore(value: 0) // to wait synchronously for task completion
        let task = URLSession.shared.dataTask(with: self.req(p, token), completionHandler: {(_, rsp, err) in
            if let h = rsp as? HTTPURLResponse {
                if h.statusCode == 401 {
                    status = .authError
                } else {
                    switch h.statusCode/100 {
                    case 2:
                        status = .success
                    case 4:
                        status = .clientError
                    default:
                        status = .generalError
                    }
                }
            }
            semaphore.signal()
        })

        task.resume()
        semaphore.wait()
        return status
    }
}

protocol WatcherDelegate {
    func onAuthError() -> Void
}

class Watcher {
    var sender = Sender()
    var token: String?
    var delegate: WatcherDelegate?
    var playbacks: [Playback] = []
    let itunes = iTunesApp()

    var lastPaused: Song?
    var lastPausedTime: TimeInterval? // Unix timestamp

    func start(_ token: String) {
        self.token = token
        let d = DistributedNotificationCenter.default()
        d.addObserver(self, selector: #selector(Watcher.trackChange),
                      name: NSNotification.Name("com.apple.iTunes.playerInfo"), object: nil)
    }

    func stop() {
        let d = DistributedNotificationCenter.default()
        d.removeObserver(self, name: NSNotification.Name("com.apple.iTunes.playerInfo"), object: nil)
    }

    private func isPlaying(_ info: Dictionary<AnyHashable, Any>) -> Bool {
        if let state = info["Player State"] as? String {
            return state == "Playing"
        }
        return false
    }

    private func isPaused(_ info: Dictionary<AnyHashable, Any>) -> Bool {
        if let state = info["Player State"] as? String {
            return state == "Paused"
        }
        return false
    }

    // The given song is considered a restart if it's the same
    // as the last known paused song and sufficient time hasn't
    // passed since the last known pause time.
    private func isRestart(_ song: Song, _ now: TimeInterval) -> Bool {
        if self.lastPaused == nil {
            // no last paused song, so the given song
            // can't be a restart
            return false
        }
        let isRecent = (now - self.lastPausedTime!) < TimeInterval(4 * 60 * 60)
        return song == self.lastPaused! && isRecent
    }

    // All of the handling in this function synchronous so
    // that we don't handle multiple notifications concurrently
    // (which will not bode well with the sending strategy in use).
    @objc private func trackChange(n: Notification)  {
        let info = n.userInfo!
        let now = Date().timeIntervalSince1970;

        if !itunes.currentTrackIsSong() {
            return
        }
        
        if isPaused(info) {
            // record the latest pause info
            self.lastPaused = Song(fromNotification: info)
            self.lastPausedTime = now
            return
        }

        // guard against unknown states (if any)
        if !isPlaying(info) {
            return
        }

        let s = Song(fromNotification: info)

        if isRestart(s, now) {
            return
        }

        // since this wasn't a restart, it means that
        // there is no last paused song anymore.
        self.lastPaused = nil
        self.lastPausedTime = 0

        playbacks.append(Playback(song: s, startTime: floor(now)))
        let status = sender.send(playbacks, token!)

        switch status {
        case .success, .clientError:
            playbacks.removeAll()
        case .authError:
            delegate?.onAuthError()
        default:
            () // nothing to do
        }
    }
}
