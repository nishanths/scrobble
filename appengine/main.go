package main

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
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
	secret     *Secret
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf(err.Error())
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
	var secret *Secret

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

		s, err := fetchSecret(ctx, ds)
		if err != nil {
			return errors.Wrapf(err, "fetching secret")
		}
		secret = &s
	}

	s := &server{
		ds:         ds,
		storage:    cloudStorage,
		tasks:      tasks,
		httpClient: &http.Client{},
		secret:     secret,
	}

	webMiddleware := func(h http.Handler) http.Handler {
		return withHTTPS(withAlleleRedirect(h))
	}

	// Register handlers.

	// web handlers
	if isDev() {
		http.HandleFunc("/", devRootHandler)
		http.HandleFunc("/u/", devUHandler)
	} else {
		http.Handle("/", webMiddleware(http.HandlerFunc(s.rootHandler)))
		http.Handle("/u/", webMiddleware(http.HandlerFunc(s.uHandler)))
	}
	http.Handle("/initializeAccount", webMiddleware(http.HandlerFunc(s.initializeAccountHandler)))
	http.Handle("/newAPIKey", webMiddleware(http.HandlerFunc(s.newAPIKeyHandler)))
	http.Handle("/setPrivacy", webMiddleware(http.HandlerFunc(s.setPrivacyHandler)))
	http.Handle("/login", webMiddleware(http.HandlerFunc(s.loginHandler)))
	http.Handle("/googleAuth", webMiddleware(http.HandlerFunc(s.googleAuthHandler)))
	http.Handle("/logout", webMiddleware(http.HandlerFunc(s.logoutHandler)))
	http.Handle("/privacy-policy", webMiddleware(http.HandlerFunc(s.privacyPolicyHandler)))

	// doc handlers
	http.Handle("/doc/api/v1/", http.StripPrefix("/doc/api/v1/", http.FileServer(http.Dir(filepath.Join("doccontent", "api", "dst")))))
	http.Handle("/doc/guide/", http.StripPrefix("/doc/guide/", http.FileServer(http.Dir(filepath.Join("doccontent", "guide", "dst")))))
	if isDev() {
		// in production this is handled by app.yaml
		http.HandleFunc("/doc/style.css", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join("doccontent", "style.css"))
		})
	}

	// API handlers
	if isDev() {
		http.HandleFunc("/api/v1/scrobbled", devScrobbledHandler)
		http.HandleFunc("/api/v1/scrobbled/color", devScrobbledColorHandler)
	} else {
		http.HandleFunc("/api/v1/scrobbled", s.scrobbledHandler)
		http.HandleFunc("/api/v1/scrobbled/color", s.scrobbledColorHandler)
	}
	http.HandleFunc("/api/v1/scrobble", s.scrobbleHandler)
	http.HandleFunc("/api/v1/account", s.accountHandler)
	http.HandleFunc("/api/v1/account/delete", s.deleteAccountHandler)
	http.HandleFunc("/api/v1/artwork", s.artworkHandler)
	http.HandleFunc("/api/v1/artwork/missing", s.artworkMissingHandler)

	// data API handlers
	http.HandleFunc("/api/v1/songs/play-count", s.songPlayCountHandler)
	http.HandleFunc("/api/v1/songs/length", s.songLengthHandler)
	http.HandleFunc("/api/v1/artists/play-count", s.artistPlayCountHandler)
	http.HandleFunc("/api/v1/artists/added", s.artistAddedHandler)

	// internal handlers
	http.Handle("/internal/fillITunesFields", s.requireTasksSecretHeader(http.HandlerFunc(s.fillITunesFieldsHandler)))
	http.Handle("/internal/markParentComplete", s.requireTasksSecretHeader(http.HandlerFunc(s.markParentCompleteHandler)))
	http.Handle("/internal/deleteEntities", s.requireTasksSecretHeader(http.HandlerFunc(s.deleteEntitiesHandler)))
	http.Handle("/internal/fillArtworkScore", s.requireTasksSecretHeader(http.HandlerFunc(s.fillArtworkScoreHandler)))
	http.Handle("/internal/computeArtistStats", s.requireTasksSecretHeader(http.HandlerFunc(s.computeArtistStatsHandler)))

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

func withAlleleRedirect(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == "scrobble.allele.cc" {
			u := *r.URL
			u.Host = "scrobble.littleroot.org"
			http.Redirect(w, r, u.String(), http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func drainAndClose(r io.ReadCloser) {
	io.Copy(ioutil.Discard, r)
	r.Close()
}
