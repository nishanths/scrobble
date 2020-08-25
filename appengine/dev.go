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
	if !validRootPath(r.URL.Path) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	const loggedIn = false

	if loggedIn {
		var loggedInArgs = RootArgs{
			Title: "Dashboard · Dev Scrobble",
			Bootstrap: BootstrapArgs{
				Host:             r.Host,
				Email:            "localdev@gmail.com",
				LogoutURL:        "/fakelogouturl",
				Account:          devFakeAccount,
				TotalSongs:       1337,
				LastScrobbleTime: time.Now().Unix(),
			},
			AppDomain: AppDomain,
		}

		if err := dashboardTmpl.Execute(w, loggedInArgs); err != nil {
			log.Errorf("failed to execute template: %v", err.Error())
		}
	} else {
		var loggedOutArgs = RootArgs{
			Title: "Apple Music scrobbling — Dev Scrobble",
			Bootstrap: BootstrapArgs{
				Host:     r.Host,
				LoginURL: "/fakeloginurl",
			},
			AppDomain: AppDomain,
		}

		if err := homeTmpl.Execute(w, loggedOutArgs); err != nil {
			log.Errorf("failed to execute template: %v", err.Error())
		}
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

func devScrobbledColorHandler(w http.ResponseWriter, r *http.Request) {
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

func devSongPlayCountHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)

	data := func() string {
		b, err := ioutil.ReadFile(filepath.Join(".devdata", "data/songs_play_count.json"))
		if err != nil {
			panic(err)
		}
		return string(b)
	}()

	io.WriteString(w, data)
}

func devSongLengthHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)

	data := func() string {
		b, err := ioutil.ReadFile(filepath.Join(".devdata", "data/songs_length.json"))
		if err != nil {
			panic(err)
		}
		return string(b)
	}()

	io.WriteString(w, data)
}

func devArtistPlayCountHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)

	data := func() string {
		b, err := ioutil.ReadFile(filepath.Join(".devdata", "data/artists_play_count.json"))
		if err != nil {
			panic(err)
		}
		return string(b)
	}()

	io.WriteString(w, data)
}

func devArtistAddedHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)

	data := func() string {
		b, err := ioutil.ReadFile(filepath.Join(".devdata", "data/artists_added.json"))
		if err != nil {
			panic(err)
		}
		return string(b)
	}()

	io.WriteString(w, data)
}
