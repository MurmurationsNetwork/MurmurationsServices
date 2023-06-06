#--------------------------
# Include other Makefiles.
#--------------------------
include ./build/geoip/mk/Makefile
include ./build/index/mk/Makefile
include ./build/library/mk/Makefile
include ./build/revalidatenode/mk/Makefile
include ./build/nodecleaner/mk/Makefile
include ./build/schemaparser/mk/Makefile
include ./build/validation/mk/Makefile

#--------------------------
# Set environment variables.
#--------------------------
DEPLOY_ENV ?= local

ifeq ($(DEPLOY_ENV), staging)
	ENV_FILE = test/e2e-staging-env.json
else
	ENV_FILE = test/e2e-local-env.json
endif

#--------------------------
# Set up the development servers
#--------------------------
.PHONY: dev
dev:
	export SOURCEPATH=$(PWD) && skaffold dev --port-forward

#--------------------------
# Runs the unit tests.
#--------------------------
.PHONY: test
test:
	export ENV=test && go test ./...

#--------------------------
# Run the end-to-end (E2E) tests using newman.
#--------------------------
.PHONY: newman-test
newman-test:
	newman run test/e2e-tests.json -e $(ENV_FILE) --verbose --delay-request 1000

# ---------------------------------------------------------------

docker-build-dataproxy:
	$(MAKE) -C services/dataproxy/ docker-build

docker-build-dataproxyupdater:
	$(MAKE) -C services/cronjob/dataproxyupdater/ docker-build

docker-build-dataproxyrefresher:
	$(MAKE) -C services/cronjob/dataproxyrefresher/ docker-build

# ---------------------------------------------------------------

# The TAG value is constructed based on the commit SHA.
# If running in a GitHub Actions environment, it uses the GITHUB_SHA.
# In a local environment, it uses the SHA of the HEAD commit.
TAG ?= $(shell git rev-parse --short $(if $(GITHUB_SHA),$(GITHUB_SHA),HEAD))

check-clean:
	@if [ -n "$(shell git status --porcelain)" ]; then \
		echo "Uncommitted changes present. Please commit them before running this command."; \
		exit 1; \
	fi

docker-tag-index: check-clean docker-build-index
	docker tag murmurations/index murmurations/index:$(TAG)

docker-tag-validation: check-clean docker-build-validation
	docker tag murmurations/validation murmurations/validation:${TAG}

docker-tag-library: check-clean docker-build-library
	docker tag murmurations/library murmurations/library:${TAG}

docker-tag-nodecleaner: check-clean docker-build-nodecleaner
	docker tag murmurations/nodecleaner murmurations/nodecleaner:${TAG}

docker-tag-schemaparser: check-clean docker-build-schemaparser
	docker tag murmurations/schemaparser murmurations/schemaparser:${TAG}

docker-tag-revalidatenode: check-clean docker-build-revalidatenode
	docker tag murmurations/revalidatenode murmurations/revalidatenode:${TAG}

docker-tag-geoip: check-clean docker-build-geoip
	docker tag murmurations/geoip murmurations/geoip:${TAG}

docker-tag-dataproxy: check-clean docker-build-dataproxy
	docker tag murmurations/dataproxy murmurations/dataproxy:${TAG}

docker-tag-dataproxyupdater: check-clean docker-build-dataproxyupdater
	docker tag murmurations/dataproxyupdater murmurations/dataproxyupdater:${TAG}

docker-tag-dataproxyrefresher: check-clean docker-build-dataproxyrefresher
	docker tag murmurations/dataproxyrefresher murmurations/dataproxyrefresher:${TAG}

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

docker-push-dataproxy: docker-tag-dataproxy
	docker push murmurations/dataproxy:latest
	docker push murmurations/dataproxy:$(TAG)

docker-push-dataproxyupdater: docker-tag-dataproxyupdater
	docker push murmurations/dataproxyupdater:latest
	docker push murmurations/dataproxyupdater:$(TAG)

docker-push-dataproxyrefresher: docker-tag-dataproxyrefresher
	docker push murmurations/dataproxyrefresher:latest
	docker push murmurations/dataproxyrefresher:$(TAG)

# ---------------------------------------------------------------

deploy-ingress:
	helm upgrade murmurations-ingress ./charts/murmurations/charts/ingress --set global.env=$(DEPLOY_ENV) --install --atomic

deploy-mq:
	helm upgrade murmurations-mq ./charts/murmurations/charts/message-queue --set global.env=$(DEPLOY_ENV) --install --atomic

deploy-index:
	helm upgrade murmurations-index ./charts/murmurations/charts/index --set global.env=$(DEPLOY_ENV),image=murmurations/index:$(TAG) --install --atomic

