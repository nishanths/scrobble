package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

var (
	indexTmpl = template.Must(template.New("").Parse(string(MustAsset("appengine/template/index.html"))))
	uTmpl     = template.Must(template.New("").Parse(string(MustAsset("appengine/template/u.html"))))
)

var defaultTxOpts = &datastore.TransactionOptions{XG: true}

type BootstrapArgs struct {
	Host      string  `json:"host"`
	Email     string  `json:"email"`
	LoginURL  string  `json:"loginURL"`
	LogoutURL string  `json:"logoutURL"`
	Account   Account `json:"account"`
}

type IndexArgs struct {
	Title     string
	Bootstrap BootstrapArgs
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

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
	exec := func(a IndexArgs) {
		if err := indexTmpl.Execute(w, a); err != nil {
			log.Errorf(ctx, "failed to execute template: %v", err.Error())
		}
	}

	u := user.Current(ctx)

	if u == nil {
		login, err := user.LoginURL(ctx, dest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		exec(IndexArgs{
			Title: title,
			Bootstrap: BootstrapArgs{
				Host:     host,
				LoginURL: login,
			},
		})
		return
	}

	logout, err := user.LogoutURL(ctx, dest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a, err := ensureAccount(ctx, u.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	exec(IndexArgs{
		Title: title,
		Bootstrap: BootstrapArgs{
			Host:      host,
			Email:     u.Email,
			LogoutURL: logout,
			Account:   a,
		},
	})
}

func ensureAccount(ctx context.Context, email string) (Account, error) {
	var account Account

	err := datastore.RunInTransaction(ctx, func(tx context.Context) error {
		aKey := datastore.NewKey(tx, KindAccount, email, 0, nil)

		if err := datastore.Get(tx, aKey, &account); err != nil {
			if err == datastore.ErrNoSuchEntity {
				// account entity does not exists; create new account entity
				if _, err := datastore.Put(tx, aKey, &account); err != nil {
					return errors.Wrapf(err, "failed to put account for %s", email)
				}
				return nil
			}

			// generic error
			return errors.Wrapf(err, "failed to get account for %s", email)
		}

		// account entity exists
		return nil
	}, defaultTxOpts)

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

func uHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	c := pathComponents(r.URL.Path)
	if len(c) != 2 { // 'u', username
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profileUsername := c[1]
	acc, _, ok := fetchAccountForUsername(ctx, profileUsername, w)
	if !ok {
		return
	}

	u := user.Current(ctx)

	// If the user is logged in, gather a logout URL and the account info.
	var logoutURL string
	var account Account
	self := false

	if u != nil {
		var err error
		logoutURL, err = user.LogoutURL(ctx, r.RequestURI)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := datastore.Get(ctx, datastore.NewKey(ctx, KindAccount, u.Email, 0, nil), &account); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		self = account.Username != "" && account.Username == profileUsername
	}

	bucketName, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get default GCS bucket name: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := uTmpl.Execute(w, UArgs{
		Title:           profileUsername,
		Host:            r.Host,
		ArtworkBaseURL:  "https://storage.googleapis.com/" + bucketName + "/" + artworkStorageDirectory,
		ProfileUsername: profileUsername,
		LogoutURL:       logoutURL,
		Account:         account,
		Self:            self,
		Private:         acc.Private,
	}); err != nil {
		log.Errorf(ctx, "failed to execute template: %v", err.Error())
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
func initializeAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	u := user.Current(ctx)
	if u == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	if ok := isAllowedUsername(username); !ok {
		w.WriteHeader(http.StatusNotAcceptable) // gross, but whatever
		return
	}

	inUse := false
	var account Account
	err := datastore.RunInTransaction(ctx, func(tx context.Context) error {
		// Ensure username uniqueness.
		uKey := datastore.NewKey(tx, KindUsername, username, 0, nil)
		if err := datastore.Get(tx, uKey, ptrStruct()); err != datastore.ErrNoSuchEntity {
			if err == nil {
				inUse = true
				return errors.New("username already in use")
			}
			return errors.Wrapf(err, "failed to get username")
		}
		if _, err := datastore.Put(tx, uKey, ptrStruct()); err != nil {
			return errors.Wrapf(err, "failed to put username")
		}

		// Ensure API key uniqueness.
		gen, err := setAPIKey(tx, generateAPIKey)
		if err != nil {
			return errors.Wrapf(err, "failed to set API key")
		}

		// Initialize the account.
		aKey := datastore.NewKey(tx, KindAccount, u.Email, 0, nil)
		if err := datastore.Get(tx, aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to get account for %s", u.Email)
		}
		if account.Username != "" {
			return fmt.Errorf("username already set for %s", u.Email)
		}
		account.Username = username
		account.APIKey = gen
		account.Private = true
		if _, err := datastore.Put(tx, aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to put account for %s", u.Email)
		}

		return nil
	}, defaultTxOpts)

	if err != nil {
		log.Errorf(ctx, "%v", err.Error())
		if inUse {
			w.WriteHeader(http.StatusNotAcceptable) // gross, but whatever
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	b, err := json.Marshal(account)
	if err != nil {
		log.Errorf(ctx, "failed to json-marshal account: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func newAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	u := user.Current(ctx)
	if u == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var apiKey string
	err := datastore.RunInTransaction(ctx, func(tx context.Context) error {
		k, err := setAPIKey(ctx, generateAPIKey)
		if err != nil {
			return err
		}

		var account Account
		aKey := datastore.NewKey(tx, KindAccount, u.Email, 0, nil)
		if err := datastore.Get(tx, aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to get account for %s", u.Email)
		}
		account.APIKey = k
		if _, err := datastore.Put(tx, aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to put account for %s", u.Email)
		}

		apiKey = k
		return nil
	}, defaultTxOpts)

	if err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(apiKey)
	if err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func setPrivacyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	u := user.Current(ctx)
	if u == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	privacy, err := strconv.ParseBool(r.FormValue("privacy"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = datastore.RunInTransaction(ctx, func(tx context.Context) error {
		var account Account
		aKey := datastore.NewKey(tx, KindAccount, u.Email, 0, nil)
		if err := datastore.Get(tx, aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to get account for %s", u.Email)
		}
		account.Private = privacy
		if _, err := datastore.Put(tx, aKey, &account); err != nil {
			return errors.Wrapf(err, "failed to put account for %s", u.Email)
		}
		return nil
	}, defaultTxOpts)

	if err != nil {
		log.Errorf(ctx, "%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Callers may provide a transaction context if they wish. The operations
// performed by setAPIKey are safe to do in a transaction.
func setAPIKey(ctx context.Context, generator func() (string, error)) (string, error) {
	const maxTries = 10
	tries := 0

	for {
		tries++
		gen, err := generator()
		if err != nil {
			return "", errors.Wrapf(err, "failed to generate API key")
		}

		dsKey := datastore.NewKey(ctx, KindAPIKey, gen, 0, nil)

		if err := datastore.Get(ctx, dsKey, ptrStruct()); err != datastore.ErrNoSuchEntity {
			if err == nil {
				if tries == maxTries {
					return "", errors.New("API key already assigned")
				}
				continue
			}
			return "", errors.Wrapf(err, "failed to get API key")
		}
		if _, err := datastore.Put(ctx, dsKey, ptrStruct()); err != nil {
			return "", errors.Wrapf(err, "failed to put API key")
		}
		return gen, nil
	}
}

func ptrStruct() *struct{} {
	return &struct{}{}
}
