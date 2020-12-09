test:
	export ENV=test && go test ./...

# ---------------------------------------------------------------

dev:
	export SOURCEPATH=$(PWD) && skaffold dev --port-forward

# ---------------------------------------------------------------

docker_build:
	$(MAKE) -C services/index/ docker-build
	$(MAKE) -C services/validation/ docker-build
	$(MAKE) -C services/library/ docker-build
	$(MAKE) -C services/cronjob/nodecleanup/ docker-build
	$(MAKE) -C services/cronjob/parseschema/ docker-build

TAG ?= $(shell git rev-parse --short ${GITHUB_SHA})

docker_push:
	docker push murmurations/index:latest
	docker push murmurations/index:$(TAG)
	docker push murmurations/validation:latest
	docker push murmurations/validation:$(TAG)
	docker push murmurations/library:latest
	docker push murmurations/library:$(TAG)
	docker push murmurations/nodecleanup:latest
	docker push murmurations/nodecleanup:$(TAG)
	docker push murmurations/parseschema:latest
	docker push murmurations/parseschema:$(TAG)
