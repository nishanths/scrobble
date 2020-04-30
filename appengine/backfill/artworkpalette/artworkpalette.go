package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"sort"
	"strings"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/RobCherry/vibrant"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const credFile = "selective-scrobble-354ed4c58385.json"
const DefaultBucketName = "selective-scrobble.appspot.com"

const (
	KindArtwork = "Artwork"
)

type Palette []Swatch

type HSL struct {
	H, S, L float64
}

type Swatch struct {
	Color      HSL
	Population int
}

// Namespace: [default]
// Key: artwork hash
type Artwork struct {
	Palette Palette
}

func artworkKey(hash string) *datastore.Key {
	return &datastore.Key{Kind: KindArtwork, Name: hash}
}

var (
	fParallel = flag.Int("parallel", 16, "number of worker goroutines")
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	ds, err := datastore.NewClient(ctx, "selective-scrobble", option.WithCredentialsFile(credFile))
	if err != nil {
		return fmt.Errorf("datastore client: %s", err)
	}
	defer ds.Close()

	cloudStorage, err := storage.NewClient(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		return fmt.Errorf("storage client: %s", err)
	}
	defer cloudStorage.Close()

	bkt := cloudStorage.Bucket(DefaultBucketName)

	g, ctx := errgroup.WithContext(ctx)
	workCh := make(chan string)

	for i := 0; i < *fParallel; i++ {
		g.Go(func() error {
			return handle(ctx, ds, bkt, workCh)
		})
	}

	g.Go(func() error {
		defer close(workCh)

		q := &storage.Query{Prefix: "aw/"}
		it := bkt.Objects(ctx, q)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("iterator: %s", err)
			}
			if attrs.Name == "aw/" {
				continue // skip directory-only entry
			}
			workCh <- attrs.Name
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func handle(ctx context.Context, ds *datastore.Client, bkt *storage.BucketHandle, ch chan string) error {
	for name := range ch {
		if err := handleOne(ctx, ds, bkt, name); err != nil {
			return err
		}
	}
	return nil
}

func handleOne(ctx context.Context, ds *datastore.Client, bkt *storage.BucketHandle, name string) error {
	obj := bkt.Object(name)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("read storage object: %s", err)
	}
	defer r.Close()

	top, err := topSwatchesFromFile(r)
	if err != nil {
		return fmt.Errorf("top swatches: %s", err)
	}

	hash := strings.TrimPrefix(name, "aw/")
	entity := &Artwork{Palette: top}
	if _, err := ds.Put(ctx, artworkKey(hash), entity); err != nil {
		return fmt.Errorf("datastore put: %s", err)
	}

	log.Printf("put %s", hash)
	return nil
}

func topSwatchesFromFile(r io.Reader) ([]Swatch, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode image: %s", err)
	}

	palette := vibrant.NewPaletteBuilder(img).Generate()
	return topSwatches(palette.Swatches()), nil

}

func topSwatches(swatches []*vibrant.Swatch) []Swatch {
	var total float64
	for _, s := range swatches {
		total += float64(s.Population())
	}

	sort.Slice(swatches, func(i, j int) bool {
		return swatches[i].Population() > swatches[j].Population()
	})

	if len(swatches) > 5 {
		swatches = swatches[:5]
	}

	out := make([]Swatch, len(swatches))
	for i, s := range swatches {
		perMile := int(float64(s.Population()) / total * 1000) // percentage of total, but using 1000 instead of 100
		hsl := vibrant.HSLModel.Convert(s.Color()).(vibrant.HSL)
		out[i] = Swatch{HSL{hsl.H, hsl.S, hsl.L}, perMile}
	}
	return out
}
