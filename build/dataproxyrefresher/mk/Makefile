docker-build-dataproxyrefresher:
	docker build -f build/dataproxyrefresher/docker/Dockerfile-prod \
		-t murmurations/$(DOCKER_TAG_PREFIX)dataproxyrefresher .

docker-tag-dataproxyrefresher: check-clean docker-build-dataproxyrefresher
	docker tag murmurations/$(DOCKER_TAG_PREFIX)dataproxyrefresher \
		murmurations/$(DOCKER_TAG_PREFIX)dataproxyrefresher:${TAG}

docker-push-dataproxyrefresher: docker-tag-dataproxyrefresher
	docker push murmurations/$(DOCKER_TAG_PREFIX)dataproxyrefresher:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)dataproxyrefresher:$(TAG)
