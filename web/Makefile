WEBPACK := node_modules/webpack-cli/bin/webpack.js
DIST    := ../appengine/web/dist

.PHONY: default
default:
	mkdir -p $(DIST)
	cp index.html $(DIST)/index.html
	$(WEBPACK) --no-color --mode=development --watch
