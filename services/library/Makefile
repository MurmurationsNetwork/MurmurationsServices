# If TAG is not provided set default value
TAG ?= murmurations/library:$(shell git rev-parse --short ${GITHUB_SHA})$(and $(shell git status -s),-dirty)

docker-build:
	cd ../../ && \
	docker build -f services/library/docker/Dockerfile-prod -t $(TAG) .
