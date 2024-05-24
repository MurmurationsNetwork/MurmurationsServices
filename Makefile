#--------------------------
# Include other Makefiles.
#--------------------------
include ./build/dataproxy/mk/Makefile
include ./build/dataproxyrefresher/mk/Makefile
include ./build/dataproxyupdater/mk/Makefile
include ./build/index/mk/Makefile
include ./build/library/mk/Makefile
include ./build/nodecleaner/mk/Makefile
include ./build/revalidatenode/mk/Makefile
include ./build/schemaparser/mk/Makefile
include ./build/validation/mk/Makefile
include ./build/maintenance/mk/Makefile

#--------------------------
# Set environment variables.
#--------------------------
DEPLOY_ENV ?= development

ifeq ($(DEPLOY_ENV), staging)
	ENV_FILE = test/e2e-env-staging.json
else ifeq ($(DEPLOY_ENV), pretest)
	ENV_FILE = test/e2e-env-pretest.json
else
	ENV_FILE = test/e2e-env-development.json
endif

VALUES_FILE=./charts/murm-queue/values-contabo.yaml

# The TAG value is constructed based on the commit SHA.
# If running in a GitHub Actions environment, it uses the GITHUB_SHA.
# In a local environment, it uses the SHA of the HEAD commit.
TAG ?= $(shell git rev-parse --short $(if $(GITHUB_SHA),$(GITHUB_SHA),HEAD))

DOCKER_TAG_PREFIX := $(if $(filter production,$(DEPLOY_ENV)),,${DEPLOY_ENV}-)

#--------------------------
# Set up the development servers
#--------------------------
.PHONY: dev
dev:
    # There's no particular deployment order for Helm, so we use --tolerate-failures-until-deadline to
    # prevent deployment failure if the required object, such as PriorityClass, has not been deployed yet.
	export SOURCEPATH=$(PWD) && skaffold dev --tolerate-failures-until-deadline=true --port-forward

#--------------------------
# Runs the unit tests.
#--------------------------
.PHONY: test
test:
	export APP_ENV=test && go test ./...

#--------------------------
# Run the end-to-end (E2E) tests using newman.
#--------------------------
.PHONY: newman-test
newman-test:
	newman run test/e2e-tests-staging.json -e $(ENV_FILE) --verbose --delay-request 10

check-clean:
	@if [ -n "$(shell git status --porcelain)" ]; then \
		echo "Uncommitted changes present. Please commit them before running this command."; \
		exit 1; \
	fi

# ---------------------------------------------------------------

deploy-murmurations-core:
	helm upgrade murmurations-core ./charts/murmurations/charts/core \
	--set global.env=$(DEPLOY_ENV) --install --atomic

deploy-ingress:
	helm upgrade murmurations-ingress ./charts/murmurations/charts/ingress \
	--set global.env=$(DEPLOY_ENV) --install --atomic

deploy-nats:
	helm repo add nats https://nats-io.github.io/k8s/helm/charts/ && \
	helm repo update && \
	helm upgrade nats nats/nats \
	--namespace murm-queue \
	--create-namespace \
	--install \
	--atomic \
	--set global.env=$(DEPLOY_ENV) \
	--version 1.1.6 \
	-f $(VALUES_FILE)

deploy-index:
	helm upgrade murmurations-index ./charts/murmurations/charts/index \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)index:$(TAG) \
	--install --atomic

deploy-validation:
	helm upgrade murmurations-validation ./charts/murmurations/charts/validation \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)validation:$(TAG) \
	--install --atomic

deploy-library:
	helm upgrade murmurations-library ./charts/murmurations/charts/library \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)library:$(TAG) \
	--install --atomic

deploy-nodecleaner:
	helm upgrade murmurations-nodecleaner ./charts/murmurations/charts/nodecleaner \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)nodecleaner:$(TAG) \
	--install --atomic

deploy-schemaparser:
	helm upgrade murmurations-schemaparser ./charts/murmurations/charts/schemaparser \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)schemaparser:$(TAG) \
	--install --atomic

deploy-revalidatenode:
	helm upgrade murmurations-revalidatenode ./charts/murmurations/charts/revalidatenode \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)revalidatenode:$(TAG) \
	--install --atomic

deploy-dataproxy:
	helm upgrade murmurations-dataproxy ./charts/murmurations/charts/dataproxy \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)dataproxy:$(TAG) \
	--install --atomic

deploy-dataproxyupdater:
	helm upgrade murmurations-dataproxyupdater ./charts/murmurations/charts/dataproxyupdater \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)dataproxyupdater:$(TAG) \
	--install --atomic

