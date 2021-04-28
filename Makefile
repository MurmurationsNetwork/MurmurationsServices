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

docker-build-geoip:
	$(MAKE) -C services/geoip/ docker-build

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

docker-tag-geoip: docker-build-geoip
	docker tag murmurations/geoip murmurations/geoip:${TAG}

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

docker-push-geoip: docker-tag-geoip
	docker push murmurations/geoip:latest
	docker push murmurations/geoip:$(TAG)

# ---------------------------------------------------------------

deploy-ingress:
	helm upgrade murmurations-ingress ./charts/murmurations/charts/ingress --set global.env=staging --install --wait --atomic

deploy-mq:
	helm upgrade murmurations-mq ./charts/murmurations/charts/message-queue --set global.env=staging --install --wait --atomic

deploy-index:
	helm upgrade murmurations-index ./charts/murmurations/charts/index --set global.env=staging,image=murmurations/index:$(TAG) --install --wait --atomic

deploy-validation:
	helm upgrade murmurations-validation ./charts/murmurations/charts/validation --set global.env=staging,image=murmurations/validation:$(TAG) --install --wait --atomic

deploy-library:
	helm upgrade murmurations-library ./charts/murmurations/charts/library --set global.env=staging,image=murmurations/library:$(TAG) --install --wait --atomic

deploy-nodecleaner:
	helm upgrade murmurations-nodecleaner ./charts/murmurations/charts/nodecleaner --set global.env=staging,image=murmurations/nodecleaner:$(TAG) --install --wait --atomic

deploy-schemaparser:
	helm upgrade murmurations-schemaparser ./charts/murmurations/charts/schemaparser --set global.env=staging,image=murmurations/schemaparser:$(TAG) --install --wait --atomic

deploy-revalidatenode:
	helm upgrade murmurations-revalidatenode ./charts/murmurations/charts/revalidatenode --set global.env=staging,image=murmurations/revalidatenode:$(TAG) --install --wait --atomic

deploy-geoip:
	helm upgrade murmurations-geoip ./charts/murmurations/charts/geoip --set global.env=staging,image=murmurations/geoip:$(TAG) --install --wait --atomic

# ---------------------------------------------------------------

# Please update to the version you want to deploy.
SPECIFIC_TAG ?= <>
ENV ?= <>

manually-deploy-ingress:
	helm upgrade murmurations-ingress ./charts/murmurations/charts/ingress --set global.env=$(ENV) --install --wait --atomic

manually-deploy-mq:
	helm upgrade murmurations-mq ./charts/murmurations/charts/message-queue --set global.env=$(ENV) --install --wait --atomic

manually-deploy-index:
	helm upgrade murmurations-index ./charts/murmurations/charts/index --set global.env=$(ENV),image=murmurations/index:$(SPECIFIC_TAG) --install --wait --atomic

manually-deploy-validation:
	helm upgrade murmurations-validation ./charts/murmurations/charts/validation --set global.env=$(ENV),image=murmurations/validation:$(SPECIFIC_TAG) --install --wait --atomic

manually-deploy-library:
	helm upgrade murmurations-library ./charts/murmurations/charts/library --set global.env=$(ENV),image=murmurations/library:$(SPECIFIC_TAG) --install --wait --atomic

manually-deploy-nodecleaner:
	helm upgrade murmurations-nodecleaner ./charts/murmurations/charts/nodecleaner --set global.env=$(ENV),image=murmurations/nodecleaner:$(SPECIFIC_TAG) --install --wait --atomic

manually-deploy-schemaparser:
	helm upgrade murmurations-schemaparser ./charts/murmurations/charts/schemaparser --set global.env=$(ENV),image=murmurations/schemaparser:$(SPECIFIC_TAG) --install --wait --atomic

manually-deploy-revalidatenode:
	helm upgrade murmurations-revalidatenode ./charts/murmurations/charts/revalidatenode --set global.env=$(ENV),image=murmurations/revalidatenode:$(SPECIFIC_TAG) --install --wait --atomic

manually-deploy-geoip:
	helm upgrade murmurations-geoip ./charts/murmurations/charts/geoip --set global.env=$(ENV),image=murmurations/geoip:$(SPECIFIC_TAG) --install --wait --atomic
