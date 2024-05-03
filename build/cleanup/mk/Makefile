docker-build-cleanup:
	docker build -f build/cleanup/docker/Dockerfile \
		-t murmurations/$(DOCKER_TAG_PREFIX)cleanup .

docker-tag-cleanup: check-clean docker-build-cleanup
	docker tag murmurations/$(DOCKER_TAG_PREFIX)cleanup \
		murmurations/$(DOCKER_TAG_PREFIX)cleanup:${TAG}

docker-push-cleanup: docker-tag-cleanup
	docker push murmurations/$(DOCKER_TAG_PREFIX)cleanup:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)cleanup:$(TAG)