deploy-dataproxyrefresher:
	helm upgrade murmurations-dataproxyrefresher ./charts/murmurations/charts/dataproxyrefresher \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)dataproxyrefresher:$(TAG) \
	--install --atomic

# ---------------------------------------------------------------
# Manual Helm deployment targets for individual services with debugging.
# ---------------------------------------------------------------

# Set the specific tag and environment for manual deployments
SPECIFIC_TAG ?= latest
MANUAL_DEPLOY_TARGETS = manually-deploy-murmurations-core \
                        manually-deploy-ingress \
                        manually-deploy-nats \
                        manually-deploy-index \
                        manually-deploy-validation \
                        manually-deploy-library \
                        manually-deploy-nodecleaner \
                        manually-deploy-schemaparser \
                        manually-deploy-revalidatenode \
                        manually-deploy-dataproxy \
                        manually-deploy-dataproxyupdater \
                        manually-deploy-dataproxyrefresher \

deploy-all-services: $(MANUAL_DEPLOY_TARGETS)

manually-deploy-murmurations-core:
	helm upgrade murmurations-core ./charts/murmurations/charts/core \
		--set global.env=$(DEPLOY_ENV) --install --atomic --debug

manually-deploy-ingress:
	helm upgrade murmurations-ingress ./charts/murmurations/charts/ingress \
	--set global.env=$(DEPLOY_ENV) --install --atomic --debug

manually-deploy-nats:
	helm repo add nats https://nats-io.github.io/k8s/helm/charts/ && \
	helm repo update && \
	helm upgrade nats nats/nats \
	--namespace murm-queue \
	--create-namespace \
	--install \
	--atomic \
	--set global.env=$(DEPLOY_ENV) \
	--version 1.1.6 \
	-f $(VALUES_FILE)

manually-deploy-index:
	helm upgrade murmurations-index ./charts/murmurations/charts/index \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)index:$(SPECIFIC_TAG) \
	--install --atomic --debug

manually-deploy-validation:
	helm upgrade murmurations-validation \
	./charts/murmurations/charts/validation \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)validation:$(SPECIFIC_TAG) \
	--install --atomic --debug

manually-deploy-library:
	helm upgrade murmurations-library ./charts/murmurations/charts/library \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)library:$(SPECIFIC_TAG) \
	--install --atomic --debug

manually-deploy-nodecleaner:
	helm upgrade murmurations-nodecleaner \
	./charts/murmurations/charts/nodecleaner \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)nodecleaner:$(SPECIFIC_TAG) \
	--install --atomic --debug

manually-deploy-schemaparser:
	helm upgrade murmurations-schemaparser \
	./charts/murmurations/charts/schemaparser \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)schemaparser:$(SPECIFIC_TAG) \
	--install --atomic --debug

manually-deploy-revalidatenode:
	helm upgrade murmurations-revalidatenode \
	./charts/murmurations/charts/revalidatenode \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)revalidatenode:$(SPECIFIC_TAG) \
	--install --atomic --debug

manually-deploy-dataproxy:
	helm upgrade murmurations-dataproxy ./charts/murmurations/charts/dataproxy \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)dataproxy:$(SPECIFIC_TAG) \
	--install --atomic --debug

manually-deploy-dataproxyupdater:
	helm upgrade murmurations-dataproxyupdater \
	./charts/murmurations/charts/dataproxyupdater \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)dataproxyupdater:$(SPECIFIC_TAG) \
	--install --atomic --debug

manually-deploy-dataproxyrefresher:
	helm upgrade murmurations-dataproxyrefresher \
	./charts/murmurations/charts/dataproxyrefresher \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)dataproxyrefresher:$(SPECIFIC_TAG) \
	--install --atomic --debug

manually-deploy-murm-logging:
	helm upgrade murm-logging ./charts/murm-logging \
	--namespace murm-logging \
	--create-namespace \
	--set global.env=$(DEPLOY_ENV) \
	--set elasticsearch.image=docker.elastic.co/elasticsearch/elasticsearch:8.12.1 \
	--set kibana.image=docker.elastic.co/kibana/kibana:8.12.1 \
	--install --atomic --debug

manually-deploy-maintenance:
	helm upgrade murmurations-maintenance ./charts/murmurations/charts/maintenance \
	--set global.env=$(DEPLOY_ENV),image=murmurations/$(DOCKER_TAG_PREFIX)maintenance:$(SPECIFIC_TAG) \
	--install --atomic --debug
