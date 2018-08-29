WEBPACK := node_modules/.bin/webpack-cli --no-color

default: dev

.PHONY: dev
dev:
	$(WEBPACK) --mode=development --watch

.PHONY: dist
dist:
	$(WEBPACK) --mode=production

.PHONY: clean
clean:
	rm -rf dist
