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
# (see appengine/server.go).

default:
	@echo "the default target does nothing!"

.PHONY: all
all: clean bootstrap bindata go web

.PHONY: indexes
indexes:
	gcloud --quiet --project $(PROJECT_ID) datastore indexes create $(INDEX_YAML)

.PHONY: deploy
deploy:
	gcloud --quiet --project $(PROJECT_ID) app deploy -v 1 $(APP_YAML)

.PHONY: bindata
bindata:
	go-bindata -pkg=main -o=appengine/gen-template.go appengine/template

.PHONY: dev-deps
dev-deps:
	go get github.com/jteeuwen/go-bindata/...

.PHONY: check-go
check-go:
	@go version
	cd appengine && go build -mod=vendor -o=/dev/null # to check successful compilation

.PHONY: go
go:
	@go version
	cd appengine && go build -mod=vendor -o=main

.PHONY: web
web:
	@cd web && $(MAKE) dist

.PHONY: ln-web
ln-web:
	mkdir -p appengine/web
	ln -s $(PWD)/web/dist $(PWD)/appengine/web
	ln -s $(PWD)/web/static $(PWD)/appengine/web
	sed -i '' 's/#29d/rgb\(21,126,251\)/g' $(PWD)/web/node_modules/nprogress/nprogress.css
	ln -s $(PWD)/web/node_modules/nprogress/nprogress.css $(PWD)/web/static/css/nprogress.css
	ln -s $(PWD)/web/node_modules/nprogress/nprogress.js $(PWD)/web/static/js/nprogress.js

.PHONY: clean
clean:
	@cd web && $(MAKE) clean
	rm -f appengine/main
	rm -rf appengine/web
	rm -rf web/static/{css/nprogress.css,js/nprogress.js}

.PHONY: bootstrap
bootstrap: ln-web _bootstrap

.PHONY: _bootstrap
_bootstrap:
	mkdir -p web/dist

.PHONY: dev
dev: go
	cd appengine && ./main
