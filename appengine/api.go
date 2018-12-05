package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

const (
	artworkStorageDirectory = "aw" // Clould Storage directory name for artwork
)

const (
	KindAccount       = "Account"
	KindUsername      = "Username" // to guarantee uniqueness
	KindAPIKey        = "APIKey"   // to guarantee uniqueness
	KindSong          = "Song"
	KindArtworkRecord = "ArtworkRecord"
	KindITunesTrack   = "ITunesTrack"
)

func songident(album, artist, title string, year int) string {
	return fmt.Sprintf("%s|%s|%s|%s",
		base64encode([]byte(album)),
		base64encode([]byte(artist)),
		base64encode([]byte(title)),
		base64encode([]byte((strconv.Itoa(year)))))
}

// Namespace: [default]
// Key: email
type Account struct {
	APIKey   string `json:"apiKey"`
	Username string `json:"username"`
	Private  bool   `json:"private"`
}

// Namespace: account
// Key: see Ident() method
type Song struct {
	// basic properties
	AlbumTitle string        `datastore:",noindex" json:"albumTitle"`
	ArtistName string        `datastore:",noindex" json:"artistName"`
	Title      string        `datastore:",noindex" json:"title"`
	TotalTime  time.Duration `datastore:",noindex" json:"totalTime"`
	Year       int           `json:"year"`

	// sorting fields
	SortAlbumTitle string `json:"-"`
	SortArtistName string `json:"-"`
	SortTitle      string `json:"-"`

	// play info
	LastPlayed int64 `json:"lastPlayed"`
	PlayCount  int   `json:"-"`

	ArtworkHash  string `datastore:",noindex" json:"artworkHash"`
	PreviewURL   string `datastore:",noindex" json:"previewURL"`
	TrackViewURL string `datastore:",noindex" json:"trackViewURL"`

	Loved bool `json:"loved"`
}

func (s *Song) Ident() string {
	return songident(s.AlbumTitle, s.ArtistName, s.Title, s.Year)
}

func (s *Song) iTunesFilled() bool {
	return s.PreviewURL != "" && s.TrackViewURL != ""
}

const headerAPIKey = "X-Scrobble-API-Key"

func accountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.Header.Get(headerAPIKey)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, _, ok := fetchAccountForKey(ctx, key, w)
	if !ok {
		return
	}

	b, err := json.Marshal(struct {
		Username string `json:"username"`
	}{acc.Username})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var accID string
	if u := user.Current(ctx); u == nil {
		key := r.Header.Get(headerAPIKey)
		if key == "" {
			log.Errorf(ctx, "not signed in and missing API key header")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var ok bool
		_, accID, ok = fetchAccountForKey(ctx, key, w)
		if !ok {
			return
		}
	} else {
		accID = u.Email
	}

	log.Infof(ctx, "deleting account %q", accID)

	if err := datastore.RunInTransaction(ctx, func(tx context.Context) error {
		// synchronously delete Username, Account entities
		// NOTE: we intentionally do not delete the APIKey entity because
		// those should always exist to guarantee non-reuse.
		aKey := datastore.NewKey(tx, KindAccount, accID, 0, nil)
		var account Account
		if err := datastore.Get(tx, aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to get account")
		}

		log.Infof(tx, "Account entity %+v for key %s", account, aKey)

		// If the account isn't initialized, the username won't be set
		// and a corresponding Username entity won't exist. So only
		// attempt to delete the Username entity if the username is
		// set.
		if account.Username != "" {
			if err := datastore.Delete(tx, datastore.NewKey(tx, KindUsername, account.Username, 0, nil)); err != nil {
				return errors.Wrapf(err, "failed to delete Username entity")
			}
		}

		if err := datastore.Delete(tx, aKey); err != nil {
			return errors.Wrapf(err, "failed to delete Account entity")
		}
		// asynchronously delete the namespace's entities
		namespace := namespaceID(accID)
		if err := deleteFunc.Call(tx, namespace, KindArtworkRecord); err != nil {
			return errors.Wrapf(err, "failed to call deleteFunc for %s,%s", namespace, KindArtworkRecord)
		}
		if err := deleteFunc.Call(tx, namespace, KindSong); err != nil {
			return errors.Wrapf(err, "failed to call deleteFunc for %s,%s", namespace, KindSong)
		}
		return nil
	}, defaultTxOpts); err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func scrobbledHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, accID, ok := fetchAccountForUsername(ctx, username, w)
	if !ok {
		return
	}

	if acc.Private && !canViewScrobbled(ctx, accID, user.Current(ctx), r.Header) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ns, ok := namespaceFromAccount(ctx, accID, w)
	if !ok {
		return
	}

	q := datastore.NewQuery(KindSong).
		Order("-LastPlayed")

	lovedOnly := r.FormValue("loved") == "true"
	if lovedOnly {
		q = q.Filter("Loved=", true)
	}

	songs := make([]Song, 0) // to marshal as empty JSON array instead of null when there are 0 songs
	_, err := q.GetAll(ns, &songs)
	if err != nil {
		log.Errorf(ns, "failed to fetch songs: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		log.Errorf(ns, "failed to write response: %v", err.Error())
	}
}

func canViewScrobbled(ctx context.Context, forAccountID string, u *user.User, h http.Header) bool {
	if u != nil && u.Email == forAccountID {
		return true
	}

	if key := h.Get(headerAPIKey); key != "" {
		if _, id, _, err := accountForKey(ctx, key); err == nil && id == forAccountID {
			return true
		}
	}

	return false
}

func scrobbleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.Header.Get(headerAPIKey)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, accID, ok := fetchAccountForKey(ctx, key, w)
	if !ok {
		return
	}

	ns, ok := namespaceFromAccount(ctx, accID, w)
	if !ok {
		return
	}

	// Parse request.
	type MediaItem struct {
		Added          float64 `json:"added"`
		AlbumTitle     string  `json:"albumTitle"`
		SortAlbumTitle string  `json:"sortAlbumTitle"`
		ArtistName     string  `json:"artistName"`
		SortArtistName string  `json:"sortArtistName"`
		Genre          string  `json:"genre"`
		HasArtwork     bool    `json:"hasArtwork"`
		Kind           string  `json:"kind"`
		LastPlayed     float64 `json:"lastPlayed"`
		PlayCount      uint    `json:"playCount"`
		ReleaseDate    float64 `json:"releaseDate"`
		SortTitle      string  `json:"sortTitle"`
		Title          string  `json:"title"`
		TotalTime      uint    `json:"totalTime"` // milliseconds
		Year           uint    `json:"year"`
		PersistentID   string  `json:"persistentID"`
		ArtworkHash    string  `json:"artworkHash"`
		Loved          bool    `json:"loved"`
	}
	var mis []MediaItem
	if err := json.NewDecoder(r.Body).Decode(&mis); err != nil {
		log.Errorf(ns, "%v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Convert to Songs.
	songs := make([]Song, len(mis))
	for i, m := range mis {
		sal := m.SortAlbumTitle
		if sal == "" {
			sal = m.AlbumTitle
		}
		sar := m.SortArtistName
		if sar == "" {
			sar = m.ArtistName
		}
		st := m.SortTitle
		if st == "" {
			st = m.Title
		}
		songs[i] = Song{
			AlbumTitle:     m.AlbumTitle,
			ArtistName:     m.ArtistName,
			Title:          m.Title,
			TotalTime:      time.Duration(m.TotalTime) * (time.Millisecond / time.Nanosecond),
			Year:           int(m.Year),
			SortAlbumTitle: sal,
			SortArtistName: sar,
			SortTitle:      st,
			LastPlayed:     int64(m.LastPlayed),
			PlayCount:      int(m.PlayCount),
			ArtworkHash:    m.ArtworkHash,
			Loved:          m.Loved,
		}
	}

	sKeys := make([]*datastore.Key, len(songs))
	var aKeys []*datastore.Key
	for i, s := range songs {
		// Create song key.
		sKeys[i] = datastore.NewKey(ns, KindSong, s.Ident(), 0, nil)
		// Create artwork hash key.
		if h := s.ArtworkHash; h != "" {
			aKeys = append(aKeys, datastore.NewKey(ns, KindArtworkRecord, h, 0, nil))
		}
	}

	// The overall goal is to put the incoming songs and remove any
	// old ones (leftover) that are not present in the set of
	// incoming songs.
	var oldIDs map[string]struct{}
	incomingIDs := setFromSlice(sKeys)
	toRemoveIDs := make(map[string]struct{})

	// Get old song IDs.
	{
		keys, err := datastore.NewQuery(KindSong).KeysOnly().GetAll(ctx, nil)
		if err != nil {
			log.Errorf(ns, "failed to query songs: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		oldIDs = setFromSlice(keys)
	}

	// Compute the IDs to remove.
	// (We do the actual removal below later.)
	for k := range oldIDs {
		if _, exists := incomingIDs[k]; !exists {
			toRemoveIDs[k] = struct{}{}
		}
	}

	// Put songs.
	{
		s := 0
		e := min(s+datastoreLimitPerOp, len(songs))
		chunk := songs[s:e]
		keysChunk := sKeys[s:e]

		for len(chunk) > 0 {
			if _, err := datastore.PutMulti(ns, keysChunk, chunk); err != nil {
				log.Errorf(ns, "failed to put songs: %v", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			s = e
			e = min(s+datastoreLimitPerOp, len(songs))
			chunk = songs[s:e]
			keysChunk = sKeys[s:e]
		}
	}

	var g errgroup.Group

	// Put artwork hash keys.
	g.Go(func() error {
		s := 0
		e := min(s+datastoreLimitPerOp, len(aKeys))
		keysChunk := aKeys[s:e]

		for len(keysChunk) > 0 {
			if _, err := datastore.PutMulti(ns, keysChunk, make([]struct{}, len(keysChunk))); err != nil {
				log.Errorf(ns, "failed to put artwork records: %v", err.Error()) // only log
			}

			s = e
			e = min(s+datastoreLimitPerOp, len(aKeys))
			keysChunk = aKeys[s:e]
		}

		return nil
	})

	// Create tasks to fill in iTunes-related fields.
	for _, s := range songs {
		ident := s.Ident()
		g.Go(func() error {
			if err := fillITunesFunc.Call(ctx, namespaceID(accID), ident); err != nil {
				log.Errorf(ctx, "failed to call fillITunesFunc for %s,%s", namespaceID(accID), ident) // only log
			}
			return nil
		})
	}

	g.Go(func() error {
		toRemoveKeys := sliceFromSet(ns, toRemoveIDs)

		s := 0
		e := min(s+datastoreLimitPerOp, len(toRemoveKeys))
		keysChunk := toRemoveKeys[s:e]

		for len(keysChunk) > 0 {
			if err := datastore.DeleteMulti(ns, keysChunk); err != nil {
				log.Errorf(ns, "failed to delete leftover song keys: %v", err.Error()) // only log
			}

			s = e
			e = min(s+datastoreLimitPerOp, len(toRemoveKeys))
			keysChunk = toRemoveKeys[s:e]
		}

		return nil
	})

	if err := g.Wait(); err != nil { // no such errors are expected (all groups return nil error)
		log.Errorf(ns, "%v", err.Error())
	}

	w.WriteHeader(http.StatusOK)
}

func setFromSlice(ss []*datastore.Key) map[string]struct{} {
	ret := make(map[string]struct{}, len(ss))
	for _, s := range ss {
		ret[s.StringID()] = struct{}{}
	}
	return ret
}

func sliceFromSet(ns context.Context, set map[string]struct{}) []*datastore.Key {
	ret := make([]*datastore.Key, 0, len(set))
	for stringID := range set {
		ret = append(ret, datastore.NewKey(ns, KindSong, stringID, 0, nil))
	}
	return ret
}

func artworkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.Header.Get(headerAPIKey)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate API key
	// TODO: make this more explicit
	_, _, ok := fetchAccountForKey(ctx, key, w)
	if !ok {
		return
	}

	artwork, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(ctx, "failed to read request body: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	format := r.FormValue("format")
	hash := artworkHash(artwork, format)

	// upload to GCS
	name, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get default GCS bucket name: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to create client: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer client.Close()

	wr := client.Bucket(name).Object(artworkStorageDirectory + "/" + hash).NewWriter(ctx)
	wr.Metadata = map[string]string{"format": format}
	if _, err := wr.Write(artwork); err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := wr.Close(); err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof(ctx, "saved artwork hash=%s", hash)

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, hash)
}

func artworkMissingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.Header.Get(headerAPIKey)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, accID, ok := fetchAccountForKey(ctx, key, w)
	if !ok {
		return
	}

	ns, ok := namespaceFromAccount(ctx, accID, w)
	if !ok {
		return
	}

	g, ns := errgroup.WithContext(ns)
	have := make(map[string]struct{})
	want := make(map[string]bool)
	var artworkKeys []*datastore.Key

	g.Go(func() error {
		name, err := file.DefaultBucketName(ns)
		if err != nil {
			return errors.Wrapf(err, "failed to get default GCS bucket name")
		}

		client, err := storage.NewClient(ns)
		if err != nil {
			return errors.Wrapf(err, "failed to create client")
		}
		defer client.Close()

		it := client.Bucket(name).Objects(ns, &storage.Query{Prefix: artworkStorageDirectory + "/"})
		for {
			o, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Errorf(ns, "%v", err.Error()) // only log
				break
			}
			have[strings.TrimPrefix(o.Name, artworkStorageDirectory+"/")] = struct{}{}
		}
		return nil
	})

	g.Go(func() error {
		keys, err := datastore.NewQuery(KindArtworkRecord).KeysOnly().GetAll(ns, nil)
		if err != nil {
			return errors.Wrapf(err, "failed to fetch artwork records")
		}
		artworkKeys = keys
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Errorf(ns, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, k := range artworkKeys {
		if _, ok := have[k.StringID()]; ok {
			continue
		}
		want[k.StringID()] = true
	}

	log.Infof(ns, "%d artwork records with missing data", len(want))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(want); err != nil {
		log.Errorf(ns, "failed to write response: %v", err.Error())
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func accountForKey(ctx context.Context, apiKey string) (Account, string, int, error) {
	var as []Account
	keys, err := datastore.NewQuery(KindAccount).Filter("APIKey=", apiKey).Limit(2).GetAll(ctx, &as)
	if err != nil {
		return Account{}, "", http.StatusInternalServerError, err
	}

	if len(keys) > 1 {
		m := fmt.Sprintf("multiple accounts for API key %q", apiKey)
		log.Criticalf(ctx, m)
		panic(m)
	}

	if len(keys) == 0 {
		return Account{}, "", http.StatusNotFound, fmt.Errorf("no accounts for API key: %s", apiKey)
	}

	return as[0], keys[0].StringID(), 0, nil
}

func accountForUsername(ctx context.Context, username string) (Account, string, int, error) {
	var as []Account
	keys, err := datastore.NewQuery(KindAccount).Filter("Username=", username).Limit(2).GetAll(ctx, &as)
	if err != nil {
		return Account{}, "", http.StatusInternalServerError, err
	}

	if len(keys) > 1 {
		m := fmt.Sprintf("multiple accounts for username %q", username)
		log.Criticalf(ctx, m)
		panic(m)
	}

	if len(keys) == 0 {
		return Account{}, "", http.StatusNotFound, fmt.Errorf("no accounts for username: %s", username)
	}

	return as[0], keys[0].StringID(), 0, nil
}

func fetchAccountForKey(ctx context.Context, apiKey string, w http.ResponseWriter) (Account, string, bool) {
	a, id, code, err := accountForKey(ctx, apiKey)
	if err != nil {
		log.Errorf(ctx, err.Error())
		w.WriteHeader(code)
		return Account{}, "", false
	}
	return a, id, true
}

func fetchAccountForUsername(ctx context.Context, username string, w http.ResponseWriter) (Account, string, bool) {
	a, id, code, err := accountForUsername(ctx, username)
	if err != nil {
		log.Errorf(ctx, err.Error())
		w.WriteHeader(code)
		return Account{}, "", false
	}
	return a, id, true
}

func namespaceID(accountID string) string {
	// allowed namespace pattern: /^[0-9A-Za-z._-]{0,100}$/
	return string(hexencode([]byte(accountID)))
}

func namespaceFromAccount(ctx context.Context, accountID string, w http.ResponseWriter) (context.Context, bool) {
	ns, err := appengine.Namespace(ctx, namespaceID(accountID))
	if err != nil {
		log.Errorf(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return nil, false
	}
	return ns, true
}

func generateAPIKey() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read rand")
	}
	return strings.ToUpper(string(hexencode(b))), nil
}

func hexencode(b []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(dst, b)
	return dst
}

func base64encode(b []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(dst, b)
	return dst
}

func artworkHash(artwork []byte, format string) string {
	h := sha1.New()
	h.Write(artwork)
	h.Write([]byte("|"))
	h.Write([]byte(format))
	sum := h.Sum(nil)

	// TODO: better way?
	var buf bytes.Buffer
	for _, b := range sum {
		buf.WriteString(fmt.Sprintf("%d", b))
	}
	return buf.String()
}
