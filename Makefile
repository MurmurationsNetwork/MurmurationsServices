test:
	export ENV=test && go test ./...

# ---------------------------------------------------------------

dev:
	export SOURCEPATH=$(PWD) && skaffold dev --port-forward

# Docker Build

index:
	$(MAKE) -C services/index/ docker-build

validation:
	$(MAKE) -C services/validation/ docker-build

library:
	$(MAKE) -C services/library/ docker-build

nodecleanup:
	$(MAKE) -C services/cronjob/nodecleanup/ docker-build

parseschema:
	$(MAKE) -C services/cronjob/parseschema/ docker-build

# Docker Publish

VERSION ?= $(shell git rev-parse --short HEAD)

publish-index:
	@echo 'publish index:$(VERSION)'
	docker push murmurations/index:$(VERSION)
