docker-build-dataproxyupdater:
	docker build -f build/dataproxyupdater/docker/Dockerfile \
		-t murmurations/$(DOCKER_TAG_PREFIX)dataproxyupdater .

docker-tag-dataproxyupdater: check-clean docker-build-dataproxyupdater
	docker tag murmurations/$(DOCKER_TAG_PREFIX)dataproxyupdater \
		murmurations/$(DOCKER_TAG_PREFIX)dataproxyupdater:${TAG}

docker-push-dataproxyupdater: docker-tag-dataproxyupdater
	docker push murmurations/$(DOCKER_TAG_PREFIX)dataproxyupdater:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)dataproxyupdater:$(TAG)
