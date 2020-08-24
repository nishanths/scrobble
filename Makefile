PWD               := $(shell pwd)
PROJECT_ID        := selective-scrobble
INDEX_YAML        := appengine/index.yaml
APP_YAML          := appengine/app.yaml

# NOTE: To deploy, typically you want `make all` followed by `make deploy`.
#
# For local development, run `make dev` and the web directory's `make dev`.
# You may optionally need `make bootstrap`.
#
# Only certain paths may be supported with fake data in local dev
# (see appengine/main.go).

default:
	@echo "the default target does nothing!"

.PHONY: all
all: clean bootstrap doc bindata go web

.PHONY: indexes
indexes:
	gcloud --quiet --project $(PROJECT_ID) datastore indexes create $(INDEX_YAML)

.PHONY: deploy
deploy:
	gcloud --quiet --project $(PROJECT_ID) app deploy -v 1 $(APP_YAML)
	echo $(shell git rev-list HEAD | head -n 1) > deployed-hash.txt

.PHONY: bindata
bindata:
	go-bindata -pkg=main -o=appengine/gen-bindata.go appengine/template

.PHONY: dev-deps
dev-deps:
	go get github.com/kevinburke/go-bindata/...
	go get gopkg.in/russross/blackfriday.v2

.PHONY: go
go:
	@go version
	cd appengine && go build -mod=vendor -o=main

.PHONY: web
web:
	@cd web && $(MAKE) dist

.PHONY: ln
ln:
	mkdir -p appengine/web
	ln -sfn $(PWD)/web/dist $(PWD)/appengine/web
	ln -sfn $(PWD)/web/static $(PWD)/appengine/web
	ln -sfn $(PWD)/doc/content $(PWD)/appengine/doccontent

.PHONY: clean
clean:
	@cd web && $(MAKE) clean
	rm -f appengine/main
	rm -rf appengine/web

.PHONY: bootstrap
bootstrap: ln _bootstrap

.PHONY: _bootstrap
_bootstrap:
	mkdir -p web/dist

.PHONY: dev
dev: doc bindata go
	cd appengine && ./main

.PHONY: test
test: go-test

.PHONY: go-test
go-test:
	cd appengine && go test -mod=vendor -race ./...

.PHONY: doc
doc:
	cd doc && $(MAKE) all
