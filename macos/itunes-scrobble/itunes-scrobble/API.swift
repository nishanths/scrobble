//
//  API.swift
//  itunes-scrobble
//
//  Created by Nishanth Shanmugham on 8/28/18.
//  Copyright © 2018 Nishanth Shanmugham. All rights reserved.
//

import Foundation
import iTunesLibrary

class API {
    static let headerAPIKey = "X-Scrobble-API-Key"
    
    static func scrobbleRequest(_ apiKey: String, _ body: Data) -> URLRequest {
        let url = URL(string: String(format: "https://%@/api/v1/scrobble", AppDelegate.baseUrl))!
        var r = URLRequest(url: url)
        r.httpMethod = "POST"
        r.setValue("application/json", forHTTPHeaderField: "Content-Type")
        r.httpBody = body
        setStandardHeaders(&r, apiKey)
        return r
    }
    
    static func accountRequest(_ apiKey: String) -> URLRequest {
        let url = URL(string: String(format: "https://%@/api/v1/account", AppDelegate.baseUrl))!
        var r = URLRequest(url: url)
        r.httpMethod = "GET"
        setStandardHeaders(&r, apiKey)
        return r
    }
    
    static func artworkRequest(_ apiKey: String, _ artworkFormat: ITLibArtworkFormat, _ artwork: Data) -> URLRequest {
        var c = URLComponents()
        c.scheme = "https"
        c.host = AppDelegate.baseUrl
        c.path = "/api/v1/artwork"
        c.queryItems = [URLQueryItem(name: "format", value: artworkFormatString(artworkFormat))]
        
        var r = URLRequest(url: c.url!)
        r.httpMethod = "POST"
        r.httpBody = artwork
        setStandardHeaders(&r, apiKey)
        return r
    }
    
    static func missingArtworkRequest(_ apiKey: String) -> URLRequest {
        let url = URL(string: String(format: "https://%@/api/v1/artwork/missing", AppDelegate.baseUrl))
        var r = URLRequest(url: url!)
        r.httpMethod = "GET"
        setStandardHeaders(&r, apiKey)
        return r
    }
    
    private static func setStandardHeaders(_ r: inout URLRequest, _ apiKey: String) {
        r.setValue(apiKey, forHTTPHeaderField: headerAPIKey)
    }
    
    struct Account: Decodable {
        var username: String
    }
    
    struct MediaItem: Encodable, Hashable {
        var hashValue: Int {
            return persistentID.hashValue
        }
        
        var added: Double?
        var albumTitle: String?
        var sortAlbumTitle: String?
        var artistName: String?
        var sortArtistName: String?
        var genre: String
        var hasArtwork: Bool
        var kind: String?
        var lastPlayed: Double?
        var playCount: UInt
        var releaseDate: Double?
        var sortTitle: String?
        var title: String
        var totalTime: UInt
        var year: UInt
        var persistentID: String
        var artworkHash: String?
        
        var loved = false
        
        init(fromITLibMediaItem i: ITLibMediaItem) {
            self.added = i.addedDate?.timeIntervalSince1970
            self.albumTitle = i.album.title
            self.sortAlbumTitle = i.album.sortTitle
            self.artistName = i.artist?.name
            self.sortArtistName = i.artist?.sortName
            self.genre = i.genre
            self.hasArtwork = i.hasArtworkAvailable
            self.kind = i.kind
            self.lastPlayed = i.lastPlayedDate?.timeIntervalSince1970
            self.playCount = uint(i.playCount)
            self.releaseDate = i.releaseDate?.timeIntervalSince1970
            self.sortTitle = i.sortTitle
            self.title = i.title
            self.totalTime = uint(i.totalTime)
            self.year = uint(i.year)
            self.persistentID = i.persistentID.stringValue
            
            // artwork hash
            guard let a = i.artwork?.imageData else { return }
            guard let f = i.artwork?.imageDataFormat else { return }
            var d = Data()
            d.append(a)
            d.append("|".data(using: .utf8)!) // TODO
            d.append(artworkFormatString(f).data(using: .utf8)!) // TODO
            self.artworkHash = sha1(d).map { String(format: "%d", $0) }.joined()
        }
    }
}