deploy-validation:
	helm upgrade murmurations-validation ./charts/murmurations/charts/validation --set global.env=$(DEPLOY_ENV),image=murmurations/validation:$(TAG) --install --atomic

deploy-library:
	helm upgrade murmurations-library ./charts/murmurations/charts/library --set global.env=$(DEPLOY_ENV),image=murmurations/library:$(TAG) --install --atomic

deploy-nodecleaner:
	helm upgrade murmurations-nodecleaner ./charts/murmurations/charts/nodecleaner --set global.env=$(DEPLOY_ENV),image=murmurations/nodecleaner:$(TAG) --install --atomic

deploy-schemaparser:
	helm upgrade murmurations-schemaparser ./charts/murmurations/charts/schemaparser --set global.env=$(DEPLOY_ENV),image=murmurations/schemaparser:$(TAG) --install --atomic

deploy-revalidatenode:
	helm upgrade murmurations-revalidatenode ./charts/murmurations/charts/revalidatenode --set global.env=$(DEPLOY_ENV),image=murmurations/revalidatenode:$(TAG) --install --atomic

deploy-geoip:
	helm upgrade murmurations-geoip ./charts/murmurations/charts/geoip --set global.env=$(DEPLOY_ENV),image=murmurations/geoip:$(TAG) --install --atomic

deploy-dataproxy:
	helm upgrade murmurations-dataproxy ./charts/murmurations/charts/dataproxy --set global.env=$(DEPLOY_ENV),image=murmurations/dataproxy:$(TAG) --install --atomic

deploy-dataproxyupdater:
	helm upgrade murmurations-dataproxyupdater ./charts/murmurations/charts/dataproxyupdater --set global.env=$(DEPLOY_ENV),image=murmurations/dataproxyupdater:$(TAG) --install --atomic

deploy-dataproxyrefresher:
	helm upgrade murmurations-dataproxyrefresher ./charts/murmurations/charts/dataproxyrefresher --set global.env=$(DEPLOY_ENV),image=murmurations/dataproxyrefresher:$(TAG) --install --atomic

# ---------------------------------------------------------------

# Please update to the version you want to deploy.
SPECIFIC_TAG ?= <>
ENV ?= <>

manually-deploy-ingress:
	helm upgrade murmurations-ingress ./charts/murmurations/charts/ingress --set global.env=$(ENV) --install --atomic

manually-deploy-mq:
	helm upgrade murmurations-mq ./charts/murmurations/charts/message-queue --set global.env=$(ENV) --install --atomic

manually-deploy-index:
	helm upgrade murmurations-index ./charts/murmurations/charts/index --set global.env=$(ENV),image=murmurations/index:$(SPECIFIC_TAG) --install --atomic

manually-deploy-validation:
	helm upgrade murmurations-validation ./charts/murmurations/charts/validation --set global.env=$(ENV),image=murmurations/validation:$(SPECIFIC_TAG) --install --atomic

manually-deploy-library:
	helm upgrade murmurations-library ./charts/murmurations/charts/library --set global.env=$(ENV),image=murmurations/library:$(SPECIFIC_TAG) --install --atomic

manually-deploy-nodecleaner:
	helm upgrade murmurations-nodecleaner ./charts/murmurations/charts/nodecleaner --set global.env=$(ENV),image=murmurations/nodecleaner:$(SPECIFIC_TAG) --install --atomic

manually-deploy-schemaparser:
	helm upgrade murmurations-schemaparser ./charts/murmurations/charts/schemaparser --set global.env=$(ENV),image=murmurations/schemaparser:$(SPECIFIC_TAG) --install --atomic

manually-deploy-revalidatenode:
	helm upgrade murmurations-revalidatenode ./charts/murmurations/charts/revalidatenode --set global.env=$(ENV),image=murmurations/revalidatenode:$(SPECIFIC_TAG) --install --atomic

manually-deploy-geoip:
	helm upgrade murmurations-geoip ./charts/murmurations/charts/geoip --set global.env=$(ENV),image=murmurations/geoip:$(SPECIFIC_TAG) --install --atomic

manually-deploy-dataproxy:
	helm upgrade murmurations-dataproxy ./charts/murmurations/charts/dataproxy --set global.env=$(ENV),image=murmurations/dataproxy:$(SPECIFIC_TAG) --install --atomic

manually-deploy-dataproxyupdater:
	helm upgrade murmurations-dataproxyupdater ./charts/murmurations/charts/dataproxyupdater --set global.env=$(ENV),image=murmurations/dataproxyupdater:$(SPECIFIC_TAG) --install --atomic

manually-deploy-dataproxyrefresher:
	helm upgrade murmurations-dataproxyrefresher ./charts/murmurations/charts/dataproxyrefresher --set global.env=$(ENV),image=murmurations/dataproxyrefresher:$(SPECIFIC_TAG) --install --atomic
