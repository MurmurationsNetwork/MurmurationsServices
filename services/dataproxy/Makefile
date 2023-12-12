DOCKER_TAG_PREFIX := $(if $(filter production,$(DEPLOY_ENV)),,${DEPLOY_ENV}-)

docker-build:
	cd ../../ && \
	docker build -f services/dataproxy/docker/Dockerfile-prod -t murmurations/$(DOCKER_TAG_PREFIX)dataproxy .
