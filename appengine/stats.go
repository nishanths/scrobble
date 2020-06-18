package main

import (
	"sort"
	"time"

	"cloud.google.com/go/datastore"
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

type ArtistDatum struct {
	ArtistName     string      `datastore:",noindex"`
	PrimaryValue   interface{} `datastore:",noindex"`
	SecondaryValue interface{} `datastore:",noindex"`
}

// namespace: account
type ArtistStats struct {
	Data         []ArtistDatum `datastore:",noindex"`
	TotalArtists int           `datastore:",noindex"`
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
		Data:         artistData,
		TotalArtists: len(m),
	}
}
