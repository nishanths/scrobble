package main

import (
	"io"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func devUHandler(w http.ResponseWriter, r *http.Request) {
	const profileUsername = "devuser"
	const signedInUsername = "devuser2"
	const bucketName = "selective-scrobble.appspot.com"

	ctx := appengine.NewContext(r)

	if err := uTmpl.Execute(w, UArgs{
		Title:           profileUsername,
		Host:            r.Host,
		ArtworkBaseURL:  "https://storage.googleapis.com/" + bucketName + "/" + artworkStorageDirectory,
		ProfileUsername: profileUsername,
		LogoutURL:       "/fake",
		Account: Account{
			Username: signedInUsername,
			APIKey:   "FAKE",
			Private:  true,
		},
		Self: profileUsername == signedInUsername,
	}); err != nil {
		log.Errorf(ctx, "failed to execute template: %v", err.Error())
	}
}

func devScrobbledHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, devScrobbledResponse)
}

const devScrobbledResponse = `
[{"albumTitle":"Eulogy - EP","artistName":"Charles Fauna","title":"Friday","totalTime":237141000000,"year":2018,"sortAlbumTitle":"Eulogy - EP","sortArtistName":"Charles Fauna","sortTitle":"Friday","lastPlayed":1535526027,"playCount":139,"artworkHash":"6516638118222518123082712556622554402160984106"},{"albumTitle":"'Like a Woman, Like a Drunkard, Like an Animal' - EP","artistName":"Caroline Vreeland","title":"Black Summer","totalTime":243692000000,"year":2018,"sortAlbumTitle":"'Like a Woman, Like a Drunkard, Like an Animal' - EP","sortArtistName":"Caroline Vreeland","sortTitle":"Black Summer","lastPlayed":1535512128,"playCount":151,"artworkHash":"7057114158121718814692129141751131152492221198163"},{"albumTitle":"Life In Technicolor ii - Single","artistName":"Coldplay","title":"Life In Technicolor ii","totalTime":247293000000,"year":2009,"sortAlbumTitle":"Life In Technicolor ii - Single","sortArtistName":"Coldplay","sortTitle":"Life In Technicolor ii","lastPlayed":1535500457,"playCount":20,"artworkHash":"201451131782132424620852172003429183179247143240249101"},{"albumTitle":"La thune - Single","artistName":"Angèle","title":"La thune","totalTime":202160000000,"year":2018,"sortAlbumTitle":"La thune - Single","sortArtistName":"Angèle","sortTitle":"La thune","lastPlayed":1535481594,"playCount":11,"artworkHash":"11823232961308921650108482348719754244304192251"},{"albumTitle":"Something More Holy - EP","artistName":"Morly","title":"Plucky","totalTime":194274000000,"year":2016,"sortAlbumTitle":"Something More Holy - EP","sortArtistName":"Morly","sortTitle":"Plucky","lastPlayed":1535436609,"playCount":52,"artworkHash":"243891076875115341152081862297645798165132190206"},{"albumTitle":"Saltwater for Strings ((reimagined by Pêtr Aleksänder)) - Single","artistName":"Geowulf \u0026 Pêtr Aleksänder","title":"Saltwater for Strings ((reimagined by Pêtr Aleksänder))","totalTime":239582000000,"year":2018,"sortAlbumTitle":"Saltwater for Strings ((reimagined by Pêtr Aleksänder)) - Single","sortArtistName":"Geowulf \u0026 Pêtr Aleksänder","sortTitle":"Saltwater for Strings ((reimagined by Pêtr Aleksänder))","lastPlayed":1535430278,"playCount":14,"artworkHash":"85382919741376913625150938222019223250442155438"}]
`
