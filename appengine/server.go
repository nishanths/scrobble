package main

import "net/http"

func RegisterHandlers() {
	http.Handle("/", withHTTPS(http.HandlerFunc(rootHandler)))
	http.Handle("/u/", withHTTPS(http.HandlerFunc(uHandler)))
	http.Handle("/initializeAccount", withHTTPS(http.HandlerFunc(initializeAccountHandler)))
	http.Handle("/newAPIKey", withHTTPS(http.HandlerFunc(newAPIKeyHandler)))
	http.Handle("/setPrivacy", withHTTPS(http.HandlerFunc(setPrivacyHandler)))

	http.HandleFunc("/api/v1/scrobbled", scrobbledHandler)
	http.HandleFunc("/api/v1/scrobble", scrobbleHandler)
	http.HandleFunc("/api/v1/account", accountHandler)
	http.HandleFunc("/api/v1/artwork", artworkHandler)
	http.HandleFunc("/api/v1/artwork/missing", artworkMissingHandler)
}

func withHTTPS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeRedirectHTTPS(w, r) {
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Redirect requests with a "http" scheme to "https".
func maybeRedirectHTTPS(w http.ResponseWriter, r *http.Request) bool {
	u := *r.URL
	if u.Scheme != "http" {
		return false
	}
	u.Scheme = "https"
	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
	return true
}
