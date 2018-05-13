package server

import "google.golang.org/appengine/datastore"

func AccountsForToken(token string) *datastore.Query {
	return datastore.NewQuery(KindAccount).Filter("APIToken=", token)
}

func Playbacks() *datastore.Query {
	return datastore.NewQuery(KindPlayback)
}
