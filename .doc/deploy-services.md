# Deploy Services

## 1. Switch Contexts

Switch to the context you want to deploy.

```
kubectl config get-contexts

kubectl config use-context CONTEXT_NAME
```

## 2. Modify Makefile

1. Open `Makefile`
2. Change `SPECIFIC_TAG` to the tag you want to deploy.
3. Change `ENV` to the environment you want to deploy.

_By default we use the latest tag but you might want to deploy a specific version._

## 3. Deploy Services

```
make manually-deploy-mq
make manually-deploy-index
make manually-deploy-validation
make manually-deploy-library
make manually-deploy-nodecleaner
make manually-deploy-schemaparser
make manually-deploy-revalidatenode
```
