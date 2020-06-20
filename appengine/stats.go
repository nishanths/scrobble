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

func statsArtistPlayCountKey(namespace string) *datastore.Key {
	return &datastore.Key{
		Kind:      KindStats,
		Name:      "artist-playCount",
		Namespace: namespace,
	}
}

func statsArtistAddedKey(namespace string) *datastore.Key {
	return &datastore.Key{
		Kind:      KindStats,
		Name:      "artist-added",
		Namespace: namespace,
	}
}

// namespace: account
type ArtistPlayCountStats struct {
	Data []ArtistPlayCountDatum `datastore:",noindex"`
}

type ArtistPlayCountDatum struct {
	ArtistName    string `datastore:",noindex" json:"artistName"`
	PlayCount     int    `datastore:",noindex" json:"playCount"`
	TotalPlayTime int    `datastore:",noindex" json:"totalPlayTime"` // in seconds
}

type ArtistAddedStats struct {
	Data []ArtistAddedDatum `datastore:",noindex"`
}

type ArtistAddedDatum struct {
	ArtistName string `datastore:",noindex" json:"artistName"`
	Added      int64  `datastore:",noindex" json:"added"`
	Song       Song   `datastore:",noindex" json:"song"` // song that was added
}

func computeArtistPlayCount(songs []Song) ArtistPlayCountStats {
	m := make(map[string]ArtistPlayCountDatum)
	for _, s := range songs {
		v := m[s.ArtistName]
		v.ArtistName = s.ArtistName
		v.PlayCount += s.PlayCount
		v.TotalPlayTime += int(s.TotalTime/time.Second) * s.PlayCount
		m[s.ArtistName] = v
	}

	slice := make([]ArtistPlayCountDatum, 0, len(m))
	for _, v := range m {
		slice = append(slice, v)
	}

	// sort desc.
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].PlayCount > slice[j].PlayCount
	})

	return ArtistPlayCountStats{
		Data: slice,
	}
}

func computeArtistAdded(songs []Song) ArtistAddedStats {
	m := make(map[string]ArtistAddedDatum)
	for _, s := range songs {
		if v, ok := m[s.ArtistName]; ok {
			if s.Added < v.Added {
				// earliest added date
				datum := m[s.ArtistName]
				datum.Added = s.Added
				m[s.ArtistName] = datum
			} else {
				// latest added song
				datum := m[s.ArtistName]
				datum.Song = s
				m[s.ArtistName] = datum
			}
		} else {
			m[s.ArtistName] = ArtistAddedDatum{s.ArtistName, s.Added, s}
		}
	}

	slice := make([]ArtistAddedDatum, 0, len(m))
	for _, v := range m {
		slice = append(slice, v)
	}

	// sort by added desc. (latest added appears first)
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Added > slice[j].Added
	})

	return ArtistAddedStats{
		Data: slice,
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
		Ancestor(parentKey)

	var songs []Song
	if _, err := s.ds.GetAll(ctx, q, &songs); err != nil {
		log.Errorf("failed to get songs: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pc := computeArtistPlayCount(songs)
	pcKey := statsArtistPlayCountKey(t.Namespace)

	a := computeArtistAdded(songs)
	aKey := statsArtistAddedKey(t.Namespace)

	if _, err := s.ds.PutMulti(ctx, []*datastore.Key{pcKey, aKey}, []interface{}{&pc, &a}); err != nil {
		log.Errorf("failed to put artist stats: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
