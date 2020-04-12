package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/nishanths/scrobble/appengine/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
)

// Datastore limit per operation, i.e. the number of entities that can be
// put/get/deleted in a single call. App Engine documentation?
const datastoreLimitPerOp = 500

type deleteEntitiesTask struct {
	Namespace string
	Kind      string
}

func (s *server) deleteEntitiesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var t deleteEntitiesTask
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.deleteEntities(ctx, t.Namespace, t.Kind); err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Deletes entities of the given kind in the namespace.
func (svr *server) deleteEntities(ctx context.Context, namespace string, kind string) error {
	log.Infof("deleting namespace=%s, kind=%s", namespace, kind)

	allKeys, err := svr.ds.GetAll(ctx, datastore.NewQuery(kind).Namespace(namespace).KeysOnly(), nil)
	if err != nil {
		return errors.Wrapf(err, "failed to get keys")
	}

	s := 0
	e := min(s+datastoreLimitPerOp, len(allKeys))
	chunk := allKeys[s:e]

	for len(chunk) > 0 {
		if err := deleteKeysChunk(ctx, chunk, svr.ds); err != nil {
			return errors.Wrapf(err, "failed to delete chunk")
		}

		s = e
		e = min(s+datastoreLimitPerOp, len(allKeys))
		chunk = allKeys[s:e]
	}

	return nil
}

func deleteKeysChunk(c context.Context, keys []*datastore.Key, ds *datastore.Client) error {
	if len(keys) > datastoreLimitPerOp {
		panic(fmt.Sprintf("length must be <= %d, got %d", datastoreLimitPerOp, len(keys)))
	}
	if len(keys) == 0 {
		return nil
	}
	return ds.DeleteMulti(c, keys)
}

func trimSongParents(ctx context.Context, namespace string, createdBefore int64, ds *datastore.Client) error {
	f := func() error {
		log.Infof("about to delete SongParents created before %d", createdBefore)

		q := datastore.NewQuery(KindSongParent).
			Namespace(namespace).
			Filter("Created <", createdBefore).
			KeysOnly()

		toDeleteSpKeys, err := ds.GetAll(ctx, q, nil)
		if err != nil {
			return errors.Wrapf(err, "failed to get SongParents")
		}

		g, gctx := errgroup.WithContext(ctx)

		for _, spKey := range toDeleteSpKeys {
			log.Infof("deleting %s as part of trimming", spKey)
		}

		// Gather and delete the Songs under each SongParent.
		for _, spKey := range toDeleteSpKeys {
			spKey := spKey // for closure
			g.Go(func() error {
				var del []*datastore.Key // accumulate keys to delete
				q := datastore.NewQuery(KindSong).Namespace(namespace).Ancestor(spKey).KeysOnly()

				for t := ds.Run(gctx, q); ; {
					songKey, err := t.Next(nil)
					if err == iterator.Done {
						break
					}
					if err != nil {
						return err
					}

					del = append(del, songKey)
					if len(del) == datastoreLimitPerOp {
						if err := deleteKeysChunk(gctx, del, ds); err != nil {
							return err
						}
						del = del[:0]
					}
				}

				// delete last remaining (if any)
				if err := deleteKeysChunk(gctx, del, ds); err != nil {
					return err
				}
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		// Now, delete the SongParents.
		s := 0
		e := min(s+datastoreLimitPerOp, len(toDeleteSpKeys))
		chunk := toDeleteSpKeys[s:e]

		for len(chunk) > 0 {
			if err := deleteKeysChunk(ctx, chunk, ds); err != nil {
				log.Errorf("failed to delete chunk: %v", err.Error())
				return errors.Wrapf(err, "failed to delete chunk")
			}

			s = e
			e = min(s+datastoreLimitPerOp, len(toDeleteSpKeys))
			chunk = toDeleteSpKeys[s:e]
		}

		return nil
	}

	return f()
}
