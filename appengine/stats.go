package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/nishanths/scrobble/appengine/log"
)

const (
	KindStats = "Stats" // namespace: account
)

func statsPlayCountArtistKey(namespace string) *datastore.Key {
	return &datastore.Key{
		Kind:      KindStats,
		Name:      "playCount-artist",
		Namespace: namespace,
	}
}

func statsLastPlayedArtistKey(namespace string) *datastore.Key {
	return &datastore.Key{
		Kind:      KindStats,
		Name:      "lastPlayed-artist",
		Namespace: namespace,
	}
}

// namespace: account
type ArtistStats struct {
	Data []ArtistDatum `datastore:",noindex"`
}

type ArtistDatum struct {
	ArtistName     string      `datastore:",noindex"`
	PrimaryValue   interface{} `datastore:",noindex"`
	SecondaryValue interface{} `datastore:",noindex"`
}

const maxArtistStatsLen = 20

func computePlayCountArtistStats(songs []Song) ArtistStats {
	type value struct {
		artistName    string
		playCount     int
		totalPlayTime int
	}

	m := make(map[string]value)
	for _, s := range songs {
		v := m[s.ArtistName]
		v.artistName = s.ArtistName
		v.playCount += s.PlayCount
		v.totalPlayTime += int(s.TotalTime/time.Second) * s.PlayCount
		m[s.ArtistName] = v
	}

	slice := make([]value, 0, len(m))
	for _, v := range m {
		slice = append(slice, v)
	}

	// sort desc.
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].playCount > slice[j].playCount
	})

	var artistData []ArtistDatum
	for i := 0; i < len(slice) && len(artistData) < maxArtistStatsLen; i++ {
		v := slice[i]
		artistData = append(artistData, ArtistDatum{
			ArtistName:     v.artistName,
			PrimaryValue:   v.playCount,
			SecondaryValue: v.totalPlayTime,
		})
	}

	return ArtistStats{
		Data: artistData,
	}
}

// songs must be sorted by last played times desc.
func computeLastPlayedArtistsStats(songs []Song) ArtistStats {
	m := make(map[string]struct{})
	var data []ArtistDatum

	for _, s := range songs {
		if _, ok := m[s.ArtistName]; !ok {
			// first time
			data = append(data, ArtistDatum{ArtistName: s.ArtistName})
			m[s.ArtistName] = struct{}{}
		}
		if len(m) == maxArtistStatsLen {
			break
		}
	}

	return ArtistStats{
		Data: data,
	}
}

type computeArtistStatsTask struct {
	Namespace       string
	SongParentIdent string
}

func (s *server) computeArtistStatsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var t computeArtistStatsTask
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Errorf("failed to json-decode task: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parentKey := songParentKey(t.Namespace, t.SongParentIdent)
	q := datastore.NewQuery(KindSong).
		Namespace(t.Namespace).
		Order("-LastPlayed").
		Ancestor(parentKey)

	var songs []Song
	if _, err := s.ds.GetAll(ctx, q, &songs); err != nil {
		log.Errorf("failed to get songs: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pc := computePlayCountArtistStats(songs)
	pcKey := statsPlayCountArtistKey(t.Namespace)

	lp := computeLastPlayedArtistsStats(songs)
	lpKey := statsLastPlayedArtistKey(t.Namespace)

	if _, err := s.ds.PutMulti(ctx, []*datastore.Key{pcKey, lpKey}, []ArtistStats{pc, lp}); err != nil {
		log.Errorf("failed to put artist stats: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
