# If TAG is not provided set default value
TAG ?= murmurations/validation:$(shell git rev-parse --short ${GITHUB_SHA})$(and $(shell git status -s),-dirty)

docker-build:
	cd ../../ && \
	docker build -f services/validation/docker/Dockerfile-prod -t $(TAG) .
