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

// TODO: WTF is this monstrosity
var (
	rootTmpl = template.Must(
		template.Must(template.New("").Parse(string(MustAsset("appengine/template/fs-snippet.tmpl")))).
			Parse(string(MustAsset("appengine/template/root.tmpl"))),
	)
	uTmpl = template.Must(
		template.Must(template.New("").Parse(string(MustAsset("appengine/template/fs-snippet.tmpl")))).
			Parse(string(MustAsset("appengine/template/u.tmpl"))),
	)
)

type BootstrapArgs struct {
	Host      string  `json:"host"`
	Email     string  `json:"email"`
	LoginURL  string  `json:"loginURL"`
	LogoutURL string  `json:"logoutURL"`
	Account   Account `json:"account"`
}

type RootArgs struct {
	Title     string
	Bootstrap BootstrapArgs
}

func (s *server) rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	host := r.Host
	dest := r.RequestURI
	title := "Scrobble"

	// helper function
	exec := func(a RootArgs) {
		if err := rootTmpl.Execute(w, a); err != nil {
			log.Errorf("failed to execute template: %v", err.Error())
		}
	}

	u, err := s.currentUser(r)

	if err != nil {
		// either generic error or ErrNoUser
		login := loginURLWithRedirect(dest)
		exec(RootArgs{
			Title: title,
			Bootstrap: BootstrapArgs{
				Host:     host,
				LoginURL: login,
			},
		})
		return
	}

	logout := logoutURLWithRedirect(dest)

	a, err := ensureAccount(ctx, u.Email, s.ds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	exec(RootArgs{
		Title: title,
		Bootstrap: BootstrapArgs{
			Host:      host,
			Email:     u.Email,
			LogoutURL: logout,
			Account:   a,
		},
	})
}

func ensureAccount(ctx context.Context, email string, ds *datastore.Client) (Account, error) {
	var account Account

	_, err := ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		aKey := datastore.NameKey(KindAccount, email, nil)

		if err := tx.Get(aKey, &account); err != nil {
			if err == datastore.ErrNoSuchEntity {
				// account entity does not exists; create new account entity
				if _, err := tx.Put(aKey, &account); err != nil {
					return errors.Wrapf(err, "failed to put account for %s", email)
				}
				return nil
			}

			// generic error
			return errors.Wrapf(err, "failed to get account for %s", email)
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	c := pathComponents(r.URL.Path)
	if len(c) != 2 && len(c) != 3 { // /u/username[/loved]
		w.WriteHeader(http.StatusBadRequest)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		self = account.Username != "" && account.Username == profileUsername
	}

	if err := uTmpl.Execute(w, UArgs{
		Title:           profileUsername,
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	u, err := s.currentUser(r)
	if err == ErrNoUser {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	if ok := isAllowedUsername(username); !ok {
		w.WriteHeader(http.StatusNotAcceptable) // gross, but whatever
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
		account.Private = true
		if _, err := tx.Put(aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to put account for %s", u.Email)
		}

		return nil
	})

	if err != nil {
		log.Errorf("%v", err.Error())
		if inUse {
			w.WriteHeader(http.StatusNotAcceptable) // gross, but whatever
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	b, err := json.Marshal(account)
	if err != nil {
		log.Errorf("failed to json-marshal account: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (s *server) newAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	u, err := s.currentUser(r)
	if err == ErrNoUser {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(apiKey)
	if err != nil {
		log.Errorf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (s *server) setPrivacyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	u, err := s.currentUser(r)
	if err == ErrNoUser {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	privacy, err := strconv.ParseBool(r.FormValue("privacy"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

const privacyPolicy = `Privacy Policy
--------------

Your privacy is important to us. It is allele's policy to respect your
privacy regarding any information we may collect from you across our
website, https://scrobble.allele.cc, and other sites we own and operate.

We only ask for personal information when we truly need it to provide a
service to you. We collect it by fair and lawful means, with your knowledge
and consent. We also let you know why we’re collecting it and how it will be
used.

We only retain collected information for as long as necessary to provide you
with your requested service. What data we store, we’ll protect within
commercially acceptable means to prevent loss and theft, as well as
unauthorized access, disclosure, copying, use or modification.

We don’t share any personally identifying information publicly or with
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, privacyPolicy)
}
