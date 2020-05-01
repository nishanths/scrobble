package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
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
	fParallel = flag.Int("parallel", 64, "number of worker goroutines")
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
	names := make(chan string)

	g.Go(func() error {
		defer close(names)

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
			names <- attrs.Name
		}
		return nil
	})

	for i := 0; i < *fParallel; i++ {
		g.Go(func() error {
			return handle(ctx, names, bkt, ds)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func handle(ctx context.Context, names chan string, bkt *storage.BucketHandle, ds *datastore.Client) error {
	for name := range names {
		if err := handleOne(ctx, name, bkt, ds); err != nil {
			return fmt.Errorf("handle name %s: %s", name, err)
		}
	}
	return nil
}

func handleOne(ctx context.Context, name string, bkt *storage.BucketHandle, ds *datastore.Client) error {
	obj := bkt.Object(name)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("read storage object: %s", err)
	}
	defer r.Close()

	img, _, err := image.Decode(r)
	if err != nil {
		return fmt.Errorf("decode image: %s", err)
	}

	palette := vibrant.NewPaletteBuilder(img).Generate()
	swatches := artwork.Swatches(palette.Swatches())

	artworkScore := artwork.ArtworkScoreFromSwatches(swatches)
	hash := strings.TrimPrefix(name, "aw/")

	log.Printf("%+v", artworkScore)

	// Put global ArtworkScore entity.
	if _, err := ds.Put(ctx, artwork.ArtworkScoreKey(hash), &artworkScore); err != nil {
		return fmt.Errorf("datastore put artwork score: %s", err)
	}

	log.Printf("put %s", hash)
	return nil
}
