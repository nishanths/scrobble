package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/nishanths/scrobble/appengine/artwork"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const credFile = "../selective-scrobble-354ed4c58385.json"

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

	g, ctx := errgroup.WithContext(ctx)
	hashes := make(chan string)

	for i := 0; i < *fParallel; i++ {
		g.Go(func() error {
			return handle(ctx, hashes, ds)
		})
	}

	g.Go(func() error {
		defer close(hashes)

		q := datastore.NewQuery(artwork.KindArtwork).KeysOnly()
		it := ds.Run(ctx, q)
		for {
			key, err := it.Next(nil)
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("iterator: %s", err)
			}
			hashes <- key.Name
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func handle(ctx context.Context, hashes chan string, ds *datastore.Client) error {
	for hash := range hashes {
		if err := handleOne(ctx, hash, ds); err != nil {
			return fmt.Errorf("handle hash %s: %s", hash, err)
		}
	}
	return nil
}

func handleOne(ctx context.Context, hash string, ds *datastore.Client) error {
	var a artwork.Artwork
	if err := ds.Get(ctx, artwork.ArtworkKey(hash), &a); err != nil {
		return fmt.Errorf("datastore get: %s", err)
	}

	as := artwork.ArtworkScoreFromSwatches(a.Palette)

	if _, err := ds.Put(ctx, artwork.ArtworkScoreKey(hash), &as); err != nil {
		return fmt.Errorf("datastore put: %s", err)
	}

	log.Printf("put %s", hash)
	return nil
}
