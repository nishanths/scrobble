package artwork

import (
	"fmt"
	"image/color"
	"sort"

	"cloud.google.com/go/datastore"
	"github.com/RobCherry/vibrant"
	"github.com/nishanths/scrobble/appengine/basiccolor"
)

const (
	KindArtworkScore  = "ArtworkScore"  // namespace: [default]
	KindArtworkRecord = "ArtworkRecord" // namespace: Account
)

func ArtworkScoreKey(hash string) *datastore.Key {
	return &datastore.Key{Kind: KindArtworkScore, Name: hash}
}

func ArtworkRecordKey(namespace string, hash string) *datastore.Key {
	return &datastore.Key{
		Kind:      KindArtworkRecord,
		Name:      hash,
		Namespace: namespace,
	}
}

// Namespace: [default]
// Key: artwork hash
type ArtworkScore struct {
	Red,
	Orange,
	Brown,
	Yellow,
	Green,
	Blue,
	Violet,
	Pink,
	Black,
	Gray,
	White int
}

// Namespace: Account
// Key: artwork hash
type ArtworkRecord struct {
	Score     ArtworkScore
	SongIdent string
}

type HSL struct {
	H, S, L float64
	A       int
}

type Swatch struct {
	Color      HSL
	Population int
}

func Swatches(swatches []*vibrant.Swatch) []Swatch {
	var total float64
	for _, s := range swatches {
		total += float64(s.Population())
	}

	sort.Slice(swatches, func(i, j int) bool {
		return swatches[i].Population() > swatches[j].Population()
	})

	out := make([]Swatch, len(swatches))
	for i, s := range swatches {
		perMile := int(float64(s.Population()) / total * 1000) // percentage of total, but using 1000 instead of 100
		hsl := vibrant.HSLModel.Convert(s.Color()).(vibrant.HSL)
		out[i] = Swatch{
			Color:      HSL{hsl.H, hsl.S, hsl.L, int(hsl.A)},
			Population: perMile,
		}
	}
	return out
}

func ArtworkScoreFromSwatches(swatches []Swatch) ArtworkScore {
	scores := make(map[basiccolor.Color]int)

	for _, swatch := range swatches {
		col := basiccolor.Closest(toColor(swatch.Color))
		scores[col] += swatch.Population
	}

	return toArtworkScore(scores)
}

const maxBasicColors = 5

func toArtworkScore(scores map[basiccolor.Color]int) ArtworkScore {
	var result ArtworkScore

	type score struct {
		color basiccolor.Color
		score int
	}

	// Get the top `maxBasicColors` {color, score} pairs
	// from the map.
	scoresSlice := make([]score, 0, len(scores))
	for c, s := range scores {
		scoresSlice = append(scoresSlice, score{c, s})
	}
	sort.Slice(scoresSlice, func(i, j int) bool {
		return scoresSlice[i].score > scoresSlice[j].score
	})
	if len(scoresSlice) > maxBasicColors {
		scoresSlice = scoresSlice[:maxBasicColors]
	}

	for _, cs := range scoresSlice {
		switch cs.color {
		case basiccolor.Red:
			result.Red = cs.score
		case basiccolor.Orange:
			result.Orange = cs.score
		case basiccolor.Brown:
			result.Brown = cs.score
		case basiccolor.Yellow:
			result.Yellow = cs.score
		case basiccolor.Green:
			result.Green = cs.score
		case basiccolor.Blue:
			result.Blue = cs.score
		case basiccolor.Violet:
			result.Violet = cs.score
		case basiccolor.Pink:
			result.Pink = cs.score
		case basiccolor.Black:
			result.Black = cs.score
		case basiccolor.Gray:
			result.Gray = cs.score
		case basiccolor.White:
			result.White = cs.score
		default:
			panic(fmt.Sprintf("unhandled basic color: %v", cs.color))
		}
	}

	return result
}

func toColor(c HSL) color.Color {
	return basiccolor.HSL{
		H: c.H,
		S: c.S,
		L: c.L,
		A: uint8(c.A),
	}
}
