package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/nishanths/scrobble/appengine/log"
	"github.com/pkg/errors"
)

var (
	homeTmpl      = template.Must(template.New("").Parse(MustAssetString("appengine/template/home.html")))
	dashboardTmpl = template.Must(template.New("").Parse(MustAssetString("appengine/template/dashboard.html")))
	uTmpl         = template.Must(template.New("").Parse(MustAssetString("appengine/template/u.html")))
)

type BootstrapArgs struct {
	Host             string  `json:"host"`
	Email            string  `json:"email"`
	LoginURL         string  `json:"loginURL"`
	LogoutURL        string  `json:"logoutURL"`
	Account          Account `json:"account"`
	TotalSongs       int     `json:"totalSongs"`       // -1 if failed to compute
	LastScrobbleTime int64   `json:"lastScrobbleTime"` // unix seconds; -1 if failed to compute
}

type RootArgs struct {
	Title     string
	Bootstrap BootstrapArgs
	AppDomain string
}

func validRootPath(p string) bool {
	if p == "/" {
		return true
	}
	c := pathComponents(p)
	if len(c) > 0 && c[0] == "dashboard" {
		return true
	}
	return false
}

func (s *server) rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if !validRootPath(r.URL.Path) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	host := r.Host
	dest := r.RequestURI

	u, err := s.currentUser(r)

	if err != nil {
		// either generic error or ErrNoUser

		// redirect any "/dashboard/" paths to "/"
		if r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		login := loginURLWithRedirect(dest)
		args := RootArgs{
			Title: "Apple Music scrobbling — Scrobble",
			Bootstrap: BootstrapArgs{
				Host:     host,
				LoginURL: login,
			},
			AppDomain: AppDomain,
		}
		if err := homeTmpl.Execute(w, args); err != nil {
			log.Errorf("failed to execute template: %v", err.Error())
		}
		return
	}

	logout := logoutURLWithRedirect(dest)

	a, err := ensureAccount(ctx, u.Email, s.ds)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	lastScrobbled, nSongs, err := songStats(ctx, s.ds, namespaceID(u.Email))
	if err != nil {
		log.Errorf("failed to count songs (continuing): %v", err.Error())
	}

	args := RootArgs{
		Title: "Dashboard · Scrobble",
		Bootstrap: BootstrapArgs{
			Host:             host,
			Email:            u.Email,
			LogoutURL:        logout,
			Account:          a,
			TotalSongs:       nSongs,
			LastScrobbleTime: lastScrobbled,
		},
		AppDomain: AppDomain,
	}
	if err := dashboardTmpl.Execute(w, args); err != nil {
		log.Errorf("failed to execute template: %v", err.Error())
	}
}

// Returns the last scrobble time and the count of scrobbled songs.
// Partial results may be returned even if the error is non-nil.
// Returns -1 for results that could not be computed.
func songStats(ctx context.Context, ds *datastore.Client, namespace string) (int64, int, error) {
	// Get the latest complete parent.
	q := datastore.NewQuery(KindSongParent).
		Namespace(namespace).
		Order("-Created").Filter("Complete=", true).
		Limit(1)

	var sp []SongParent
	parentKeys, err := ds.GetAll(ctx, q, &sp)
	if err != nil {
		return -1, -1, errors.Wrapf(err, "failed to do SongParent query")
	}
	if len(parentKeys) == 0 {
		// no songs
		return 0, 0, nil
	}

	q = datastore.NewQuery(KindSong).
		Namespace(namespace).
		Ancestor(parentKeys[0])
	count, err := ds.Count(ctx, q)
	if err != nil {
		return sp[0].Created, -1, errors.Wrapf(err, "failed to count songs")
	}
	return sp[0].Created, count, nil
}

func ensureAccount(ctx context.Context, accID string, ds *datastore.Client) (Account, error) {
	var account Account

	_, err := ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		aKey := datastore.NameKey(KindAccount, accID, nil)

		if err := tx.Get(aKey, &account); err != nil {
			if err == datastore.ErrNoSuchEntity {
				// account entity does not exists; create new account entity
				if _, err := tx.Put(aKey, &account); err != nil {
					return errors.Wrapf(err, "failed to put account for %s", accID)
				}
				return nil
			}

			// generic error
			return errors.Wrapf(err, "failed to get account for %s", accID)
		}

		// account entity exists
		return nil
	})

	return account, err
}

type UArgs struct {
	Title           string  `json:"title"`
	Host            string  `json:"host"`
	ArtworkBaseURL  string  `json:"artworkBaseURL"`
	ProfileUsername string  `json:"profileUsername"`
	LogoutURL       string  `json:"logoutURL"`
	Account         Account `json:"account"`
	Self            bool    `json:"self"`
	Private         bool    `json:"private"`
}

func (s *server) uHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	c := pathComponents(r.URL.Path)
	if len(c) < 2 || len(c) > 6 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	profileUsername := c[1]
	acc, _, ok := fetchAccountForUsername(ctx, profileUsername, s.ds, w)
	if !ok {
		return
	}

	// If the user is logged in, gather a logout URL and the account info.
	var logoutURL string
	var account Account
	self := false

	u, err := s.currentUser(r)
	if err == nil {
		logoutURL = logoutURLWithRedirect(r.RequestURI)
		if err := s.ds.Get(ctx, datastore.NameKey(KindAccount, u.Email, nil), &account); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		self = account.Username != "" && account.Username == profileUsername
	}

	if err := uTmpl.Execute(w, UArgs{
		Title:           profileUsername + "'s scrobbles",
		Host:            r.Host,
		ArtworkBaseURL:  "https://storage.googleapis.com/" + DefaultBucketName + "/" + artworkStorageDirectory,
		ProfileUsername: profileUsername,
		LogoutURL:       logoutURL,
		Account:         account,
		Self:            self,
		Private:         acc.Private,
	}); err != nil {
		log.Errorf("failed to execute template: %v", err.Error())
	}
}

