package main

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
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

func deleteKeysChunk(ctx context.Context, keys []*datastore.Key) error {
	if len(keys) > datastoreLimitPerOp {
		panic(fmt.Sprintf("length must be <= %d, got %d", datastoreLimitPerOp, len(keys)))
	}
	return datastore.DeleteMulti(ctx, keys)
}
