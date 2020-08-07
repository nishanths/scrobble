package main

import (
	"encoding/json"
	"io"
	"net/http"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/nishanths/scrobble/appengine/log"
)

func (s *server) songPlayCountHandler(w http.ResponseWriter, r *http.Request) {
	s.songDataHandler("-PlayCount").ServeHTTP(w, r)
}

func (s *server) songLengthHandler(w http.ResponseWriter, r *http.Request) {
	s.songDataHandler("-TotalTime").ServeHTTP(w, r)
}

func (s *server) songDataHandler(fieldName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeSuccessRsp := func(s []SongResponse) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(s); err != nil {
				log.Errorf("failed to write response: %v", err.Error())
			}
		}

		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()

		username := r.FormValue("username")
		if username == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		limit, hasLimit := parseLimit(r.FormValue("limit"))

		acc, accID, ok := fetchAccountForUsername(ctx, username, s.ds, w)
		if !ok {
			return
		}

		if acc.Private && !s.canViewScrobbled(ctx, accID, r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		namespace := namespaceID(accID)

		// get latest (complete) song parent
		q := datastore.NewQuery(KindSongParent).
			Namespace(namespace).
			Order("-Created").Filter("Complete=", true).
			Limit(1).KeysOnly()

		parentKeys, err := s.ds.GetAll(ctx, q, nil)
		if err != nil {
			log.Errorf("%v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(parentKeys) == 0 {
			writeSuccessRsp(make([]SongResponse, 0))
			return
		}

		// get songs, sorted by provided field
		q = datastore.NewQuery(KindSong).
			Namespace(namespace).
			Order(fieldName).
			Ancestor(parentKeys[0])
		if hasLimit {
			q = q.Limit(limit)
		}

		songs := make([]SongResponse, 0) // "make" to json-marshal as empty array instead of null when there are 0 songs
		keys, err := s.ds.GetAll(ctx, q, &songs)
		if err != nil {
			log.Errorf("failed to fetch songs: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// fill in ident
		for i := range songs {
			songs[i].Ident = keys[i].Name
		}

		writeSuccessRsp(songs)
	})
}

func (s *server) artistPlayCountHandler(w http.ResponseWriter, r *http.Request) {
	s.artistDataHandler(statsArtistPlayCountStoragePath).ServeHTTP(w, r)
}

func (s *server) artistAddedHandler(w http.ResponseWriter, r *http.Request) {
	s.artistDataHandler(statsArtistAddedStoragePath).ServeHTTP(w, r)
}

func (s *server) artistDataHandler(storagePathFunc func(namespace string) string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()

		username := r.FormValue("username")
		if username == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		acc, accID, ok := fetchAccountForUsername(ctx, username, s.ds, w)
		if !ok {
			return
		}

		if acc.Private && !s.canViewScrobbled(ctx, accID, r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		namespace := namespaceID(accID)

		rd, err := s.storage.Bucket(DefaultBucketName).Object(storagePathFunc(namespace)).NewReader(ctx)
		if err == storage.ErrObjectNotExist {
			log.Infof("no artist stats for %s", storagePathFunc(namespace))
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if err != nil {
			log.Errorf("failed to get: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := io.Copy(w, rd); err != nil {
			log.Errorf("failed to write response: %v", err.Error())
			return
		}
		if err := rd.Close(); err != nil {
			log.Errorf("failed to close storage reader: %v", err.Error())
		}
	})
}
