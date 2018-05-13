package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/nishanths/applemusic"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	appengineuser "google.golang.org/appengine/user"
)

func RegisterHandlers() {
	http.HandleFunc("/", rootHandler)

	// available, but not used in the app
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.HandleFunc("/api/scrobble", scrobbleHandler)
	http.HandleFunc("/api/scrobbles", scrobblesHandler)
	http.HandleFunc("/api/account", accountHandler)
	http.HandleFunc("/api/token/generate", accountHandler)
}

var disallowedUsernames = []string{
	"login", "logout", "api", "account", "internal", "task", "cron", "username",
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

type AccountResponse struct {
	Username string `json:"username"` // json tag suitable for swift clients
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

// SongX is an enhaced version of Song.
type SongX struct {
	Song
	ArtworkURL string
	ArtistURL  string
	AlbumURL   string
}

// Namespace: account
// Key: ID key, auto-generated
type Playback struct {
	Song      Song
	StartTime int64 // Unix timestamp, seconds
}

type CoalescedPlayback struct {
	Song       SongX
	StartTimes []int64
}

type PlaybackResponse struct {
	Playbacks []CoalescedPlayback
	Cursor    string
}

func ensureAccount(w http.ResponseWriter, r *http.Request) (*datastore.Key, Account, error) {
	ctx := appengine.NewContext(r)

	u := appengineuser.Current(ctx)
	if u == nil {
		panic(fmt.Sprintf("invalid state? no user while handling %q", r.URL.Path))
	}

	var acc Account
	key := datastore.NewKey(ctx, KindAccount, u.Email, 0, nil)

	// Ensure the Account entity.
	err := datastore.RunInTransaction(ctx, func(tx context.Context) error {
		err := datastore.Get(tx, key, &acc)
		if err == nil {
			return nil // already exists
		}
		if err != datastore.ErrNoSuchEntity {
			return err
		}
		tok, err := generateAPIToken()
		if err != nil {
			return err
		}
		acc.APIToken = tok
		key, err = datastore.Put(tx, key, &acc)
		if err != nil {
			return err
		} else {
			return nil
		}
	}, nil)

	return key, acc, err
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u := appengineuser.Current(ctx)
	if u != nil {
		// already logged in
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	url, err := appengineuser.LoginURL(ctx, "/")
	if err != nil {
		log.Errorf(ctx, "failed to make url: %s", err)
		http.Error(w, "failed to make login url", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<a href="%s">Continue with Google</a>`, url)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u := appengineuser.Current(ctx)
	if u == nil {
		// already logged out
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	url, err := appengineuser.LogoutURL(ctx, "/")
	if err != nil {
		log.Errorf(ctx, "failed to make url: %s", err)
		http.Error(w, "failed to make logout url", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<a href="%s">Sign out</a>`, url)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/" {
		log.Debugf(ctx, "handling index path")
		handleIndex(ctx, w, r)
		return
	}

	parts := strings.Split(r.URL.Path, "/") // parts[0] should be ""
	log.Debugf(ctx, "handling username path using username %q", parts[1])
	handleUsername(ctx, w, r, parts[1])
}

func handleIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := appengineuser.Current(ctx)

	// if logged in and a username has been set, redirect to the '/username' page
	// if logged in without a username, set up a username
	if u != nil {
		key, acc, err := ensureAccount(w, r)
		if err != nil {
			log.Errorf(ctx, "failed to ensure account: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Debugf(ctx, "successfully ensured account %s", key)

		if acc.Username == "" {
			fmt.Fprintf(w, `Pick a username`)
			return
		}

		http.Redirect(w, r, "/"+acc.Username, http.StatusTemporaryRedirect)
		return
	}

	// show the home page
	url, err := appengineuser.LoginURL(ctx, "/")
	if err != nil {
		log.Errorf(ctx, "failed to make url: %s", err)
		http.Error(w, "failed to make login url", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<a href="%s">Continue with Google</a>`, url)
}

func handleUsername(ctx context.Context, w http.ResponseWriter, r *http.Request, maybeUsername string) {
	key := datastore.NewKey(ctx, KindUsername, maybeUsername, 0, nil)
	var u Username

	if err := datastore.Get(ctx, key, &u); err != nil {
		log.Errorf(ctx, "failed to get username %s: %s", key, err)
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

	_ = ns

	// TODO: render the page
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

func scrobbleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token, ok := parseAuthenticationToken(ctx, r.Header, w)
	if !ok {
		return
	}

	// Get the account for the token
	akeys, err := AccountsForToken(token).KeysOnly().GetAll(ctx, nil)
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

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token, ok := parseAuthenticationToken(ctx, r.Header, w)
	if !ok {
		return
	}

	var as []Account
	keys, err := AccountsForToken(token).GetAll(ctx, &as)
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

const appleMusicFetchTimeout = 5 * time.Second

func scrobblesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if appengine.IsDevAppServer() {
		devScrobblesHandler(w, r)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cursorStart, err := datastore.DecodeCursor(r.FormValue("cursor"))
	if err != nil {
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
	it := Playbacks().Order("-StartTime").Start(cursorStart).Run(ns)
	const queryLimit = 45
	coalesceCount := 0
	var cursorNext *datastore.Cursor

	for {
		var p Playback
		_, err := it.Next(&p)
		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "failed to get next item: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// TODO: the coalescing can happen inline here, instead
		// of in the separate function after
		if len(ps) == 0 || !equal(ps[len(ps)-1].Song, p.Song) {
			coalesceCount++
		}
		ps = append(ps, p)
		if coalesceCount > queryLimit {
			var err error
			cursorNext, err = it.Cursor()
			if err != nil {
				log.Errorf(ctx, "failed to get cursor: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			break
		}
	}

	cs := coalesce(ps)

	// Gather the unique urli and urlp values, then fetch
	// apple music info for them
	type urlpi struct{ urlp, urli string }
	var uniquesMu sync.Mutex
	uniques := make(map[urlpi]applemusic.Info)

	cannotAppleFetch := func(s Song) bool { return s.Urlp == "" || s.Urli == "" }

	for _, c := range cs {
		if cannotAppleFetch(c.Song.Song) {
			continue
		}
		k := urlpi{c.Song.Urlp, c.Song.Urli}
		uniques[k] = applemusic.Info{} // whatever -- replaced below
	}

	var g errgroup.Group

	for k := range uniques {
		k := k
		g.Go(func() error {
			ctx, cancel := context.WithTimeout(ctx, appleMusicFetchTimeout)
			defer cancel()
			info, err := fetchAppleMusicInfo(ctx, k.urlp, k.urli)
			if err != nil {
				// just log...
				log.Warningf(ctx, "failed to fetch apple music info: %s", err)
				return nil
			}
			uniquesMu.Lock()
			defer uniquesMu.Unlock()
			uniques[k] = info
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Errorf(ctx, "failed to fetch apple music info: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Enhance.
	for i := range cs {
		if cannotAppleFetch(cs[i].Song.Song) {
			continue
		}
		k := urlpi{cs[i].Song.Urlp, cs[i].Song.Urli}
		cs[i].Song.ArtworkURL = uniques[k].Artwork.HttpsURL
		cs[i].Song.AlbumURL = uniques[k].AlbumURL
		cs[i].Song.ArtistURL = uniques[k].ArtistURL
	}

	rsp := PlaybackResponse{
		Playbacks: cs,
		Cursor: func() string {
			// TODO: don't use the raw cursor
			if cursorNext != nil {
				return cursorNext.String()
			}
			return ""
		}(),
	}

	if rsp.Playbacks == nil {
		rsp.Playbacks = make([]CoalescedPlayback, 0) // json-marshal to empty array instead of null
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rsp); err != nil {
		log.Errorf(ctx, "failed to json-encode response")
		return
	}
}

// Coalesces playbacks. The returned items will not have their
// SongX specific fields set.
// TODO: the types here are gross.
func coalesce(ps []Playback) []CoalescedPlayback {
	if len(ps) == 0 {
		return nil
	}

	cs := []CoalescedPlayback{
		{SongX{Song: ps[0].Song}, []int64{ps[0].StartTime}},
	}
	current := ps[0]

	for i := 1; i < len(ps); i++ {
		p := ps[i]
		if !equal(current.Song, p.Song) {
			cs = append(cs, CoalescedPlayback{SongX{Song: p.Song}, []int64{p.StartTime}})
			current = p
		} else {
			cs[len(cs)-1].StartTimes = append(cs[len(cs)-1].StartTimes, p.StartTime)
		}
	}

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

func appleMusicPreviewURL(urlp, urli string) string {
	return fmt.Sprintf("https://itunes.apple.com/us/album/%s?i=%s", urlp, urli)
}

func fetchAppleMusicInfo(ctx context.Context, urlp, urli string) (applemusic.Info, error) {
	c := urlfetch.Client(ctx)
	u := appleMusicPreviewURL(urlp, urli)

	rsp, err := c.Get(u)
	if err != nil {
		return applemusic.Info{}, errors.Wrapf(err, "failed to get %q", u)
	}

	defer drainAndClose(rsp.Body)

	if rsp.StatusCode/100 != 2 {
		return applemusic.Info{}, errors.Errorf("non-2xx status code %d for %q", rsp.StatusCode, u)
	}

	info, err := applemusic.ParseHTML(rsp.Body)
	if err != nil {
		return applemusic.Info{}, errors.Wrapf(err, "failed to parse apple music html")
	}

	log.Debugf(ctx, "got info for %s %+v", u, info)
	return info, nil
}

func drainAndClose(r io.ReadCloser) {
	io.Copy(ioutil.Discard, r)
	r.Close()
}
