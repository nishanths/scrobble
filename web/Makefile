WEBPACK := npx webpack-cli
TSFMT := npx tsfmt

default:
	@echo "the default target does nothing!"

.PHONY: dev
dev: sass
	$(WEBPACK) --mode=development --watch

.PHONY: dist
dist: sass
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
		xargs $(TSFMT) --replace --
	# lint
	npx eslint . --ext .ts,.tsx

.PHONY: sb
sb:
	npx start-storybook -p 9001 -c .storybook \
		-s static # https://storybook.js.org/docs/configurations/serving-static-files/

.PHONY: sass
sass:
	# content page wants same styles as root page, so need a compiled copy
	# of the css
	npx sass --embed-source-map scss/root.scss static/css/root.css
