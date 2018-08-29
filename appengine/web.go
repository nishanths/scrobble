package main

import (
	"html/template"
	"net/http"
)

var (
	indexTmpl = template.Must(template.New("").Parse(string(MustAsset("appengine/template/index.html"))))
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, struct{}{})
}
