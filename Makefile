.PHONY: deploy

PROJECT-ID := selective-scrobble

default:
	@printf "the default target does nothing"

deploy:
	gcloud app deploy --project $(PROJECT-ID) appengine/app.yaml
