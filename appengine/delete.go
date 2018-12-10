package main

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
)

// Datastore limit per operation, i.e. the number of entities that can be
// put/get/deleted in a single call. App Engine documentation?
const datastoreLimitPerOp = 500

// Deletes entities of the given kind in the namespace.
var deleteFunc = delay.Func("delete", func(ctx context.Context, namespace string, kind string) error {
	log.Infof(ctx, "deleting namespace=%s, kind=%s", namespace, kind)

	ns, err := appengine.Namespace(ctx, namespace)
	if err != nil {
		log.Errorf(ctx, "failed to make namespace: %v", err.Error())
		return errors.Wrapf(err, "failed to make namespace")
	}

	allKeys, err := datastore.NewQuery(kind).KeysOnly().GetAll(ns, nil)
	if err != nil {
		log.Errorf(ctx, "failed to get keys: %v", err.Error())
		return errors.Wrapf(err, "failed to get keys")
	}

	s := 0
	e := min(s+datastoreLimitPerOp, len(allKeys))
	chunk := allKeys[s:e]

	for len(chunk) > 0 {
		if err := deleteKeysChunk(ns, chunk); err != nil {
			log.Errorf(ns, "failed to delete chunk: %v", err.Error())
			return errors.Wrapf(err, "failed to delete chunk")
		}

		s = e
		e = min(s+datastoreLimitPerOp, len(allKeys))
		chunk = allKeys[s:e]
	}

	return nil
})

func deleteKeysChunk(c context.Context, keys []*datastore.Key) error {
	if len(keys) > datastoreLimitPerOp {
		panic(fmt.Sprintf("length must be <= %d, got %d", datastoreLimitPerOp, len(keys)))
	}
	if len(keys) == 0 {
		return nil
	}
	return datastore.DeleteMulti(c, keys)
}

func trimSongParents(ns context.Context) error {
	f := func(complete bool, trimAt, nDelete int) error {
		// Get the oldest up to the limit.
		q := datastore.NewQuery(KindSongParent).
			Order("Created").Filter("Complete=", complete).
			Limit(trimAt).KeysOnly()

		spKeys, err := q.GetAll(ns, nil)
		if err != nil {
			return errors.Wrapf(err, "failed to get SongParents")
		}

		if len(spKeys) < trimAt {
			return nil // don't trim
		}

		toDeleteSpKeys := spKeys[:nDelete]
		g, gns := errgroup.WithContext(ns)

		for _, spKey := range toDeleteSpKeys {
			log.Infof(ns, "deleting %s as part of trimming", spKey)
		}

		// Gather and delete the Songs under each SongParent.
		for _, spKey := range toDeleteSpKeys {
			spKey := spKey // for closure
			g.Go(func() error {
				var del []*datastore.Key // accumulate keys to delete
				q := datastore.NewQuery(KindSong).Ancestor(spKey).KeysOnly()

				for t := q.Run(gns); ; {
					songKey, err := t.Next(nil)
					if err == datastore.Done {
						break
					}
					if err != nil {
						return err
					}

					del = append(del, songKey)
					if len(del) == datastoreLimitPerOp {
						if err := deleteKeysChunk(gns, del); err != nil {
							return err
						}
						del = del[:0]
					}
				}

				// delete last remaining (if any)
				if err := deleteKeysChunk(gns, del); err != nil {
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
			if err := deleteKeysChunk(ns, chunk); err != nil {
				log.Errorf(ns, "failed to delete chunk: %v", err.Error())
				return errors.Wrapf(err, "failed to delete chunk")
			}

			s = e
			e = min(s+datastoreLimitPerOp, len(toDeleteSpKeys))
			chunk = toDeleteSpKeys[s:e]
		}

		return nil
	}

	if err := f(true, 10, 5); err != nil {
		return err
	}
	if err := f(false, 5, 3); err != nil {
		return err
	}
	return nil
}
