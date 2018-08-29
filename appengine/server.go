package main

import "net/http"

func RegisterHandlers() {
	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/api/v1/scrobbled", scrobbledHandler)
	http.HandleFunc("/api/v1/scrobble", scrobbleHandler)
	http.HandleFunc("/api/v1/account", accountHandler)
	http.HandleFunc("/api/v1/artwork", artworkHandler)
	http.HandleFunc("/api/v1/artwork/missing", artworkMissingHandler)
}
