package main

import (
	"bytes"
	"context"
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

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
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
)

// Namespace: [default]
// Key: email
type Account struct {
	APIKey   string
	Username string
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
	SortAlbumTitle string `datastore:",noindex" json:"sortAlbumTitle"`
	SortArtistName string `datastore:",noindex" json:"sortArtistName"`
	SortTitle      string `datastore:",noindex" json:"sortTitle"`

	// play info
	LastPlayed int64 `json:"lastPlayed"`
	PlayCount  int   `json:"playCount"`

	ArtworkHash string `datastore:",noindex" json:"artworkHash"`
}

func (s *Song) Ident() string {
	return fmt.Sprintf("%s|%s|%s|%s",
		base64encode([]byte(s.AlbumTitle)),
		base64encode([]byte(s.ArtistName)),
		base64encode([]byte(s.Title)),
		base64encode([]byte(strconv.Itoa(s.Year))))
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
		w.WriteHeader(http.StatusUnauthorized)
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

func scrobbledHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, accID, ok := fetchAccountForUsername(ctx, username, w)
	if !ok {
		return
	}

	ns, ok := namespaceFromAccount(ctx, accID, w)
	if !ok {
		return
	}

	var songs []Song
	_, err := datastore.NewQuery(KindSong).GetAll(ns, &songs)
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

func scrobbleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.Header.Get(headerAPIKey)
	if key == "" {
		w.WriteHeader(http.StatusUnauthorized)
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
		ArtworkHash    string  `json:"artworkHash"` // md5(artworkData + '|' + artworkFormatString)
	}
	var mis []MediaItem
	if err := json.NewDecoder(r.Body).Decode(&mis); err != nil {
		log.Errorf(ns, "%v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	songs := make([]Song, len(mis))
	for i, m := range mis {
		songs[i] = Song{
			AlbumTitle:     m.AlbumTitle,
			ArtistName:     m.ArtistName,
			Title:          m.Title,
			TotalTime:      time.Duration(m.TotalTime) * (time.Millisecond / time.Nanosecond),
			Year:           int(m.Year),
			SortAlbumTitle: m.SortAlbumTitle,
			SortArtistName: m.SortArtistName,
			SortTitle:      m.SortTitle,
			LastPlayed:     int64(m.LastPlayed),
			PlayCount:      int(m.PlayCount),
			ArtworkHash:    m.ArtworkHash,
		}
	}

	const n = 500 // datastore limit per operation
	s := 0
	e := min(s+n, len(songs))
	chunk := songs[s:e]

	for len(chunk) > 0 {
		keys := make([]*datastore.Key, len(chunk))
		var aKeys []*datastore.Key

		for i := range chunk {
			keys[i] = datastore.NewKey(ns, KindSong, chunk[i].Ident(), 0, nil)
			if h := chunk[i].ArtworkHash; h != "" {
				aKeys = append(aKeys, datastore.NewKey(ns, KindArtworkRecord, h, 0, nil))
			}
		}

		if _, err := datastore.PutMulti(ns, keys, chunk); err != nil {
			log.Errorf(ns, "failed to put songs: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := datastore.PutMulti(ns, aKeys, make([]struct{}, len(aKeys))); err != nil {
			log.Errorf(ns, "failed to put artwork records: %v", err.Error()) // only log
		}

		s = e
		e = min(s+n, len(songs))
		chunk = songs[s:e]
	}

	w.WriteHeader(http.StatusOK)
}

func artworkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.Header.Get(headerAPIKey)
	if key == "" {
		w.WriteHeader(http.StatusUnauthorized)
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
		w.WriteHeader(http.StatusUnauthorized)
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

func fetchAccountForKey(ctx context.Context, apiKey string, w http.ResponseWriter) (Account, string, bool) {
	var as []Account
	keys, err := datastore.NewQuery(KindAccount).Filter("APIKey=", apiKey).Limit(1).GetAll(ctx, &as)
	if err != nil {
		log.Errorf(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return Account{}, "", false
	}

	if len(keys) == 0 {
		log.Infof(ctx, "no accounts for API key: %s", apiKey)
		w.WriteHeader(http.StatusUnauthorized)
		return Account{}, "", false
	}

	return as[0], keys[0].StringID(), true
}

func fetchAccountForUsername(ctx context.Context, username string, w http.ResponseWriter) (Account, string, bool) {
	var as []Account
	keys, err := datastore.NewQuery(KindAccount).Filter("Username=", username).Limit(1).GetAll(ctx, &as)
	if err != nil {
		log.Errorf(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return Account{}, "", false
	}

	if len(keys) == 0 {
		log.Infof(ctx, "no accounts for username: %s", username)
		w.WriteHeader(http.StatusNotFound)
		return Account{}, "", false
	}

	return as[0], keys[0].StringID(), true
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

func generateAPIToken() (string, error) {
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
