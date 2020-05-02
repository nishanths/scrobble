WEBPACK := npx webpack-cli
TSFMT := npx tsfmt

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
	find . -type f \( -name '*.ts' -o -name '*.tsx' \) \
		! -name '*.d.ts' \
		! -regex '.*/node_modules/.*' | \
		xargs $(TSFMT) --replace --

.PHONY: sb
sb:
	npm run storybook
