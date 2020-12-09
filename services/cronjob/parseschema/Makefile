# If TAG is not provided set default value
TAG ?= murmurations/parseschema:$(shell git rev-parse --short HEAD)$(and $(shell git status -s),-dirty)

docker-build:
	cd ../../../ && \
	docker build -f services/cronjob/parseschema/docker/Dockerfile-prod -t $(TAG) .
