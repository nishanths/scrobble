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
		Filter(fmt.Sprintf("Score.%s >=", datastoreFieldNameForColor(inputColor)), 400).
		KeysOnly()
	if hasLimit {
		q = q.Limit(limit)
	}

	keys, err := s.ds.GetAll(ctx, q, nil)
	if err != nil {
		log.Errorf("failed to get query: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prepare response.
	hashes := make([]string, len(keys))
	for i, k := range keys {
		hashes[i] = k.Name
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hashes); err != nil {
		log.Errorf("failed to write response: %v", err.Error())
	}
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
