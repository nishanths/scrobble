package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/nishanths/scrobble/appengine/artwork"
	"github.com/nishanths/scrobble/appengine/basiccolor"
	"github.com/nishanths/scrobble/appengine/log"
	"github.com/pkg/errors"
)

var validColors = [...]basiccolor.Color{
	basiccolor.Red,
	basiccolor.Orange,
	basiccolor.Brown,
	basiccolor.Yellow,
	basiccolor.Green,
	basiccolor.Blue,
	basiccolor.Violet,
	basiccolor.Pink,
	basiccolor.Black,
	basiccolor.Gray,
	basiccolor.White,
}

func (s *server) artworkColorHandler(w http.ResponseWriter, r *http.Request) {
	writeSuccessRsp := func(s []SongResponse) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Errorf("failed to write response: %v", err.Error())
		}
	}

	ctx := r.Context()

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	limit, hasLimit := parseLimit(r.FormValue("limit"))

	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse input color.
	c := r.FormValue("color")
	if c == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inputColor basiccolor.Color
	var ok bool
	for _, b := range validColors {
		if strings.EqualFold(b.String(), c) {
			inputColor = b
			ok = true
			break
		}
	}
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, accID, ok := fetchAccountForUsername(ctx, username, s.ds, w)
	if !ok {
		return
	}

	if acc.Private && !s.canViewScrobbled(ctx, accID, r) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	namespace := namespaceID(accID)

	// query datastore for ArtworkRecords matching the input color.
	q := datastore.NewQuery(artwork.KindArtworkRecord).
		Namespace(namespace).
		Order(fmt.Sprintf("-Score.%s", datastoreFieldNameForColor(inputColor))).
		Filter(fmt.Sprintf("Score.%s >=", datastoreFieldNameForColor(inputColor)), 300).
		Project("SongIdent")

	if hasLimit {
		q = q.Limit(limit)
	}

	var artworkRecords []artwork.ArtworkRecord
	_, err := s.ds.GetAll(ctx, q, &artworkRecords)
	if err != nil {
		log.Errorf("failed to get query: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the latest complete parent (required to construct Song keys
	// for fetching).
	q = datastore.NewQuery(KindSongParent).
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
		// WTF?
		log.Warningf("unexpectedly found 0 SongParent keys")
		writeSuccessRsp(make([]SongResponse, 0))
		return
	}

	// construct SongKeys.
	songKeys := make([]*datastore.Key, len(artworkRecords))
	for i, ar := range artworkRecords {
		songKeys[i] = songKey(namespace, ar.SongIdent, parentKeys[0])
	}

	// Fetch songs by SongKeys.
	// Some of the songs may not be found if the artwork is from a newer generation.
	// Ignore such errors.
	songs := make([]Song, len(songKeys))
	err = s.ds.GetMulti(ctx, songKeys, songs)
	if err != nil && !isAllErrNoSuchEntity(err) {
		log.Errorf("failed to songs: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// omit items that correspond to ErrNoSuchEntity.
	var foundSongs []Song
	if err != nil {
		merr := err.(datastore.MultiError)
		for i, e := range merr {
			if e == nil {
				foundSongs = append(foundSongs, songs[i])
			}
		}
	} else {
		foundSongs = songs
	}

	// construct response
	rsp := make([]SongResponse, len(foundSongs))
	for i, s := range foundSongs {
		rsp[i] = SongResponse{
			Song:  s,
			Ident: songident(s.AlbumTitle, s.ArtistName, s.Title, s.Year),
		}
	}

	writeSuccessRsp(rsp)
}

func isAllErrNoSuchEntity(err error) bool {
	if err == nil {
		panic("err must be non-nil")
	}
	merr, ok := err.(datastore.MultiError)
	if !ok {
		return err == datastore.ErrNoSuchEntity
	}
	for _, e := range merr {
		if e == nil || e == datastore.ErrNoSuchEntity {
			continue
		}
		return false
	}
	return true
}

func datastoreFieldNameForColor(c basiccolor.Color) string {
	return c.String()
}

type fillArtworkScoreTask struct {
	Namespace string
	Hash      string
}

func (s *server) fillArtworkScoreHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var t fillArtworkScoreTask
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Errorf("failed to json-decode task: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Attempt to fill in artwork scores from the global repository of artwork scores.
	var as artwork.ArtworkScore
	err := s.ds.Get(ctx, artwork.ArtworkScoreKey(t.Hash), &as)
	if err == datastore.ErrNoSuchEntity {
		w.WriteHeader(http.StatusOK)
		return
	}
	if err != nil {
		log.Errorf("failed datastore get: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := s.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var ar artwork.ArtworkRecord
		arKey := artwork.ArtworkRecordKey(t.Namespace, t.Hash)
		err := tx.Get(arKey, &ar)
		if err != nil {
			return errors.Wrapf(err, "failed datastore get")
		}

		ar.Score = as

		_, err = tx.Put(arKey, &ar)
		if err != nil {
			return errors.Wrapf(err, "failed datastore put")
		}
		return nil
	}); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
