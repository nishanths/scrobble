WEBPACK := npx webpack-cli
ESLINT_BASE := npx eslint . --ext .ts,.tsx

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

.PHONY: fmt
fmt:
	# tsfmt
	find . -type f \( -name '*.ts' -o -name '*.tsx' \) \
		! -name '*.d.ts' \
		! -regex '.*/node_modules/.*' \
		! -regex '.*/static/.*' | \
		xargs npx tsfmt --replace --
	# lint
	$(ESLINT_BASE)

.PHONY: fix
fix:
	$(ESLINT_BASE) --fix

.PHONY: sb
sb:
	npx start-storybook -p 9001 -c .storybook \
		-s static # https://storybook.js.org/docs/configurations/serving-static-files/

.PHONY: deps
deps:
	# NOTE: Needs node version 14; node version >=16 is known to fail.
	# Typically sass fails with the following error:
	#
	#  Error: Node Sass does not yet support your current environment: OS X
	#  64-bit with Unsupported runtime (102)
	#  ... https://github.com/sass/node-sass/releases/tag/v4.14.1
	#
	# To fix:
	#   make sure you have nodenv init-ed in your shell
	#   if necessary, re-open a new terminal tab, so that nodenv picks up the .node-version file in this repo.
	npm i
	# remove bad types file
	rm node_modules/flexsearch/index.d.ts
