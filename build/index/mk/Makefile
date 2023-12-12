docker-build-index:
	docker build -f build/index/docker/Dockerfile-prod \
		-t murmurations/$(DOCKER_TAG_PREFIX)index .

docker-tag-index: check-clean docker-build-index
	docker tag murmurations/$(DOCKER_TAG_PREFIX)index \
		murmurations/$(DOCKER_TAG_PREFIX)index:$(TAG)

docker-push-index: docker-tag-index
	docker push murmurations/$(DOCKER_TAG_PREFIX)index:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)index:$(TAG)
