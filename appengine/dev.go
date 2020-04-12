package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/nishanths/scrobble/appengine/log"
)

const devSignedInUsername = "devuser"

var devFakeAccount = Account{
	Username: devSignedInUsername,
	APIKey:   "FAKE",
	Private:  true,
}

func devRootHandler(w http.ResponseWriter, r *http.Request) {
	a := RootArgs{
		Title: "Dev Scrobble",
		Bootstrap: BootstrapArgs{
			Host:      r.Host,
			Email:     "fake@gmail.com",
			LogoutURL: "/fake",
			Account:   devFakeAccount,
		},
	}

	if err := rootTmpl.Execute(w, a); err != nil {
		log.Errorf("failed to execute template: %v", err.Error())
	}
}

func devUHandler(w http.ResponseWriter, r *http.Request) {
	const profileUsername = devSignedInUsername
	const bucketName = "selective-scrobble.appspot.com"

	if err := uTmpl.Execute(w, UArgs{
		Title:           profileUsername,
		Host:            r.Host,
		ArtworkBaseURL:  "https://storage.googleapis.com/" + bucketName + "/" + artworkStorageDirectory,
		ProfileUsername: profileUsername,
		LogoutURL:       "/fake",
		Account:         devFakeAccount,
		Self:            profileUsername == devSignedInUsername,
	}); err != nil {
		log.Errorf("failed to execute template: %v", err.Error())
	}
}

func devScrobbledHandler(w http.ResponseWriter, r *http.Request) {
	var devScrobbledResponse = func() string {
		b, err := ioutil.ReadFile(filepath.Join(".devdata", "scrobbled.json"))
		if err != nil {
			panic(err)
		}
		return string(b)
	}()

	io.WriteString(w, devScrobbledResponse)
}
