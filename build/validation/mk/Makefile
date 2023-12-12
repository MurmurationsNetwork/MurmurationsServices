docker-build-validation:
	docker build -f build/validation/docker/Dockerfile-prod \
		-t murmurations/$(DOCKER_TAG_PREFIX)validation .

docker-tag-validation: check-clean docker-build-validation
	docker tag murmurations/$(DOCKER_TAG_PREFIX)validation \
		murmurations/$(DOCKER_TAG_PREFIX)validation:${TAG}

docker-push-validation: docker-tag-validation
	docker push murmurations/$(DOCKER_TAG_PREFIX)validation:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)validation:$(TAG)
