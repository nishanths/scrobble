.PHONY: default
default:
	@echo "the default target does nothing!"

.PHONY: clean
clean:
	rm -rf content/api/dst content/guide/dst

.PHONY: build
build: clean go
	./cmd/docgen \
		-template=cmd/page.template.html \
		-src=content/api/src \
		-dst=content/api/dst \
		-indexTitle="Scrobble · API" \
		-titlePrefix="/api/v1/" \
		-linkPrefix="/doc/api/v1"

	./cmd/docgen \
		-template=cmd/page.template.html \
		-src=content/guide/src \
		-dst=content/guide/dst \
		-indexTitle="Scrobble · Guide" \

.PHONY: go
go:
	cd cmd && go build -o docgen .
