package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/nishanths/scrobble/appengine/log"
)

const devSignedInUsername = "devuser"

var devFakeAccount = Account{
	Username: devSignedInUsername,
	APIKey:   "FAKEAPIKEY",
	Private:  false,
}

func devRootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	a := RootArgs{
		Title: "Dev Scrobble",
		Bootstrap: BootstrapArgs{
			Host:      r.Host,
			Email:     "localdev@gmail.com",
			LogoutURL: "/fakelogouturl",
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
		Title:           profileUsername + "'s scrobbles",
		Host:            r.Host,
		ArtworkBaseURL:  "https://storage.googleapis.com/" + bucketName + "/" + artworkStorageDirectory,
		ProfileUsername: profileUsername,
		LogoutURL:       "/fakelogouturl",
		Account:         devFakeAccount,
		Self:            profileUsername == devSignedInUsername,
	}); err != nil {
		log.Errorf("failed to execute template: %v", err.Error())
	}
}

func devScrobbledHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	file := "scrobbled_all.json"

	if r.FormValue("loved") == "true" {
		file = "scrobbled_loved.json" // jq "[.songs[] | select(.loved == true)]" scrobbled_all.json
	}

	if r.FormValue("song") != "" {
		file = "scrobbled_ident.json"
	}

	data := func() string {
		b, err := ioutil.ReadFile(filepath.Join(".devdata", file))
		if err != nil {
			panic(err)
		}
		return string(b)
	}()

	io.WriteString(w, data)
}

func devArtworkColorHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)

	data := func() string {
		b, err := ioutil.ReadFile(filepath.Join(".devdata", "scrobbled_color_white.json"))
		if err != nil {
			panic(err)
		}
		return string(b)
	}()

	io.WriteString(w, data)
}
