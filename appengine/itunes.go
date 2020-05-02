package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/nishanths/scrobble/appengine/log"
	"github.com/pkg/errors"
)

type fillITunesFieldsTask struct {
	Namespace       string
	SongParentIdent string
	SongIdent       string
}

func (s *server) fillITunesFieldsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var t fillITunesFieldsTask
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.fillITunesFields(ctx, t.Namespace, t.SongParentIdent, t.SongIdent); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (svr *server) fillITunesFields(ctx context.Context, namespace string, songParentIdent string, songIdent string) error {
	done := false
	sKey := songKey(namespace, songIdent, songParentKey(namespace, songParentIdent))

	// Only fill from existing data.
	if _, err := svr.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var s Song
		if err := tx.Get(sKey, &s); err != nil {
			if err == datastore.ErrNoSuchEntity {
				done = true
				return nil // deleted sometime in between?
			}
			return errors.Wrapf(err, "failed to get song %s", sKey)
		}
		if s.iTunesFilled() {
			done = true
			return nil
		}

		var track ITunesTrack
		trackKey := &datastore.Key{Kind: KindITunesTrack, Name: songIdent}
		if err := tx.Get(trackKey, &track); err != nil && err != datastore.ErrNoSuchEntity {
			return errors.Wrapf(err, "failed to get track %s", trackKey)
		} else if err == nil {
			// update the song and we are done
			s.PreviewURL = track.PreviewURL
			s.TrackViewURL = track.TrackViewURL
			if _, err := tx.Put(sKey, &s); err != nil {
				return errors.Wrapf(err, "failed to put song %s", sKey)
			}
			done = true
			return nil
		}

		// else err == datastore.ErrNoSuchEntity (will fetch from iTunes API below)
		return nil
	}); err != nil {
		log.Errorf("%v", err.Error())
		return err
	}

	if done {
		return nil
	}

	// a barebones attempt at staggering
	time.Sleep(time.Duration(rand.Intn(60e3)) * time.Millisecond)

	// Try both.
	if _, err := svr.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var s Song
		if err := tx.Get(sKey, &s); err != nil {
			if err == datastore.ErrNoSuchEntity {
				return nil // deleted sometime in between?
			}
			return errors.Wrapf(err, "failed to get song %s", sKey)
		}
		if s.iTunesFilled() {
			return nil
		}

		// Fetch data from iTunes API
		tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		searchTerm := strings.Join([]string{s.Title, s.ArtistName}, " ") // including the album name produces poorer results
		tracks, retry, err := iTunesSearchSong(tctx, svr.httpClient, searchTerm)
		if err != nil {
			log.Errorf("failed to search iTunes for %q: %s (retry=%v)", searchTerm, err, retry)
			if retry {
				return err // returning a non-nil error causes the task to retry
			}
			return nil // don't bother
		}

		matchingIdx := -1
		for idx := range tracks {
			if !tracks[idx].ReleaseDate.IsZero() && tracks[idx].Ident() == songIdent && tracks[idx].HasRequiredURLs() {
				matchingIdx = idx
				break
			}
		}

		if matchingIdx == -1 {
			log.Infof("no matching tracks found for %q (len=%d)", searchTerm, len(tracks))
			return nil
		}

		// store track for use in future calls
		trackKey := &datastore.Key{Kind: KindITunesTrack, Name: songIdent}
		if _, err := tx.Put(trackKey, &tracks[matchingIdx]); err != nil {
			return errors.Wrapf(err, "failed to put track %s", trackKey)
		}

		s.PreviewURL = tracks[matchingIdx].PreviewURL
		s.TrackViewURL = tracks[matchingIdx].TrackViewURL
		if _, err := tx.Put(sKey, &s); err != nil {
			return errors.Wrapf(err, "failed to put song %s", sKey)
		}
		return nil
	}); err != nil {
		log.Errorf("%v", err.Error())
		return err
	}

	return nil
}

type ITunesTrack struct {
	ArtistName     string    `datastore:",noindex"`
	TrackName      string    `datastore:",noindex"`
	CollectionName string    `datastore:",noindex"`
	TrackViewURL   string    `datastore:",noindex"`
	PreviewURL     string    `datastore:",noindex"`
	ReleaseDate    time.Time `datastore:",noindex"`
}

func (i *ITunesTrack) HasRequiredURLs() bool {
	return i.PreviewURL != "" && i.PreviewURL != "—" &&
		i.TrackViewURL != "" && i.TrackViewURL != "—"
}

// NOTE: track ident == song ident, by definition
func (i *ITunesTrack) Ident() string {
	if i.ReleaseDate.IsZero() {
		panic("zero ReleaseDate")
	}
	return songident(i.CollectionName, i.ArtistName, i.TrackName, i.ReleaseDate.Year())
}

// See https://affiliate.itunes.apple.com/resources/documentation/itunes-store-web-service-search-api/
func iTunesSearchSong(ctx context.Context, httpClient *http.Client, searchTerm string) ([]ITunesTrack, bool, error) {
	vals := make(url.Values)
	vals.Set("term", searchTerm)
	vals.Set("media", "music")
	vals.Set("entity", "song")
	vals.Set("limit", "50")
	vals.Set("version", "2")
	vals.Set("explicit", "Yes")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://itunes.apple.com/search?%s", vals.Encode()), nil)
	if err != nil {
		return nil, true, errors.Wrapf(err, "failed to build itunes request")
	}
	req = req.WithContext(ctx)
	rsp, err := httpClient.Do(req)
	if err != nil {
		// assume all errors are transient, indicate to retry
		return nil, true, errors.Wrapf(err, "failed to do itunes request")
	}
	defer drainAndClose(rsp.Body)

	if rsp.StatusCode != http.StatusOK {
		// indicate retry on 403s, which appears to be the "too many requests" code
		return nil, rsp.StatusCode == http.StatusForbidden, fmt.Errorf("bad status code %d", rsp.StatusCode)
	}

	type auxITunesTrack struct {
		ArtistName     string `json:"artistName"`
		TrackName      string `json:"trackName"`
		CollectionName string `json:"collectionName"`
		TrackViewURL   string `json:"trackViewUrl"`
		PreviewURL     string `json:"previewUrl"`
		ReleaseDate    string `json:"releaseDate"`
	}
	type iTunesSearchResponse struct {
		ResultCount int              `json:"resultCount"`
		Results     []auxITunesTrack `json:"results"`
	}

	var v iTunesSearchResponse
	if err := json.NewDecoder(rsp.Body).Decode(&v); err != nil {
		return nil, false, errors.Wrapf(err, "failed to json-decode itunes response")
	}

	if v.ResultCount == 0 {
		return nil, false, nil
	}

	ret := make([]ITunesTrack, len(v.Results))
	for i := range v.Results {
		ret[i].ArtistName = v.Results[i].ArtistName
		ret[i].TrackName = v.Results[i].TrackName
		ret[i].CollectionName = v.Results[i].CollectionName
		ret[i].TrackViewURL = v.Results[i].TrackViewURL
		ret[i].PreviewURL = v.Results[i].PreviewURL
		// Parse and set the ReleaseDate field.
		if release := v.Results[i].ReleaseDate; release != "" {
			t, err := time.Parse(time.RFC3339, release)
			if err == nil {
				ret[i].ReleaseDate = t
			}
		}
	}
	return ret, false, nil
}
