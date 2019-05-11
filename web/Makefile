WEBPACK := node_modules/.bin/webpack-cli

default:
	@echo "the default target does nothing!"

.PHONY: dev
dev:
	$(WEBPACK) --mode=development --watch

.PHONY: dist
dist:
	$(WEBPACK) --no-color --mode=production

.PHONY: clean
clean:
	rm -rf dist

.PHONY: bootstrap
	npm i
