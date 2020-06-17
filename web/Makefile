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
	npm i
	# remove bad types file
	rm node_modules/flexsearch/index.d.ts
