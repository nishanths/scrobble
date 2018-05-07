package main

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/scrobble", scrobbleHandler)
	http.HandleFunc("/api/account", accountHandler)

	appengine.Main()
}

const (
	KindSong     = "Song"
	KindAccount  = "Account"
	KindPlayback = "Playback"
)

// Key: name key, email
type Account struct {
	APIToken string // transform (email ++ salt)
	Salt     []byte `datastore:",noindex"`
	Username string
}

type Song struct {
	Duration                   int64 // Unix timestamp, seconds
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

const theUsername = "nishanths"
const theEmail = "nishanth.gerrard@gmail.com"

func namespace(email string) string {
	dst := make([]byte, hex.EncodedLen(len(email)))
	hex.Encode(dst, []byte(email))
	return string(dst)
}

func scrobbleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ns, err := appengine.Namespace(ctx, namespace(theEmail))
	if err != nil {
		log.Errorf(ctx, "failed to make namespace context: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: check authentication

	var ps []Playback
	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
		log.Errorf(ctx, "failed to json-unmarshal request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	keys := make([]*datastore.Key, len(ps))
	for i := range ps {
		keys[i] = datastore.NewIncompleteKey(ns, KindPlayback, nil)
	}

	// TODO: this is horrible at handling partial failures.
	// Need a non auto-generated key to handle partial
	// failues correctly.
	if _, err := datastore.PutMulti(ns, keys, ps); err != nil {
		log.Errorf(ns, "failed to put to datastore: %s", err)
		// This is only correct from a client's perspective if all of the
		// entities failed to be put (e.g. due to a transient network
		// error).
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func makeAPIToken(email, salt string) (string, error) {
	panic("not done")
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	b, err := json.Marshal(AccountResponse{Username: theUsername})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
