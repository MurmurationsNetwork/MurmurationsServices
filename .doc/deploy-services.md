# Deploy Services

## 1. Switch Contexts

Switch to the context you want to deploy.

```
kubectl config get-contexts

kubectl config use-context CONTEXT_NAME
```

## 2. Modify Makefile

Open `Makefile` and change `SPECIFIC_TAG` to the tag you want to deploy.

*By default we use the latest but you might want to deploy a specific version.*

## 3. Deploy Services

```

make helm-production-mq
make helm-production-ingress
make helm-production-index
make helm-production-validation
make helm-production-library
make helm-production-nodecleaner
make helm-production-schemaparser
make helm-production-revalidatenode
...
```
