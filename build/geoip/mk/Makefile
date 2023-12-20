docker-build-geoip:
	docker build -f build/geoip/docker/Dockerfile \
		-t murmurations/$(DOCKER_TAG_PREFIX)geoip .

docker-tag-geoip: check-clean docker-build-geoip
	docker tag murmurations/$(DOCKER_TAG_PREFIX)geoip \
		murmurations/$(DOCKER_TAG_PREFIX)geoip:${TAG}

docker-push-geoip: docker-tag-geoip
	docker push murmurations/$(DOCKER_TAG_PREFIX)geoip:latest
	docker push murmurations/$(DOCKER_TAG_PREFIX)geoip:$(TAG)
