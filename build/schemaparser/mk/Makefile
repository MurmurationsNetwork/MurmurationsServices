docker-build-schemaparser:
	docker build -f build/schemaparser/docker/Dockerfile \
		-t murmurations/$(DOCKER_TAG_PREFIX)schemaparser .

docker-tag-schemaparser: check-clean docker-build-schemaparser
	docker tag murmurations/$(DOCKER_TAG_PREFIX)schemaparser \
		murmurations/$(DOCKER_TAG_PREFIX)schemaparser:${TAG}

docker-push-schemaparser: docker-tag-schemaparser
	docker push murmurations/$(DOCKER_TAG_PREFIX)schemaparser:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)schemaparser:$(TAG)
