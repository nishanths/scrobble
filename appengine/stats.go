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
	statsStorageDirectory = "stats" // Cloud Storage directory name for stats
)

func statsArtistPlayCountStoragePath(namespace string) string {
	return statsStorageDirectory + "/" + namespace + "/" + "artist-playCount"
}

func statsArtistAddedStoragePath(namespace string) string {
	return statsStorageDirectory + "/" + namespace + "/" + "artist-added"
}

// namespace: account
type ArtistPlayCountStats struct {
	Data []ArtistPlayCountDatum `datastore:",noindex" json:"data"`
}

type ArtistPlayCountDatum struct {
	ArtistName    string `datastore:",noindex" json:"artistName"`
	PlayCount     int    `datastore:",noindex" json:"playCount"`
	TotalPlayTime int    `datastore:",noindex" json:"totalPlayTime"` // in seconds
}

type ArtistAddedStats struct {
	Data []ArtistAddedDatum `datastore:",noindex" json:"data"`
}

type ArtistAddedDatum struct {
	ArtistName   string `datastore:",noindex" json:"artistName"`
	Added        int64  `datastore:",noindex" json:"added"`        // earliest song added date for artist library-wide
	LatestSong   Song   `datastore:",noindex" json:"latestSong"`   // latest added song for artist library-wide
	EarliestSong Song   `datastore:",noindex" json:"earliestSong"` // earliest added song for artist library-wide
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
				// record the earliest added date, earliest added song (progressively)
				datum := m[s.ArtistName]
				datum.Added = s.Added
				datum.EarliestSong = s
				m[s.ArtistName] = datum
			} else {
				// record the latest added song (progressively)
				datum := m[s.ArtistName]
				datum.LatestSong = s
				m[s.ArtistName] = datum
			}
		} else {
			m[s.ArtistName] = ArtistAddedDatum{s.ArtistName, s.Added, s, s}
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

	{
		pc := computeArtistPlayCount(songs)
		pcPath := statsArtistPlayCountStoragePath(t.Namespace)

		wr := s.storage.Bucket(DefaultBucketName).Object(pcPath).NewWriter(ctx)
		if err := json.NewEncoder(wr).Encode(pc); err != nil {
			log.Errorf("failed to json-write artist stats: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := wr.Close(); err != nil {
			log.Errorf("failed to close storage writer: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	{
		a := computeArtistAdded(songs)
		aPath := statsArtistAddedStoragePath(t.Namespace)

		wr := s.storage.Bucket(DefaultBucketName).Object(aPath).NewWriter(ctx)
		if err := json.NewEncoder(wr).Encode(a); err != nil {
			log.Errorf("failed to json-encode artist stats: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := wr.Close(); err != nil {
			log.Errorf("failed to close storage writer: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
