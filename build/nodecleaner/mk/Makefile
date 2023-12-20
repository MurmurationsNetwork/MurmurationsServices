docker-build-nodecleaner:
	docker build -f build/nodecleaner/docker/Dockerfile \
		-t murmurations/$(DOCKER_TAG_PREFIX)nodecleaner .

docker-tag-nodecleaner: check-clean docker-build-nodecleaner
	docker tag murmurations/$(DOCKER_TAG_PREFIX)nodecleaner \
		murmurations/$(DOCKER_TAG_PREFIX)nodecleaner:${TAG}

docker-push-nodecleaner: docker-tag-nodecleaner
	docker push murmurations/$(DOCKER_TAG_PREFIX)nodecleaner:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)nodecleaner:$(TAG)
