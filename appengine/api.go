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

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/nishanths/scrobble/appengine/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine/user"
	"google.golang.org/genproto/googleapis/cloud/tasks/v2beta2"
)

const (
	artworkStorageDirectory = "aw" // Cloud Storage directory name for artwork
)

const (
	KindAccount       = "Account"
	KindUsername      = "Username" // stored for uniqueness checking
	KindAPIKey        = "APIKey"   // stored for uniqueness checking
	KindSongParent    = "SongParent"
	KindSong          = "Song"
	KindArtworkRecord = "ArtworkRecord" // for fast key-only determination of missing artwork
	KindITunesTrack   = "ITunesTrack"
	KindSecret        = "Secret"
)

func songparentident(t time.Time, u uuid.UUID) string {
	return fmt.Sprintf("%d|%s", t.Unix(), u.String())
}

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
// Key: <unix seconds>|<uuid>
type SongParent struct {
	Complete bool
	Created  int64 // unix seconds
}

func songParentKey(namespace, ident string) *datastore.Key {
	return &datastore.Key{Kind: KindSongParent, Name: ident, Namespace: namespace}
}

// Namespace: account
// Key: see Ident() method
type Song struct {
	// basic properties
	AlbumTitle  string        `datastore:",noindex" json:"albumTitle"`
	ArtistName  string        `datastore:",noindex" json:"artistName"`
	Title       string        `datastore:",noindex" json:"title"`
	TotalTime   time.Duration `datastore:",noindex" json:"totalTime"`
	Year        int           `json:"year"`
	ReleaseDate int64         `datastore:",noindex" json:"releaseDate"` // unix seconds

	// sorting fields
	SortAlbumTitle string `json:"-"`
	SortArtistName string `json:"-"`
	SortTitle      string `json:"-"`

	// play info
	LastPlayed int64 `json:"lastPlayed"` // unix seconds
	PlayCount  int   `json:"-"`

	ArtworkHash string `datastore:",noindex" json:"artworkHash"`

	// The following two fields may be empty in responses to clients, if the
	// data hasn't been obtained for external sources (e.g., iTunes) yet.
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

func songKey(namespace, ident string, parent *datastore.Key) *datastore.Key {
	return &datastore.Key{Kind: KindSong, Name: ident, Parent: parent, Namespace: namespace}
}

type SongResponse struct {
	Song
	Ident string `json:"ident"`
}

const headerAPIKey = "X-Scrobble-API-Key"

func (s *server) accountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.Header.Get(headerAPIKey)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, _, ok := fetchAccountForKey(ctx, key, s.ds, w)
	if !ok {
		return
	}

	b, err := json.Marshal(struct {
		Username string `json:"username"`
	}{acc.Username})

	if err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (s *server) deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
		_, accID, ok = fetchAccountForKey(ctx, key, s.ds, w)
		if !ok {
			return
		}
	} else {
		accID = u.Email
	}

	tasksSecret, err := tasksSecret(ctx, s.ds)
	if err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof(ctx, "deleting account %q", accID)

	if _, err := s.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		// synchronously delete Username, Account entities
		// NOTE: we intentionally do not delete the APIKey entity because
		// those should always exist to guarantee non-reuse.
		aKey := &datastore.Key{Kind: KindAccount, Name: accID}
		var account Account
		if err := tx.Get(aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to get account")
		}

		log.Infof(ctx, "Account entity %+v for key %s", account, aKey)

		// If the account isn't initialized, the username won't be set
		// and a corresponding Username entity won't exist. So only
		// attempt to delete the Username entity if the username is
		// set.
		if account.Username != "" {
			if err := tx.Delete(&datastore.Key{Kind: KindUsername, Name: account.Username}); err != nil {
				return errors.Wrapf(err, "failed to delete Username entity")
			}
		}

		if err := tx.Delete(aKey); err != nil {
			return errors.Wrapf(err, "failed to delete Account entity")
		}

		// asynchronously delete the namespace's entities
		// (deletion order should not matter here)
		namespace := namespaceID(accID)

		{
			createReq, err := jsonPostTask("/internal/deleteEntities", deleteEntitiesTask{
				Namespace: namespace,
				Kind:      KindArtworkRecord,
			}, tasksSecret)
			if err != nil {
				return errors.Wrapf(err, "failed to build task")
			}
			if _, err := s.tasks.CreateTask(ctx, createReq); err != nil {
				return errors.Wrapf(err, "failed to add task for %s,%s", namespace, KindArtworkRecord)
			}
		}

		{
			createReq, err := jsonPostTask("/internal/deleteEntities", deleteEntitiesTask{
				Namespace: namespace,
				Kind:      KindSongParent,
			}, tasksSecret)
			if err != nil {
				return errors.Wrapf(err, "failed to build task")
			}
			if _, err := s.tasks.CreateTask(ctx, createReq); err != nil {
				return errors.Wrapf(err, "failed to add task for %s,%s", namespace, KindSongParent)
			}
		}

		{
			createReq, err := jsonPostTask("/internal/deleteEntities", deleteEntitiesTask{
				Namespace: namespace,
				Kind:      KindSong,
			}, tasksSecret)
			if err != nil {
				return errors.Wrapf(err, "failed to build task")
			}
			if _, err := s.tasks.CreateTask(ctx, createReq); err != nil {
				return errors.Wrapf(err, "failed to add task for %s,%s", namespace, KindSong)
			}
		}
		return nil
	}); err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (svr *server) scrobbledHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")

	writeSuccessRsp := func(s []SongResponse) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Errorf(ctx, "failed to write response: %v", err.Error())
		}
	}

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lovedOnly := r.FormValue("loved") == "true"
	songIdent := r.FormValue("song")

	if songIdent != "" && lovedOnly {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, accID, ok := fetchAccountForUsername(ctx, username, svr.ds, w)
	if !ok {
		return
	}

	if acc.Private && !canViewScrobbled(ctx, svr.ds, accID, user.Current(ctx), r.Header) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	namespace := namespaceID(accID)

	// Get the latest parent.
	q := datastore.NewQuery(KindSongParent).
		Namespace(namespace).
		Order("-Created").Filter("Complete=", true).
		Limit(1).KeysOnly()

	parentKeys, err := svr.ds.GetAll(ctx, q, nil)
	if err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(parentKeys) == 0 {
		// no songs, respond with empty JSON array
		writeSuccessRsp(make([]SongResponse, 0))
		return
	}

	if songIdent != "" {
		// Get song by ident.
		key := songKey(namespace, songIdent, parentKeys[0])
		var s SongResponse
		if err := svr.ds.Get(ctx, key, &s); err != nil {
			log.Errorf(ctx, "failed to fetch song %s: %v", key, err.Error())
			if err == datastore.ErrNoSuchEntity {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		s.Ident = key.Name
		writeSuccessRsp([]SongResponse{s})
	} else {
		// Get all songs.
		q = datastore.NewQuery(KindSong).
			Namespace(namespace).
			Order("-LastPlayed").
			Ancestor(parentKeys[0])

		if lovedOnly {
			q = q.Filter("Loved=", true)
		}

		songs := make([]SongResponse, 0) // use "make" to marshal as empty JSON array instead of null when there are 0 songs
		keys, err := svr.ds.GetAll(ctx, q, &songs)
		if err != nil {
			log.Errorf(ctx, "failed to fetch songs: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range songs {
			songs[i].Ident = keys[i].Name
		}
		writeSuccessRsp(songs)
	}
}

func canViewScrobbled(ctx context.Context, ds *datastore.Client, forAccountID string, u *user.User, h http.Header) bool {
	if u != nil && u.Email == forAccountID {
		// a logged in user can view their own account's scrobbles
		return true
	}

	if key := h.Get(headerAPIKey); key != "" {
		if _, id, _, err := accountForKey(ctx, key, ds); err == nil && id == forAccountID {
			return true
		}
	}

	return false
}

func (svr *server) scrobbleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.Header.Get(headerAPIKey)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, accID, ok := fetchAccountForKey(ctx, key, svr.ds, w)
	if !ok {
		return
	}

	namespace := namespaceID(accID)

	tasksSecret, err := tasksSecret(ctx, svr.ds)
	if err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
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
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now()

	// Convert to Songs.
	var songs []Song
	for _, m := range mis {
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

		// There is a bug (somewhere) that leads to some songs having last
		// played times far into the future, for instance, in the year 2040.
		// So ignore such songs.
		if last := time.Unix(int64(m.LastPlayed), 0); last.Sub(now) > 365*24*time.Hour {
			log.Warningf(ctx, "skipping song with future LastPlayed = %s", last)
			continue
		}

		songs = append(songs, Song{
			AlbumTitle:     m.AlbumTitle,
			ArtistName:     m.ArtistName,
			Title:          m.Title,
			TotalTime:      time.Duration(m.TotalTime) * (time.Millisecond / time.Nanosecond),
			Year:           int(m.Year),
			ReleaseDate:    int64(m.ReleaseDate),
			SortAlbumTitle: sal,
			SortArtistName: sar,
			SortTitle:      st,
			LastPlayed:     int64(m.LastPlayed),
			PlayCount:      int(m.PlayCount),
			ArtworkHash:    m.ArtworkHash,
			Loved:          m.Loved,
		})
	}

	newParentIdent := songparentident(now, uuid.New())
	newParentKey := songParentKey(namespace, newParentIdent)

	// Create new incomplete SongParent.
	if _, err := svr.ds.Put(ctx, newParentKey, &SongParent{
		Complete: false,
		Created:  now.Unix(),
	}); err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sKeys := make([]*datastore.Key, len(songs))
	var aKeys []*datastore.Key
	for i, s := range songs {
		// Create song key.
		sKeys[i] = songKey(namespace, s.Ident(), newParentKey)
		// Create artwork hash key.
		if h := s.ArtworkHash; h != "" {
			aKeys = append(aKeys, &datastore.Key{Kind: KindArtworkRecord, Name: h, Namespace: namespace})
		}
	}

	// Put songs.
	{
		s := 0
		e := min(s+datastoreLimitPerOp, len(songs))
		chunk := songs[s:e]
		keysChunk := sKeys[s:e]

		for len(chunk) > 0 {
			if _, err := svr.ds.PutMulti(ctx, keysChunk, chunk); err != nil {
				log.Errorf(ctx, "failed to put songs: %v", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			s = e
			e = min(s+datastoreLimitPerOp, len(songs))
			chunk = songs[s:e]
			keysChunk = sKeys[s:e]
		}
	}

	var g errgroup.Group // TODO: why is this group not deriving from request context?

	// Put artwork hash keys.
	g.Go(func() error {
		s := 0
		e := min(s+datastoreLimitPerOp, len(aKeys))
		keysChunk := aKeys[s:e]

		for len(keysChunk) > 0 {
			if _, err := svr.ds.PutMulti(ctx, keysChunk, make([]struct{}, len(keysChunk))); err != nil {
				return errors.Wrapf(err, "failed to put artwork records")
			}

			s = e
			e = min(s+datastoreLimitPerOp, len(aKeys))
			keysChunk = aKeys[s:e]
		}

		return nil
	})

	// Create tasks to fill in iTunes-related fields.
	for _, s := range songs {
		songIdent := s.Ident()
		g.Go(func() error {
			createReq, err := jsonPostTask("/internal/fillITunesFields", fillITunesFieldsTask{
				Namespace:       namespace,
				SongParentIdent: newParentIdent,
				SongIdent:       songIdent,
			}, tasksSecret)
			if err != nil {
				return errors.Wrapf(err, "failed to build task")
			}
			if _, err := svr.tasks.CreateTask(ctx, createReq); err != nil {
				return errors.Wrapf(err, "failed to add fillITunesFields tasks for %s,%s,%s", namespaceID(accID), newParentIdent, songIdent)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// We can mark the parent as complete, now that the songs have all been put.
	// But we want the iTunes related fields to also be filled in (which happens
	// asynchronously with staggered internal delaying), so wait for a bit
	// before marking the parent as complete.
	//
	// TODO: make this deterministic instead of using a delay?
	if err := func() error {
		payload := markParentCompleteTask{
			Namespace:       namespaceID(accID),
			SongParentIdent: newParentIdent,
		}
		createReq, err := jsonPostTask("/internal/markParentComplete", payload, tasksSecret)
		if err != nil {
			return errors.Wrapf(err, "failed to build task")
		}
		createReq.Task.ScheduleTime, _ = ptypes.TimestampProto(time.Now().Add(2 * time.Minute))

		if _, err := svr.tasks.CreateTask(ctx, createReq); err != nil {
			return errors.Wrapf(err, "failed to add task")
		}
		return nil
	}(); err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) artworkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
	_, _, ok := fetchAccountForKey(ctx, key, s.ds, w)
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
	wr := s.storage.Bucket(DefaultBucketName).Object(artworkStorageDirectory + "/" + hash).NewWriter(ctx)
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

func (s *server) artworkMissingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.Header.Get(headerAPIKey)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, accID, ok := fetchAccountForKey(ctx, key, s.ds, w)
	if !ok {
		return
	}

	namespace := namespaceID(accID)

	g, gctx := errgroup.WithContext(ctx)
	have := make(map[string]struct{})
	want := make(map[string]bool)
	var artworkKeys []*datastore.Key

	g.Go(func() error {
		it := s.storage.Bucket(DefaultBucketName).Objects(gctx, &storage.Query{Prefix: artworkStorageDirectory + "/"})
		for {
			o, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Errorf(gctx, "%v", err.Error()) // only log
				break
			}
			have[strings.TrimPrefix(o.Name, artworkStorageDirectory+"/")] = struct{}{}
		}
		return nil
	})

	g.Go(func() error {
		q := datastore.NewQuery(KindArtworkRecord).Namespace(namespace).KeysOnly()
		keys, err := s.ds.GetAll(gctx, q, nil)
		if err != nil {
			return errors.Wrapf(err, "failed to fetch artwork records")
		}
		artworkKeys = keys
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, k := range artworkKeys {
		if _, ok := have[k.Name]; ok {
			continue
		}
		want[k.Name] = true
	}

	log.Infof(ctx, "%d artwork records with missing data", len(want))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(want); err != nil {
		log.Errorf(ctx, "failed to write response: %v", err.Error())
	}
}

type markParentCompleteTask struct {
	Namespace       string
	SongParentIdent string
}

func (s *server) markParentCompleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var t markParentCompleteTask
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.markParentComplete(ctx, t.Namespace, t.SongParentIdent); err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Marks a SongParent as complete, and deletes really old SongParents and their
// child Songs.
func (s *server) markParentComplete(ctx context.Context, namespace string, songParentIdent string) error {
	// Create parent link.
	newParentKey := songParentKey(namespace, songParentIdent)
	if _, err := s.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var sp SongParent
		if err := tx.Get(newParentKey, &sp); err != nil {
			return errors.Wrapf(err, "failed to get %s", newParentKey)
		}
		sp.Complete = true
		if _, err := tx.Put(newParentKey, &sp); err != nil {
			return errors.Wrapf(err, "failed to put %s", newParentKey)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := trimSongParents(ctx, namespace, s.ds); err != nil {
		return err
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func accountForKey(ctx context.Context, apiKey string, ds *datastore.Client) (Account, string, int, error) {
	q := datastore.NewQuery(KindAccount).Filter("APIKey=", apiKey).Limit(2)

	var as []Account
	keys, err := ds.GetAll(ctx, q, &as)
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

	return as[0], keys[0].Name, 0, nil
}

func accountForUsername(ctx context.Context, username string, ds *datastore.Client) (Account, string, int, error) {
	q := datastore.NewQuery(KindAccount).Filter("Username=", username).Limit(2)

	var as []Account
	keys, err := ds.GetAll(ctx, q, &as)
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

	return as[0], keys[0].Name, 0, nil
}

func fetchAccountForKey(ctx context.Context, apiKey string, ds *datastore.Client, w http.ResponseWriter) (Account, string, bool) {
	a, id, code, err := accountForKey(ctx, apiKey, ds)
	if err != nil {
		log.Errorf(ctx, err.Error())
		w.WriteHeader(code)
		return Account{}, "", false
	}
	return a, id, true
}

func fetchAccountForUsername(ctx context.Context, username string, ds *datastore.Client, w http.ResponseWriter) (Account, string, bool) {
	a, id, code, err := accountForUsername(ctx, username, ds)
	if err != nil {
		log.Errorf(ctx, err.Error())
		w.WriteHeader(code)
		return Account{}, "", false
	}
	return a, id, true
}

const headerTasksSecret = "X-Scrobble-Tasks-Secret"

func tasksSecret(ctx context.Context, ds *datastore.Client) (string, error) {
	var secret Secret
	if err := ds.Get(ctx, datastore.NameKey(KindSecret, "singleton", nil), &secret); err != nil {
		return "", errors.Wrapf(err, "failed to get from datastore")
	}
	if secret.TasksSecret == "" {
		panic("empty TasksSecret")
	}
	return secret.TasksSecret, nil
}

func (s *server) requireTasksSecretHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		got := r.Header.Get(headerTasksSecret)
		if got == "" {
			log.Errorf(ctx, "missing tasks secret header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		want, err := tasksSecret(ctx, s.ds)
		if err != nil {
			log.Errorf(ctx, "tasks secret: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if want != got {
			log.Errorf(ctx, "bad tasks secret header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func jsonPostTask(path string, payload interface{}, secret string) (*tasks.CreateTaskRequest, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to json-marshal payload")
	}

	return &tasks.CreateTaskRequest{
		Parent: DefaultQueueName,
		Task: &tasks.Task{
			PayloadType: &tasks.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &tasks.AppEngineHttpRequest{
					HttpMethod:  tasks.HttpMethod_POST,
					RelativeUrl: path,
					Headers:     map[string]string{headerTasksSecret: secret},
					Payload:     p,
				},
			},
		},
	}, nil
}

func namespaceID(accountID string) string {
	// allowed namespace pattern: /^[0-9A-Za-z._-]{0,100}$/
	return string(hexencode([]byte(accountID)))
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
