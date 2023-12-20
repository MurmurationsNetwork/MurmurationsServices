docker-build-dataproxy:
	docker build -f build/dataproxy/docker/Dockerfile \
		-t murmurations/$(DOCKER_TAG_PREFIX)dataproxy .

docker-tag-dataproxy: check-clean docker-build-dataproxy
	docker tag murmurations/$(DOCKER_TAG_PREFIX)dataproxy \
		murmurations/$(DOCKER_TAG_PREFIX)dataproxy:${TAG}

docker-push-dataproxy: docker-tag-dataproxy
	docker push murmurations/$(DOCKER_TAG_PREFIX)dataproxy:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)dataproxy:$(TAG)
