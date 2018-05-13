DEV_APPSERVER     := dev_appserver.py
DEV_APPSERVER_CMD := $(DEV_APPSERVER) appengine/app.yaml --enable_watching_go_path=false
TSC_COMMAND       := tsc -p web/tsconfig.json
PWD               := $(shell pwd)

.PHONY: default
default: serve

.PHONY: serve
serve:
	$(DEV_APPSERVER_CMD)

.PHONY: watch
watch: ts scss
	$(TSC_COMMAND) --watch & \
	scss --watch web/scss/main.scss:web/dist/css/main.css

.PHONY: ts
ts:
	$(TSC_COMMAND)

.PHONY: scss
scss:
	mkdir -p web/dist/css
	scss web/scss/main.scss web/dist/css/main.css

.PHONY: ln
ln:
	mkdir -p $(PWD)/appengine/web
	ln -s $(PWD)/web/dist $(PWD)/appengine/web/dist
	ln -s $(PWD)/web/static $(PWD)/appengine/web/static

.PHONY: bindata
bindata:
	go-bindata -pkg=server -o=server/template.go web/template

.PHONY: go-deps
go-deps:
	go get github.com/jteeuwen/go-bindata/...
