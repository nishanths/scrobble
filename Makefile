PWD        := $(shell pwd)
PROJECT-ID := selective-scrobble

default:
	@echo "the default target does nothing!"

.PHONY: all
all: clean _bootstrap bindata web ln-web

.PHONY: deploy
deploy:
	gcloud --quiet --project $(PROJECT-ID) app deploy appengine/app.yaml

.PHONY: bindata
bindata:
	go-bindata -pkg=main -o=appengine/template.go appengine/template

.PHONY: deps
deps:
	go get github.com/jteeuwen/go-bindata/...

.PHONY: web
web:
	@cd web && $(MAKE) dist

.PHONY: ln-web
ln-web:
	mkdir -p appengine/web
	ln -s $(PWD)/web/dist $(PWD)/appengine/web
	ln -s $(PWD)/web/static $(PWD)/appengine/web

.PHONY: clean
clean:
	@cd web && $(MAKE) clean
	rm -rf appengine/web

.PHONY: _bootstrap
_bootstrap:
	mkdir -p web/dist
