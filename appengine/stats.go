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

func statsAddedArtistKey(namespace string) *datastore.Key {
	return &datastore.Key{
		Kind:      KindStats,
		Name:      "added-artist",
		Namespace: namespace,
	}
}

// namespace: account
type PlayCountArtistStats struct {
	Data []PlayCountArtistDatum `datastore:",noindex"`
}

type PlayCountArtistDatum struct {
	ArtistName    string `datastore:",noindex" json:"artistName"`
	PlayCount     int    `datastore:",noindex" json:"playCount"`
	TotalPlayTime int    `datastore:",noindex" json:"totalPlayTime"` // in seconds
}

type AddedArtistStats struct {
	Data []AddedArtistDatum `datastore:",noindex"`
}

type AddedArtistDatum struct {
	ArtistName string `datastore:",noindex" json:"artistName"`
	Added      int64  `datastore:",noindex" json:"added"`
}

const maxArtistStatsLen = 20

func computePlayCountArtistStats(songs []Song) PlayCountArtistStats {
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

	var artistData []PlayCountArtistDatum
	for i := 0; i < len(slice) && len(artistData) < maxArtistStatsLen; i++ {
		v := slice[i]
		artistData = append(artistData, PlayCountArtistDatum{
			ArtistName:    v.artistName,
			PlayCount:     v.playCount,
			TotalPlayTime: v.totalPlayTime,
		})
	}

	return PlayCountArtistStats{
		Data: artistData,
	}
}

// songs must be sorted by added times desc.
func computeAddedArtistsStats(songs []Song) AddedArtistStats {
	m := make(map[string]struct{})
	var data []AddedArtistDatum

	for _, s := range songs {
		if _, ok := m[s.ArtistName]; !ok {
			// first time
			data = append(data, AddedArtistDatum{
				ArtistName: s.ArtistName,
				Added:      s.Added,
			})
			m[s.ArtistName] = struct{}{}
		}
		if len(m) == maxArtistStatsLen {
			break
		}
	}

	return AddedArtistStats{
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
		Order("-Added").
		Ancestor(parentKey)

	var songs []Song
	if _, err := s.ds.GetAll(ctx, q, &songs); err != nil {
		log.Errorf("failed to get songs: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pc := computePlayCountArtistStats(songs)
	pcKey := statsPlayCountArtistKey(t.Namespace)

	a := computeAddedArtistsStats(songs)
	aKey := statsAddedArtistKey(t.Namespace)

	if _, err := s.ds.PutMulti(ctx, []*datastore.Key{pcKey, aKey}, []interface{}{&pc, &a}); err != nil {
		log.Errorf("failed to put artist stats: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
