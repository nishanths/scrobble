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
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/RobCherry/vibrant"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/nishanths/scrobble/appengine/artwork"
	"github.com/nishanths/scrobble/appengine/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
	tasks "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

const (
	artworkStorageDirectory = "aw" // Cloud Storage directory name for artwork
)

const (
	KindAccount     = "Account"     // namespace: [default]
	KindUsername    = "Username"    // namespace: [default]; stored for uniqueness checking
	KindAPIKey      = "APIKey"      // namespace: [default]; stored for uniqueness checking
	KindITunesTrack = "ITunesTrack" // namespace: [default]
	KindSecret      = "Secret"      // namespace: [default]

	KindSongParent = "SongParent" // namespace: Account
	KindSong       = "Song"       // namespace: Account
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
	PlayCount  int   `json:"playCount"`

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
		log.Errorf("%v", err.Error())
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
	if u, err := s.currentUser(r); err != nil {
		key := r.Header.Get(headerAPIKey)
		if key == "" {
			log.Errorf("not signed-in and missing API key header")
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

	log.Infof("deleting account %q", accID)

	if _, err := s.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		// synchronously delete Username, Account entities
		// NOTE: we intentionally do not delete the APIKey entity because
		// those should always exist to guarantee non-reuse.
		aKey := &datastore.Key{Kind: KindAccount, Name: accID}
		var account Account
		if err := tx.Get(aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to get account")
		}

		log.Infof("Account entity %+v for key %s", account, aKey)

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
		deleteKinds := []string{artwork.KindArtworkRecord, KindSongParent, KindSong}

		for _, k := range deleteKinds {
			createReq, err := jsonPostTask("/internal/deleteEntities", deleteEntitiesTask{
				Namespace: namespace,
				Kind:      k,
			}, s.secret.TasksSecret)
			if err != nil {
				return errors.Wrapf(err, "failed to build task")
			}
			if _, err := s.tasks.CreateTask(ctx, createReq); err != nil {
				return errors.Wrapf(err, "failed to add task for %s,%s", namespace, k)
			}
		}

		return nil
	}); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func parseLimit(lim string) (int, bool) {
	i, err := strconv.Atoi(lim)
	if err != nil {
		return -1, false
	}
	if i < 0 {
		return -1, false
	}
	return i, true
}

func clampLimit(lim int, hasLimit bool, max int) int {
	if !hasLimit || lim > max {
		return max
	}
	return lim
}

type ScrobbledResponse struct {
	// Total is the total number of scrobbled songs.
	// It is valid only if the request wasn't for a specific song (i.e., "song"
	// query parameter wasn't set).
	Total int            `json:"total"`
	Songs []SongResponse `json:"songs"`
}

func (svr *server) scrobbledHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")

	writeSuccessRsp := func(s ScrobbledResponse) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Errorf("failed to write response: %v", err.Error())
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
	limit, hasLimit := parseLimit(r.FormValue("limit"))

	// songIdent must be mutually exclusive with loved.
	// Also, songIdent must be mutually exclusive with limit.
	if songIdent != "" && (lovedOnly || hasLimit) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, accID, ok := fetchAccountForUsername(ctx, username, svr.ds, w)
	if !ok {
		return
	}

	if acc.Private && !svr.canViewScrobbled(ctx, accID, r) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	namespace := namespaceID(accID)

	// Get the latest complete parent.
	q := datastore.NewQuery(KindSongParent).
		Namespace(namespace).
		Order("-Created").Filter("Complete=", true).
		Limit(1).KeysOnly()

	parentKeys, err := svr.ds.GetAll(ctx, q, nil)
	if err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(parentKeys) == 0 {
		// no songs
		writeSuccessRsp(ScrobbledResponse{
			Total: 0,
			Songs: make([]SongResponse, 0),
		})
		return
	}

	if songIdent != "" {
		// Get song by ident.
		key := songKey(namespace, songIdent, parentKeys[0])
		var s SongResponse
		err := svr.ds.Get(ctx, key, &s)
		if err == datastore.ErrNoSuchEntity {
			// 200 with empty songs list
			writeSuccessRsp(ScrobbledResponse{
				Total: -1,
				Songs: make([]SongResponse, 0),
			})
			return
		}
		if err != nil {
			log.Errorf("failed to fetch song %s: %v", key, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.Ident = key.Name
		writeSuccessRsp(ScrobbledResponse{
			Total: -1,
			Songs: []SongResponse{s},
		})
	} else {
		var songs []SongResponse
		var total int
		g, gctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			// Get all songs.
			q = datastore.NewQuery(KindSong).
				Namespace(namespace).
				Order("-LastPlayed").
				Ancestor(parentKeys[0])

			if lovedOnly {
				q = q.Filter("Loved=", true)
			}

			if hasLimit {
				q = q.Limit(limit)
			}

			songs = make([]SongResponse, 0) // "make" to json-marshal as empty array instead of null when there are 0 songs
			keys, err := svr.ds.GetAll(gctx, q, &songs)
			if err != nil {
				return errors.Wrapf(err, "failed to fetch songs")
			}
			for i := range songs {
				songs[i].Ident = keys[i].Name
			}

			if !hasLimit {
				total = len(songs)
			}

			return nil
		})

		if hasLimit {
			// need to compute total count with a separate query
			g.Go(func() error {
				q = datastore.NewQuery(KindSong).
					Namespace(namespace).
					Ancestor(parentKeys[0])

				if lovedOnly {
					q = q.Filter("Loved=", true)
				}

				n, err := svr.ds.Count(gctx, q)
				if err != nil {
					return errors.Wrapf(err, "failed to count total songs")
				}
				total = n
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			log.Errorf("%v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeSuccessRsp(ScrobbledResponse{
			Total: total,
			Songs: songs,
		})
	}
}

func (s *server) canViewScrobbled(ctx context.Context, forAccountID string, r *http.Request) bool {
	u, err := s.currentUser(r)
	if err == nil && u.Email == forAccountID {
		// a logged in user can view their own account's scrobbles
		return true
	}

	if key := r.Header.Get(headerAPIKey); key != "" {
		if _, id, _, err := accountForKey(ctx, key, s.ds); err == nil && id == forAccountID {
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
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now()

	// Convert to Songs (along with de-duplication)
	var songs []Song
	songIdentsSet := make(map[string]struct{}) // for de-duplication
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
			log.Warningf("skipping song with future LastPlayed = %s", last)
			continue
		}

		s := Song{
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
		}

		// Handle de-duplication
		if _, ok := songIdentsSet[s.Ident()]; ok {
			// NOTE: This happens in practice; macOS client seems bad?
			log.Warningf("duplicate incoming song ident %v; skipping", s.Ident())
			continue
		}
		songIdentsSet[s.Ident()] = struct{}{}

		songs = append(songs, s)
	}

	// Create new incomplete SongParent.
	newParentIdent := songparentident(now, uuid.New())
	newParentKey := songParentKey(namespace, newParentIdent)
	newParentCreated := now.Unix()

	if _, err := svr.ds.Put(ctx, newParentKey, &SongParent{
		Complete: false,
		Created:  newParentCreated,
	}); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var sKeys []*datastore.Key
	incomingArtworkHashes := make(map[string]string)
	for _, s := range songs {
		// Create song key.
		sKeys = append(sKeys, songKey(namespace, s.Ident(), newParentKey))
		// Collect incoming artwork hash.
		if h := s.ArtworkHash; h != "" {
			// the same artwork hash can be there for >1 songs; any one of those
			// songs' idents will do here
			incomingArtworkHashes[h] = s.Ident()
		}
	}

	// Put songs.
	// (This should serially succeed before we proceed further.)
	{
		s := 0
		e := min(s+datastoreLimitPerOp, len(songs))
		chunk := songs[s:e]
		keysChunk := sKeys[s:e]

		for len(chunk) > 0 {
			if _, err := svr.ds.PutMulti(ctx, keysChunk, chunk); err != nil {
				log.Errorf("failed to put songs: %v", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			s = e
			e = min(s+datastoreLimitPerOp, len(songs))
			chunk = songs[s:e]
			keysChunk = sKeys[s:e]
		}
	}

	// Fetch currently existing artwork record hashes.
	q := datastore.NewQuery(artwork.KindArtworkRecord).Namespace(namespace).KeysOnly()
	currentArtworkKeys, err := svr.ds.GetAll(ctx, q, nil)
	if err != nil {
		log.Errorf("failed to fetch artwork records: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	currentArtworkHashes := make(map[string]struct{}, len(currentArtworkKeys))
	for _, k := range currentArtworkKeys {
		currentArtworkHashes[k.Name] = struct{}{}
	}

	// Diff current and incoming artwork hashes.
	addHashes, removeHashes := diffStringMaps(currentArtworkHashes, incomingArtworkHashes)

	var g errgroup.Group

	// Adjust artwork records based on added / removed hashes.
	g.Go(func() error {
		keys := make([]*datastore.Key, 0, len(addHashes))
		entities := make([]artwork.ArtworkRecord, 0, len(addHashes))
		for k, v := range addHashes {
			keys = append(keys, artwork.ArtworkRecordKey(namespace, k))
			entities = append(entities, artwork.ArtworkRecord{SongIdent: v})
		}
		_, err := svr.ds.PutMulti(ctx, keys, entities)
		return err
	})
	g.Go(func() error {
		keys := make([]*datastore.Key, 0, len(removeHashes))
		for k := range removeHashes {
			keys = append(keys, artwork.ArtworkRecordKey(namespace, k))
		}
		return svr.ds.DeleteMulti(ctx, keys)
	})

	// Create artwork score fill-in tasks for the newly added hashes.
	for k := range addHashes {
		g.Go(func() error {
			createReq, err := jsonPostTask("/internal/fillArtworkScore", fillArtworkScoreTask{
				Namespace: namespace,
				Hash:      k,
			}, svr.secret.TasksSecret)
			if err != nil {
				return errors.Wrapf(err, "failed to build task")
			}
			if _, err := svr.tasks.CreateTask(ctx, createReq); err != nil {
				return errors.Wrapf(err, "failed to add fillArtworkScore tasks for %s,%s", namespaceID(accID), k)
			}
			return nil
		})
	}

	// Create tasks to fill in iTunes-related fields.
	for _, s := range songs {
		songIdent := s.Ident()
		g.Go(func() error {
			createReq, err := jsonPostTask("/internal/fillITunesFields", fillITunesFieldsTask{
				Namespace:       namespace,
				SongParentIdent: newParentIdent,
				SongIdent:       songIdent,
			}, svr.secret.TasksSecret)
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
		log.Errorf("%v", err.Error())
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
			Namespace:         namespaceID(accID),
			SongParentIdent:   newParentIdent,
			SongParentCreated: newParentCreated,
		}
		createReq, err := jsonPostTask("/internal/markParentComplete", payload, svr.secret.TasksSecret)
		if err != nil {
			return errors.Wrapf(err, "failed to build task")
		}
		createReq.Task.ScheduleTime, _ = ptypes.TimestampProto(time.Now().Add(2 * time.Minute))

		if _, err := svr.tasks.CreateTask(ctx, createReq); err != nil {
			return errors.Wrapf(err, "failed to add task")
		}
		return nil
	}(); err != nil {
		log.Errorf("%v", err.Error())
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

	// validate API key, get account ID
	_, accID, ok := fetchAccountForKey(ctx, key, s.ds, w)
	if !ok {
		return
	}

	namespace := namespaceID(accID)

	artworkBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("failed to read request body: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	format := r.FormValue("format")
	hash := artworkHash(artworkBytes, format)

	// upload to GCS
	wr := s.storage.Bucket(DefaultBucketName).Object(artworkStorageDirectory + "/" + hash).NewWriter(ctx)
	wr.Metadata = map[string]string{"format": format}
	if _, err := wr.Write(artworkBytes); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := wr.Close(); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// compute artwork score
	img, _, err := image.Decode(bytes.NewReader(artworkBytes))
	if err != nil {
		log.Errorf("failed to decode image: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pal := vibrant.NewPaletteBuilder(img).Generate()
	swatches := artwork.Swatches(pal.Swatches())
	artworkScore := artwork.ArtworkScoreFromSwatches(swatches)

	// Save score to global and namespaced entities.
	if _, err := s.ds.Put(ctx, artwork.ArtworkScoreKey(hash), &artworkScore); err != nil {
		log.Errorf("failed datastore put: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := s.ds.Put(ctx, artwork.ArtworkRecordKey(namespace, hash), &artwork.ArtworkRecord{
		Score: artworkScore,
	}); err != nil {
		log.Errorf("failed datastore put: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("saved artwork hash=%s", hash)

	hashJson, err := json.Marshal(hash)
	if err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(hashJson)
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
				log.Errorf("failed to iterate artwork bucket: %v", err.Error()) // only log
				break
			}
			have[strings.TrimPrefix(o.Name, artworkStorageDirectory+"/")] = struct{}{}
		}
		return nil
	})

	g.Go(func() error {
		q := datastore.NewQuery(artwork.KindArtworkRecord).Namespace(namespace).KeysOnly()
		keys, err := s.ds.GetAll(gctx, q, nil)
		if err != nil {
			return errors.Wrapf(err, "failed to fetch artwork records")
		}
		artworkKeys = keys
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, k := range artworkKeys {
		if _, ok := have[k.Name]; ok {
			continue
		}
		want[k.Name] = true
	}

	log.Infof("missing %d artwork images", len(want))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(want); err != nil {
		log.Errorf("failed to write response: %v", err.Error())
	}
}

type markParentCompleteTask struct {
	Namespace         string
	SongParentIdent   string
	SongParentCreated int64
}

func (s *server) markParentCompleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var t markParentCompleteTask
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.markParentComplete(ctx, t.Namespace, t.SongParentIdent, t.SongParentCreated); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Marks a SongParent as complete, and deletes any other SongParents and their
// child Songs.
func (s *server) markParentComplete(ctx context.Context, namespace, songParentIdent string, songParentCreated int64) error {
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

	if err := trimSongParents(ctx, namespace, songParentCreated, s.ds); err != nil {
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
		log.Criticalf(m)
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
		log.Criticalf(m)
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
		log.Errorf(err.Error())
		w.WriteHeader(code)
		return Account{}, "", false
	}
	return a, id, true
}

func fetchAccountForUsername(ctx context.Context, username string, ds *datastore.Client, w http.ResponseWriter) (Account, string, bool) {
	a, id, code, err := accountForUsername(ctx, username, ds)
	if err != nil {
		log.Errorf(err.Error())
		w.WriteHeader(code)
		return Account{}, "", false
	}
	return a, id, true
}

func jsonPostTask(path string, payload interface{}, secret string) (*tasks.CreateTaskRequest, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to json-marshal payload")
	}

	return &tasks.CreateTaskRequest{
		Parent: DefaultQueueName,
		Task: &tasks.Task{
			MessageType: &tasks.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &tasks.AppEngineHttpRequest{
					HttpMethod:  tasks.HttpMethod_POST,
					RelativeUri: path,
					Headers:     map[string]string{headerTasksSecret: secret},
					Body:        p,
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

// Diffs the given sets and returns the elements being added and the elements
// being removed.
func diffStringMaps(old map[string]struct{}, new map[string]string) (added map[string]string, removed map[string]struct{}) {
	added = make(map[string]string)
	removed = make(map[string]struct{})

	for k, v := range new {
		if _, ok := old[k]; !ok {
			added[k] = v // not present in old, so being newly added
		}
	}

	for k := range old {
		if _, ok := new[k]; !ok {
			removed[k] = struct{}{} // not present in new, so being removed
		}
	}

	return
}
