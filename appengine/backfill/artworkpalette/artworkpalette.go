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
	"strings"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/RobCherry/vibrant"
	"github.com/nishanths/scrobble/appengine/artwork"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const credFile = "../selective-scrobble-354ed4c58385.json"
const DefaultBucketName = "selective-scrobble.appspot.com"

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
	entity := &artwork.Artwork{Palette: top}
	if _, err := ds.Put(ctx, artwork.ArtworkKey(hash), entity); err != nil {
		return fmt.Errorf("datastore put: %s", err)
	}

	log.Printf("put %s", hash)
	return nil
}

func topSwatchesFromFile(r io.Reader) ([]artwork.Swatch, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode image: %s", err)
	}

	palette := vibrant.NewPaletteBuilder(img).Generate()
	return artwork.TopSwatches(palette.Swatches()), nil
}
