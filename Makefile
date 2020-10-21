dev:
	skaffold dev --port-forward

index-dev:
	$(MAKE) -C services/index/ docker-build-dev

validation-dev:
	$(MAKE) -C services/validation/ docker-build-dev
