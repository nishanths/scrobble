package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

const host = "scrobble.allele.cc"

var (
	indexTmpl = template.Must(template.New("").Parse(string(MustAsset("appengine/template/index.html"))))
	uTmpl     = template.Must(template.New("").Parse(string(MustAsset("appengine/template/u.html"))))
)

type BootstrapArgs struct {
	Host        string  `json:"host"`
	Email       string  `json:"email"`
	LoginURL    string  `json:"loginURL"`
	LogoutURL   string  `json:"logoutURL"`
	DownloadURL string  `json:"downloadURL"`
	Account     Account `json:"account"`
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
	dest := "https://" + host + r.RequestURI
	title := "Scrobble"
	download := "TODO"

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
				Host:        host,
				LoginURL:    login,
				DownloadURL: download,
			},
		})
		return
	}

	logout, err := user.LogoutURL(ctx, dest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var a Account
	if err := datastore.Get(ctx, datastore.NewKey(ctx, KindAccount, u.Email, 0, nil), &a); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	exec(IndexArgs{
		Title: title,
		Bootstrap: BootstrapArgs{
			Host:        host,
			Email:       u.Email,
			LogoutURL:   logout,
			DownloadURL: download,
			Account:     a,
		},
	})
}

type UArgs struct {
	Title           string  `json:"title"`
	ProfileUsername string  `json:"profileUsername"`
	LogoutURL       string  `json:"logoutURL"`
	Account         Account `json:"account"`
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

	profileAcc, profileAccID, ok := fetchAccountForUsername(ctx, profileUsername, w)
	if !ok {
		return
	}

	u := user.Current(ctx)

	// Don't allow accessing private accounts.
	if profileAcc.Private && !canViewProfile(profileAccID, u) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If the user is logged in, gather a logout URL and the account info.
	var logoutURL string
	var account Account
	if u != nil {
		var err error
		logoutURL, err = user.LogoutURL(ctx, "https://"+r.URL.Host+r.RequestURI)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := datastore.Get(ctx, datastore.NewKey(ctx, KindAccount, u.Email, 0, nil), &account); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err := uTmpl.Execute(w, UArgs{
		Title:           profileUsername,
		ProfileUsername: profileUsername,
		LogoutURL:       logoutURL,
		Account:         account,
	}); err != nil {
		log.Errorf(ctx, "failed to execute template: %v", err.Error())
	}
}

func canViewProfile(profileAccID string, u *user.User) bool {
	return u != nil || u.Email == profileAccID
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

func setUsernameHandler(w http.ResponseWriter, r *http.Request) {
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
	err := datastore.RunInTransaction(ctx, func(tx context.Context) error {
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

		var a Account
		aKey := datastore.NewKey(tx, KindAccount, u.Email, 0, nil)
		if err := datastore.Get(tx, aKey, &a); err != nil {
			return errors.Wrapf(err, "failed to get account for %s", u.Email)
		}
		if a.Username != "" {
			return fmt.Errorf("username already set for %s", u.Email)
		}
		a.Username = username
		if _, err := datastore.Put(tx, aKey, &a); err != nil {
			return errors.Wrapf(err, "failed to put account for %s", u.Email)
		}

		if _, err := setAPIKey(tx, generateAPIKey); err != nil {
			return errors.Wrapf(err, "failed to set API key")
		}

		return nil
	}, nil)

	if err != nil {
		if inUse {
			w.WriteHeader(http.StatusNotAcceptable) // gross, but whatever
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func newAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	k, err := setAPIKey(ctx, generateAPIKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.WriteString(w, k)
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
