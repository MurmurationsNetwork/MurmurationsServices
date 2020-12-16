test:
	export ENV=test && go test ./...

# ---------------------------------------------------------------

dev:
	export SOURCEPATH=$(PWD) && skaffold dev --port-forward

# ---------------------------------------------------------------

docker-build-index:
	$(MAKE) -C services/index/ docker-build

docker-build-validation:
	$(MAKE) -C services/validation/ docker-build

docker-build-library:
	$(MAKE) -C services/library/ docker-build

docker-build-nodecleaner:
	$(MAKE) -C services/cronjob/nodecleaner/ docker-build

docker-build-schemaparser:
	$(MAKE) -C services/cronjob/schemaparser/ docker-build

# ---------------------------------------------------------------

TAG ?= $(shell git rev-parse --short ${GITHUB_SHA})$(and $(shell git status -s),-dirty)

docker-tag-index: docker-build-index
	docker tag murmurations/index murmurations/index:${TAG}

docker-tag-validation: docker-build-validation
	docker tag murmurations/validation murmurations/validation:${TAG}

docker-tag-library: docker-build-library
	docker tag murmurations/library murmurations/library:${TAG}

docker-tag-nodecleaner: docker-build-nodecleaner
	docker tag murmurations/nodecleaner murmurations/nodecleaner:${TAG}

docker-tag-schemaparser: docker-build-schemaparser
	docker tag murmurations/schemaparser murmurations/schemaparser:${TAG}

# ---------------------------------------------------------------

docker-push-index: docker-tag-index
	docker push murmurations/index:latest
	docker push murmurations/index:$(TAG)

docker-push-validation: docker-tag-validation
	docker push murmurations/validation:latest
	docker push murmurations/validation:$(TAG)

docker-push-library: docker-tag-library
	docker push murmurations/library:latest
	docker push murmurations/library:$(TAG)

docker-push-nodecleaner: docker-tag-nodecleaner
	docker push murmurations/nodecleaner:latest
	docker push murmurations/nodecleaner:$(TAG)

docker-push-schemaparser: docker-tag-schemaparser
	docker push murmurations/schemaparser:latest
	docker push murmurations/schemaparser:$(TAG)

# ---------------------------------------------------------------

helm-staging-core:
	helm upgrade murmdev-core ./charts/murmdev/charts/core --install --wait --atomic

helm-staging-index:
	helm upgrade murmdev-index ./charts/murmdev/charts/index --set image=murmurations/index:$(TAG) --install --wait --atomic

helm-staging-validation:
	helm upgrade murmdev-validation ./charts/murmdev/charts/validation --set image=murmurations/validation:$(TAG) --install --wait --atomic

helm-staging-library:
	helm upgrade murmdev-library ./charts/murmdev/charts/library --set image=murmurations/library:$(TAG) --install --wait --atomic

helm-staging-nodecleaner:
	helm upgrade murmdev-nodecleaner ./charts/murmdev/charts/nodecleaner --set image=murmurations/nodecleaner:$(TAG) --install --wait --atomic

helm-staging-schemaparser:
	helm upgrade murmdev-schemaparser ./charts/murmdev/charts/schemaparser --set image=murmurations/schemaparser:$(TAG) --install --wait --atomic
