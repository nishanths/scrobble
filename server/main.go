package main

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func main() {
	http.HandleFunc("/", rootHandler)
	// TODO: these can become authenticated handlers via middleware
	http.HandleFunc("/api/scrobble", scrobbleHandler)
	http.HandleFunc("/api/account", accountHandler)
	// http.HandleFunc("/token/generate", tokenHandler)

	appengine.Main()
}

const (
	KindAccount  = "Account"
	KindPlayback = "Playback"
)

// Key: name key, email
type Account struct {
	APIToken string
	Username string
}

type Song struct {
	Duration                   int64 // milliseconds
	Genre, Name, Artist, Album string
	Year                       int64
	Urlp, Urli                 string
}

// Key: ID key, auto-generated
type Playback struct {
	Song      Song
	StartTime int64 // Unix timestamp, seconds
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Acceptable namespaces must match /^[0-9A-Za-z._-]{0,100}$/
// so we can't use an email address as is. So encode it.
func namespace(email string) string {
	return string(hexencode([]byte(email)))
}

func hexencode(b []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(dst, b)
	return dst
}

func accountsForToken(token string) *datastore.Query {
	return datastore.NewQuery(KindAccount).Filter("APIToken=", token)
}

func scrobbleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token, ok := parseAuthenticationToken(ctx, r.Header, w)
	if !ok {
		return
	}

	// Get the account for the token
	akeys, err := accountsForToken(token).KeysOnly().GetAll(ctx, nil)
	if err != nil {
		log.Errorf(ctx, "failed to execute accounts query: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(akeys) == 0 {
		log.Errorf(ctx, "found 0 accounts for token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if len(akeys) != 1 {
		log.Criticalf(ctx, "found multiple accounts for token %q", token)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Switch to the right namespace
	email := akeys[0].StringID()

	ns, err := appengine.Namespace(ctx, namespace(email))
	if err != nil {
		log.Errorf(ctx, "failed to make namespace context: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Add the playback items
	var ps []Playback
	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
		log.Errorf(ns, "failed to json-unmarshal request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pkeys := make([]*datastore.Key, len(ps))
	for i := range ps {
		pkeys[i] = datastore.NewIncompleteKey(ns, KindPlayback, nil)
	}

	// TODO: this is horrible at handling partial failures.
	// Need a non auto-generated key to handle partial
	// failues correctly.
	if _, err := datastore.PutMulti(ns, pkeys, ps); err != nil {
		log.Errorf(ns, "failed to put to datastore: %s", err)
		// This is only correct from a client's perspective if all of the
		// entities failed to be put (e.g. due to a transient network
		// error).
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func generateAPIToken() (string, error) {
	b := make([]byte, 16)
	_, err := cryptorand.Read(b)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read rand")
	}
	return string(hexencode(b)), nil
}

// Parses the token from the header.  If the token could not be parsed,
// returns ("", false), logs an errors, and writes to the response writer.
// Otherwise returns (token, true) with no side-effects.
func parseAuthenticationToken(ctx context.Context, h http.Header, w http.ResponseWriter) (string, bool) {
	header := h.Get("Authentication")
	if header == "" {
		log.Errorf(ctx, "missing Authentication header")
		w.WriteHeader(http.StatusBadRequest)
		return "", false
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		log.Errorf(ctx, "header %q has wrong format", header)
		w.WriteHeader(http.StatusBadRequest)
		return "", false
	}

	if parts[0] != "Token" {
		log.Errorf(ctx, "header %q has wrong format", header)
		w.WriteHeader(http.StatusBadRequest)
		return "", false
	}

	tok := strings.TrimSpace(parts[1])
	if tok == "" {
		log.Errorf(ctx, "header %q has wrong format", header)
		w.WriteHeader(http.StatusBadRequest)
		return "", false
	}

	return tok, true
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token, ok := parseAuthenticationToken(ctx, r.Header, w)
	if !ok {
		return
	}

	var as []Account
	keys, err := accountsForToken(token).GetAll(ctx, &as)
	if err != nil {
		log.Errorf(ctx, "failed to execute accounts query: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(keys) == 0 {
		log.Errorf(ctx, "found 0 accounts for token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if len(keys) != 1 {
		log.Criticalf(ctx, "found multiple accounts for token %q", token)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(AccountResponse{Username: as[0].Username})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
