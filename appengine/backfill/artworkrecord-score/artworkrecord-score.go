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

const credFile = "../selective-scrobble-8033d812577b.json"

var (
	namespaces = []string{
		"6e697368616e74682e6765727261726440676d61696c2e636f6d",
	}
)

var (
	fParallel = flag.Int("parallel", 64, "number of worker goroutines")
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

type work struct {
	namespace string
	hash      string
}

func run(ctx context.Context) error {
	ds, err := datastore.NewClient(ctx, "selective-scrobble", option.WithCredentialsFile(credFile))
	if err != nil {
		return fmt.Errorf("datastore client: %s", err)
	}
	defer ds.Close()

	g, ctx := errgroup.WithContext(ctx)
	workCh := make(chan work)

	g.Go(func() error {
		defer close(workCh)

		for _, namespace := range namespaces {
			q := datastore.NewQuery(artwork.KindArtworkRecord).Namespace(namespace).KeysOnly()
			it := ds.Run(ctx, q)
			for {
				key, err := it.Next(nil)
				if err == iterator.Done {
					break
				}
				if err != nil {
					return fmt.Errorf("iterator: %s", err)
				}
				workCh <- work{key.Namespace, key.Name}
			}
		}
		return nil
	})

	for i := 0; i < *fParallel; i++ {
		g.Go(func() error {
			return handle(ctx, ds, workCh)
		})
	}

	return g.Wait()
}

func handle(ctx context.Context, ds *datastore.Client, workCh chan work) error {
	for k := range workCh {
		if err := handleOne(ctx, ds, k); err != nil {
			return fmt.Errorf("handling %s", err)
		}
	}
	return nil
}

func handleOne(ctx context.Context, ds *datastore.Client, w work) error {
	namespace := w.namespace
	hash := w.hash

	var as artwork.ArtworkScore
	err := ds.Get(ctx, artwork.ArtworkScoreKey(hash), &as)
	if err == datastore.ErrNoSuchEntity {
		log.Printf("no global ArtworkScore for %s", hash)
		return nil
	}
	if err != nil {
		return fmt.Errorf("datastore get artwork score: %s", err)
	}

	_, err = ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var ar artwork.ArtworkRecord
		err := tx.Get(artwork.ArtworkRecordKey(namespace, hash), &ar)
		if err != nil {
			return fmt.Errorf("datastore get artwork record: %s", err)
		}
		ar.Score = as
		_, err = tx.Put(artwork.ArtworkRecordKey(namespace, hash), &ar)
		if err != nil {
			return fmt.Errorf("datastore put artwork record: %s", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	log.Printf("put %v", w)
	return nil
}
