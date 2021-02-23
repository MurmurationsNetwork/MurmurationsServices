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

docker-build-revalidatenode:
	$(MAKE) -C services/cronjob/revalidatenode/ docker-build

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

docker-tag-revalidatenode: docker-build-revalidatenode
	docker tag murmurations/revalidatenode murmurations/revalidatenode:${TAG}

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

docker-push-revalidatenode: docker-tag-revalidatenode
	docker push murmurations/revalidatenode:latest
	docker push murmurations/revalidatenode:$(TAG)

# ---------------------------------------------------------------

helm-staging-ingress:
	helm upgrade murmurations-ingress ./charts/murmurations/charts/ingress --set global.env=staging --install --wait --atomic

helm-staging-index:
	helm upgrade murmurations-index ./charts/murmurations/charts/index --set global.env=staging,image=murmurations/index:$(TAG) --install --wait --atomic

helm-staging-validation:
	helm upgrade murmurations-validation ./charts/murmurations/charts/validation --set global.env=staging,image=murmurations/validation:$(TAG) --install --wait --atomic

helm-staging-library:
	helm upgrade murmurations-library ./charts/murmurations/charts/library --set global.env=staging,image=murmurations/library:$(TAG) --install --wait --atomic

helm-staging-nodecleaner:
	helm upgrade murmurations-nodecleaner ./charts/murmurations/charts/nodecleaner --set global.env=staging,image=murmurations/nodecleaner:$(TAG) --install --wait --atomic

helm-staging-schemaparser:
	helm upgrade murmurations-schemaparser ./charts/murmurations/charts/schemaparser --set global.env=staging,image=murmurations/schemaparser:$(TAG) --install --wait --atomic

helm-staging-revalidatenode:
	helm upgrade murmurations-revalidatenode ./charts/murmurations/charts/revalidatenode --set global.env=staging,image=murmurations/revalidatenode:$(TAG) --install --wait --atomic

# ---------------------------------------------------------------

# Please make sure `latest` is the version you want to deploy to production,
# otherwise please update to the correct version.
SPECIFIC_TAG ?= latest

helm-production-ingress:
	helm upgrade murmurations-ingress ./charts/murmurations/charts/ingress --set global.env=production --install --wait --atomic

helm-production-index:
	helm upgrade murmurations-index ./charts/murmurations/charts/index --set global.env=production,image=murmurations/index:$(SPECIFIC_TAG) --install --wait --atomic

helm-production-validation:
	helm upgrade murmurations-validation ./charts/murmurations/charts/validation --set global.env=production,image=murmurations/validation:$(SPECIFIC_TAG) --install --wait --atomic

helm-production-library:
	helm upgrade murmurations-library ./charts/murmurations/charts/library --set global.env=production,image=murmurations/library:$(SPECIFIC_TAG) --install --wait --atomic

helm-production-nodecleaner:
	helm upgrade murmurations-nodecleaner ./charts/murmurations/charts/nodecleaner --set global.env=production,image=murmurations/nodecleaner:$(SPECIFIC_TAG) --install --wait --atomic

helm-production-schemaparser:
	helm upgrade murmurations-schemaparser ./charts/murmurations/charts/schemaparser --set global.env=production,image=murmurations/schemaparser:$(SPECIFIC_TAG) --install --wait --atomic

helm-production-revalidatenode:
	helm upgrade murmurations-revalidatenode ./charts/murmurations/charts/revalidatenode --set global.env=production,image=murmurations/revalidatenode:$(SPECIFIC_TAG) --install --wait --atomic
