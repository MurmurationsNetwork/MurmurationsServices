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

TAG ?= $(shell git rev-parse --short ${GITHUB_SHA})$(and $(shell git status -s),-dirty)

docker_tag: docker_build
	docker tag murmurations/index murmurations/index:${TAG}
	docker tag murmurations/validation murmurations/validation:${TAG}
	docker tag murmurations/library murmurations/library:${TAG}
	docker tag murmurations/nodecleanup murmurations/nodecleanup:${TAG}
	docker tag murmurations/parseschema murmurations/parseschema:${TAG}

docker_push: docker_tag
	docker push murmurations/index
	docker push murmurations/index:$(TAG)
	docker push murmurations/validation
	docker push murmurations/validation:$(TAG)
	docker push murmurations/library
	docker push murmurations/library:$(TAG)
	docker push murmurations/nodecleanup
	docker push murmurations/nodecleanup:$(TAG)
	docker push murmurations/parseschema
	docker push murmurations/parseschema:$(TAG)
