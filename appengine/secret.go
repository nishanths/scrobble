package main

import (
	"context"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/nishanths/scrobble/appengine/log"
	"github.com/pkg/errors"
)

type Secret struct {
	CookieHashKey      string `datastore:",noindex"`
	CookieBlockKey     string `datastore:",noindex"`
	GoogleClientID     string `datastore:",noindex"`
	GoogleClientSecret string `datastore:",noindex"`
	TasksSecret        string `datastore:",noindex"`
}

func fetchSecret(ctx context.Context, ds *datastore.Client) (Secret, error) {
	var secret Secret
	if err := ds.Get(ctx, datastore.NameKey(KindSecret, "singleton", nil), &secret); err != nil {
		return Secret{}, errors.Wrapf(err, "failed to get from datastore")
	}
	return secret, nil
}

const headerTasksSecret = "X-Scrobble-Tasks-Secret"

func (s *server) requireTasksSecretHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get(headerTasksSecret)
		if got == "" {
			log.Errorf("missing tasks secret header in request")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		want := s.secret.TasksSecret
		if want == "" {
			panic("empty tasks secret")
		}

		if want != got {
			log.Errorf("bad tasks secret header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
