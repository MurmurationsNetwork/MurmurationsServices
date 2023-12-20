docker-build-revalidatenode:
	docker build -f build/revalidatenode/docker/Dockerfile \
		-t murmurations/$(DOCKER_TAG_PREFIX)revalidatenode .

docker-tag-revalidatenode: check-clean docker-build-revalidatenode
	docker tag murmurations/$(DOCKER_TAG_PREFIX)revalidatenode \
		murmurations/$(DOCKER_TAG_PREFIX)revalidatenode:${TAG}

docker-push-revalidatenode: docker-tag-revalidatenode
	docker push murmurations/$(DOCKER_TAG_PREFIX)revalidatenode:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)revalidatenode:$(TAG)
