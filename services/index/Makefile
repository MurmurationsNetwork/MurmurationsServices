# If TAG is not provided set default value
TAG ?= murmurations/index:$(shell git rev-parse --short HEAD)$(and $(shell git status -s),-dirty)

docker-build:
	cd ../../ && \
	docker build -f services/index/docker/Dockerfile-prod -t $(TAG) .
