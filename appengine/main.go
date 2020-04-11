package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta2"
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/nishanths/scrobble/appengine/log"
	"github.com/pkg/errors"
)

const (
	DefaultBucketName = "selective-scrobble.appspot.com"
	DefaultQueueName  = "projects/selective-scrobble/locations/us-east1/queues/default"
)

type server struct {
	ds         *datastore.Client
	storage    *storage.Client
	tasks      *cloudtasks.Client
	httpClient *http.Client
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Errorf(ctx, err.Error())
	}
}

func run(ctx context.Context) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var ds *datastore.Client         // nil in dev
	var cloudStorage *storage.Client // ...
	var tasks *cloudtasks.Client

	if !isDev() {
		var err error

		ds, err = datastore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
		if err != nil {
			return errors.Wrapf(err, "datastore client")
		}
		defer ds.Close()

		cloudStorage, err = storage.NewClient(ctx)
		if err != nil {
			return errors.Wrapf(err, "storage client")
		}
		defer cloudStorage.Close()

		tasks, err = cloudtasks.NewClient(ctx)
		if err != nil {
			return errors.Wrapf(err, "tasks client")
		}
		defer tasks.Close()
	}

	s := &server{
		ds:         ds,
		storage:    cloudStorage,
		tasks:      tasks,
		httpClient: &http.Client{},
	}

	// Register handlers.
	if isDev() {
		http.HandleFunc("/", devRootHandler)
		http.HandleFunc("/u/", devUHandler)
	} else {
		http.Handle("/", withHTTPS(http.HandlerFunc(s.rootHandler)))
		http.Handle("/u/", withHTTPS(http.HandlerFunc(s.uHandler)))
	}
	http.Handle("/initializeAccount", withHTTPS(http.HandlerFunc(s.initializeAccountHandler)))
	http.Handle("/newAPIKey", withHTTPS(http.HandlerFunc(s.newAPIKeyHandler)))
	http.Handle("/setPrivacy", withHTTPS(http.HandlerFunc(s.setPrivacyHandler)))

	if isDev() {
		http.HandleFunc("/api/v1/scrobbled", devScrobbledHandler)
	} else {
		http.HandleFunc("/api/v1/scrobbled", s.scrobbledHandler)
	}
	http.HandleFunc("/api/v1/scrobble", s.scrobbleHandler)
	http.HandleFunc("/api/v1/account", s.accountHandler)
	http.HandleFunc("/api/v1/account/delete", s.deleteAccountHandler)
	http.HandleFunc("/api/v1/artwork", s.artworkHandler)
	http.HandleFunc("/api/v1/artwork/missing", s.artworkMissingHandler)

	http.Handle("/internal/fillITunesFields", s.requireTasksSecretHeader(http.HandlerFunc(s.fillITunesFieldsHandler)))
	http.Handle("/internal/markParentComplete", s.requireTasksSecretHeader(http.HandlerFunc(s.markParentCompleteHandler)))
	http.Handle("/internal/deleteEntities", s.requireTasksSecretHeader(http.HandlerFunc(s.deleteEntitiesHandler)))

	if isDev() {
		// in production these are handled by app.yaml
		http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir(filepath.Join("web", "dist")))))
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join("web", "static")))))
	}

	// Serve.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return errors.Wrapf(err, "ListenAndServe")
	}

	panic("should not be reachable")
}

func isDev() bool {
	return os.Getenv("GAE_DEPLOYMENT_ID") == ""
}

func withHTTPS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeRedirectHTTPS(w, r) {
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Redirect requests with a "http" scheme to "https".
func maybeRedirectHTTPS(w http.ResponseWriter, r *http.Request) bool {
	u := *r.URL
	if u.Scheme != "http" {
		return false
	}
	u.Scheme = "https"
	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
	return true
}
