package main

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/nishanths/applemusic"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func main() {
	http.HandleFunc("/", rootHandler)

	// TODO: these can become authenticated handlers via middleware
	http.HandleFunc("/api/scrobble", scrobbleHandler)
	http.HandleFunc("/api/account", accountHandler)
	http.HandleFunc("/api/scrobbles", scrobblesHandler)

	// http.HandleFunc("/token/generate", tokenHandler)

	appengine.Main()
}

const (
	KindAccount  = "Account"
	KindUsername = "Username"
	KindPlayback = "Playback"
)

// Namespace: [default]
// Key: name key, email
type Account struct {
	APIToken string
	Username string
}

// Namespace: [default]
// Key: name key, username
type Username struct {
	Email string
}

type Song struct {
	Duration                   int64 // milliseconds
	Genre, Name, Artist, Album string
	Year                       int64
	Urlp, Urli                 string
}

// Namespace: account
// Key: ID key, auto-generated
type Playback struct {
	Song      Song
	StartTime int64 // Unix timestamp, seconds
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	maybeUsername := parts[1]
	key := datastore.NewKey(ctx, KindUsername, maybeUsername, 0, nil)
	var u Username

	if err := datastore.Get(ctx, key, &u); err != nil {
		if err == datastore.ErrNoSuchEntity {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ns, err := appengine.Namespace(ctx, namespace(u.Email))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ps []Playback
	if _, err := playbacks().Order("-StartTime").GetAll(ns, &ps); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	for _, p := range ps {
		io.WriteString(w, fmt.Sprintf("%+v\n", p))
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

func playbacks() *datastore.Query {
	return datastore.NewQuery(KindPlayback)
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
		log.Errorf(ctx, "failed to json-unmarshal request: %s", err)
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
		log.Errorf(ctx, "failed to put to datastore: %s", err)
		// This is only correct from a client's perspective if all of the
		// entities failed to be put (e.g. due to a transient network
		// error).
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

const artworkFetchTimeout = 5 * time.Second

func scrobblesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := datastore.NewKey(ctx, KindUsername, username, 0, nil)
	var u Username

	if err := datastore.Get(ctx, key, &u); err != nil {
		if err == datastore.ErrNoSuchEntity {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ns, err := appengine.Namespace(ctx, namespace(u.Email))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ps []Playback
	if _, err := playbacks().Order("-StartTime").GetAll(ns, &ps); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cs := coalesce(ps)

	// use "make" for two reasons:
	// - it marshals as empty array instead of null if there are 0 elements
	// - and obviously, to assign in the goroutines below
	prs := make([]PlaybackResponse, len(cs))
	var g errgroup.Group

	for i := range prs {
		if cs[i].Song.Urlp == "" || cs[i].Song.Urli == "" {
			// TODO: assign more to prs[i]
			prs[i] = PlaybackResponse{cs[i], ""}
			return
		}

		i := i
		g.Go(func() error {
			ctx, cancel := context.WithTimeout(ctx, artworkFetchTimeout)
			defer cancel()
			info, err := fetchAppleMusicInfo(ctx, cs[i].Song.Urlp, cs[i].Song.Urli)
			if err != nil {
				return errors.WithStack(err)
			}
			// TODO: assign more to prs[i]
			prs[i] = PlaybackResponse{cs[i], info.Artwork.HttpsURL}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Errorf(ctx, "failed to make playback response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(prs, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

type CoalescedPlayback struct {
	Song                          Song
	FirstStartTime, LastStartTime int64
	Count                         int
}

type PlaybackResponse struct {
	CoalescedPlayback
	Artwork string
}

func coalesce(ps []Playback) []CoalescedPlayback {
	if len(ps) == 0 {
		return nil
	}

	if len(ps) == 1 {
		return []CoalescedPlayback{{ps[0].Song, ps[0].StartTime, ps[0].StartTime, 1}}
	}

	var cs []CoalescedPlayback
	prev := ps[1]
	count := 1
	i := 2

	for ; i < len(ps); i++ {
		p := ps[i]
		if !equal(prev.Song, p.Song) {
			cs = append(cs, CoalescedPlayback{prev.Song, prev.StartTime, ps[i-1].StartTime, count})
			// reset
			prev = p
			count = 1
		} else {
			count++
		}
	}

	cs = append(cs, CoalescedPlayback{prev.Song, prev.StartTime, ps[i-1].StartTime, count})
	return cs
}

func equal(lhs, rhs Song) bool {
	// if urli and urlp are present in both, use
	// that as the determining factor
	if lhs.Urlp != "" && lhs.Urli != "" && rhs.Urlp != "" && rhs.Urli != "" {
		return lhs.Urlp == rhs.Urlp &&
			lhs.Urli == rhs.Urli
	}
	// otherwise do a ghetto comparison of the attributes
	return lhs.Duration == rhs.Duration &&
		lhs.Genre == rhs.Genre &&
		lhs.Name == rhs.Name &&
		lhs.Artist == rhs.Artist &&
		lhs.Album == rhs.Album &&
		lhs.Year == rhs.Year
}

func fetchAppleMusicInfo(ctx context.Context, urlp, urli string) (applemusic.Info, error) {
	c := urlfetch.Client(ctx)
	u := fmt.Sprintf("https://itunes.apple.com/us/album/%s?i=%s", urlp, urli)

	rsp, err := c.Get(u)
	if err != nil {
		return applemusic.Info{}, errors.Wrapf(err, "failed to get %q", u)
	}

	defer drainAndClose(rsp.Body)

	if rsp.StatusCode/100 != 2 {
		return applemusic.Info{}, errors.Errorf("non-200 status code %d for %q", rsp.StatusCode, u)
	}

	info, err := applemusic.ParseHTML(rsp.Body)
	if err != nil {
		return applemusic.Info{}, errors.Wrapf(err, "failed to parse apple music html")
	}

	return info, nil
}

func drainAndClose(r io.ReadCloser) {
	io.Copy(ioutil.Discard, r)
	r.Close()
}

type AccountResponse struct {
	Username string `json:"username"` // json tag suitable for swift clients
}