func pathComponents(path string) []string {
	var c []string
	parts := strings.Split(path, "/")
	for _, p := range parts {
		if p != "" {
			c = append(c, p)
		}
	}
	return c
}

// Sets the username for the account and initializes the account.
func (s *server) initializeAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	u, err := s.currentUser(r)
	if err == ErrNoUser {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	if ok := isAllowedUsername(username); !ok {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable) // need a unique code. this is gross, but whatever.
		return
	}

	inUse := false
	var account Account
	_, err = s.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		// Ensure username uniqueness.
		uKey := datastore.NameKey(KindUsername, username, nil)
		if err := tx.Get(uKey, ptrStruct()); err != datastore.ErrNoSuchEntity {
			if err == nil {
				inUse = true
				return errors.New("username already in use")
			}
			return errors.Wrapf(err, "failed to get username")
		}
		if _, err := tx.Put(uKey, ptrStruct()); err != nil {
			return errors.Wrapf(err, "failed to put username")
		}

		// Ensure API key uniqueness.
		gen, err := setAPIKey(tx, generateAPIKey)
		if err != nil {
			return errors.Wrapf(err, "failed to set API key")
		}

		// Initialize the account.
		aKey := datastore.NameKey(KindAccount, u.Email, nil)
		if err := tx.Get(aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to get account for %s", u.Email)
		}
		if account.Username != "" {
			return fmt.Errorf("username already set for %s", u.Email)
		}
		account.Username = username
		account.APIKey = gen
		account.Private = false
		if _, err := tx.Put(aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to put account for %s", u.Email)
		}

		return nil
	})

	if err != nil {
		log.Errorf("%v", err.Error())
		if inUse {
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable) // need a unique code. this is gross, but whatever.
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.Marshal(account)
	if err != nil {
		log.Errorf("failed to json-marshal account: %v", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (s *server) newAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	u, err := s.currentUser(r)
	if err == ErrNoUser {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var apiKey string
	_, err = s.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		k, err := setAPIKey(tx, generateAPIKey)
		if err != nil {
			return err
		}

		var account Account
		aKey := datastore.NameKey(KindAccount, u.Email, nil)
		if err := tx.Get(aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to get account for %s", u.Email)
		}
		account.APIKey = k
		if _, err := tx.Put(aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to put account for %s", u.Email)
		}

		apiKey = k
		return nil
	})

	if err != nil {
		log.Errorf("%v", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(apiKey)
	if err != nil {
		log.Errorf("%v", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (s *server) setPrivacyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	u, err := s.currentUser(r)
	if err == ErrNoUser {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	privacy, err := strconv.ParseBool(r.FormValue("privacy"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = s.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var account Account
		aKey := datastore.NameKey(KindAccount, u.Email, nil)
		if err := tx.Get(aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to get account for %s", u.Email)
		}
		account.Private = privacy
		if _, err := tx.Put(aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to put account for %s", u.Email)
		}
		return nil
	})

	if err != nil {
		log.Errorf("%v", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// Callers provide a transaction.
func setAPIKey(tx *datastore.Transaction, generator func() (string, error)) (string, error) {
	const maxTries = 10
	tries := 0

	for {
		tries++
		gen, err := generator()
		if err != nil {
			return "", errors.Wrapf(err, "failed to generate API key")
		}

		dsKey := datastore.NameKey(KindAPIKey, gen, nil)

		if err := tx.Get(dsKey, ptrStruct()); err != datastore.ErrNoSuchEntity {
			if err == nil {
				if tries == maxTries {
					return "", errors.New("API key already assigned")
				}
				continue
			}
			return "", errors.Wrapf(err, "failed to get API key")
		}
		if _, err := tx.Put(dsKey, ptrStruct()); err != nil {
			return "", errors.Wrapf(err, "failed to put API key")
		}
		return gen, nil
	}
}

func ptrStruct() *struct{} {
	return &struct{}{}
}

const privacyPolicy = `Your privacy is important to us. It is littleroot's policy to respect your
privacy regarding any information we may collect from you across our
website, https://` + AppDomain + `, and other sites we own and operate.

We only ask for personal information when we truly need it to provide a
service to you. We collect it by fair and lawful means, with your knowledge
and consent. We also let you know why we're collecting it and how it will be
used.

We only retain collected information for as long as necessary to provide you
with your requested service. What data we store, we'll protect within
commercially acceptable means to prevent loss and theft, as well as
unauthorized access, disclosure, copying, use or modification.

We don't share any personally identifying information publicly or with
third-parties, except when required to by law.

Our website may link to external sites that are not operated by us. Please
be aware that we have no control over the content and practices of these
sites, and cannot accept responsibility or liability for their respective
privacy policies.

You are free to refuse our request for your personal information, with the
understanding that we may be unable to provide you with some of your desired
services.

Your continued use of our website will be regarded as acceptance of our
practices around privacy and personal information. If you have any questions
about how we handle user data and personal information, feel free to contact
us.

This policy is effective as of 12 April 2020.

[Privacy Policy created with GetTerms.](https://getterms.io/)
`

func (s *server) privacyPolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	io.WriteString(w, privacyPolicy)
}
