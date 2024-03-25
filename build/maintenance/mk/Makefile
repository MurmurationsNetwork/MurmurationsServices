docker-build-maintenance:
	docker build -f build/maintenance/docker/Dockerfile \
		-t murmurations/$(DOCKER_TAG_PREFIX)maintenance .

docker-tag-maintenance: check-clean docker-build-maintenance
	docker tag murmurations/$(DOCKER_TAG_PREFIX)maintenance \
		murmurations/$(DOCKER_TAG_PREFIX)maintenance:${TAG}

docker-push-maintenance: docker-tag-maintenance
	docker push murmurations/$(DOCKER_TAG_PREFIX)maintenance:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)maintenance:$(TAG)