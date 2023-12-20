docker-build-library:
	docker build -f build/library/docker/Dockerfile \
		-t murmurations/$(DOCKER_TAG_PREFIX)library .

docker-tag-library: check-clean docker-build-library
	docker tag murmurations/$(DOCKER_TAG_PREFIX)library \
		murmurations/$(DOCKER_TAG_PREFIX)library:${TAG}

docker-push-library: docker-tag-library
	docker push murmurations/$(DOCKER_TAG_PREFIX)library:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)library:$(TAG)
