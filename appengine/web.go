package main

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const host = "scrobble.allele.cc"

var (
	indexTmpl = template.Must(template.New("").Parse(string(MustAsset("appengine/template/index.html"))))
)

type BaseArgs struct {
	Title string
	Host  string
}

type IndexArgs struct {
	BaseArgs
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if err := indexTmpl.Execute(w, IndexArgs{
		BaseArgs{Title: "Scrobble", Host: host},
	}); err != nil {
		log.Errorf(ctx, "failed to execute template: %v", err.Error())
	}
}

func uHandler(w http.ResponseWriter, r *http.Request) {
}
