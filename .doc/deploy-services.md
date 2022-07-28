# Service Deployment

By default, both the staging (`test-index`) and production (`index`) environments will deploy automatically when a commit is made to the `test` or `main` branches respectively.

New PRs will be deployed to the pretest (`pretest-index`) environment.

## Deploy Services Manually

### 1. Switch Contexts

Switch to the context you want to deploy.

```
kubectl config get-contexts

kubectl config use-context CONTEXT_NAME
```

### 2. Modify Makefile

1. Open `Makefile`
2. Change `SPECIFIC_TAG` to the tag you want to deploy. (check out [dockerhub](https://hub.docker.com/r/murmurations/index/tags) or [github action](https://github.com/MurmurationsNetwork/MurmurationsServices/runs/3836865026?check_suite_focus=true#step:4:191) to find the tag)
3. Change `ENV` to the environment you want to deploy. (`production`)

_By default we use the latest tag but you might want to deploy a specific version._

### 3. Deploy Services

```
make manually-deploy-mq
make manually-deploy-index
make manually-deploy-validation
make manually-deploy-library
make manually-deploy-nodecleaner
make manually-deploy-schemaparser
make manually-deploy-revalidatenode
```
