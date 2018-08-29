PWD        := $(shell pwd)
PROJECT-ID := selective-scrobble

default:
	@echo "the default target does nothing!"

.PHONY: all
all: clean _bootstrap bindata ln-web

.PHONY: deploy
deploy:
	gcloud app deploy --project $(PROJECT-ID) appengine/app.yaml

.PHONY: bindata
bindata:
	go-bindata -pkg=main -o=appengine/template.go appengine/template

.PHONY: deps
deps:
	go get github.com/jteeuwen/go-bindata/...

.PHONY: ln-web
ln-web:
	mkdir -p appengine/web
	ln -s $(PWD)/web/dist $(PWD)/appengine/web
	ln -s $(PWD)/web/static $(PWD)/appengine/web

.PHONY: clean
clean:
	rm -rf web/dist
	rm -rf appengine/web

.PHONY: _bootstrap
_bootstrap:
	mkdir -p web/dist
