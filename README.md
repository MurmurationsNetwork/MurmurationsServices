<div align="center">
<br/>
<h1>Murmurations Services</h1>
</div>

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/MurmurationsNetwork/MurmurationsServices/CI?style=flat-square)](https://github.com/MurmurationsNetwork/MurmurationsServices/actions?query=workflow:CI)
[![Go Report Card](https://goreportcard.com/badge/github.com/MurmurationsNetwork/MurmurationsServices)](https://goreportcard.com/report/github.com/MurmurationsNetwork/MurmurationsServices)

## Run locally

*Setup*

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop) and enable Kubernetes

2. Install [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/)

3. Install [Skaffold](https://skaffold.dev/docs/install/)

4. Install [Helm](https://helm.sh/docs/intro/install/)

5. Add `127.0.0.1 index.murmurations.dev` & `127.0.0.1 library.murmurations.dev` to your host file

*After completing the setup*

1. Create secrets for each service

2. Run `make dev` to start services

3. Try `index.murmurations.dev/ping` or `library.murmurations.dev/ping`

## Create Secrets

```bash
kubectl \
  create secret generic index-mongo-secret \
  --from-literal="MONGO_INITDB_ROOT_USERNAME=index-admin" \
  --from-literal="MONGO_INITDB_ROOT_PASSWORD=password"

kubectl \
  create secret generic library-mongo-secret \
  --from-literal="MONGO_INITDB_ROOT_USERNAME=library-admin" \
  --from-literal="MONGO_INITDB_ROOT_PASSWORD=password"

kubectl \
  create secret generic index-secret \
  --from-literal="MONGO_USERNAME=index-admin" \
  --from-literal="MONGO_PASSWORD=password"

kubectl \
  create secret generic library-secret \
  --from-literal="MONGO_USERNAME=library-admin" \
  --from-literal="MONGO_PASSWORD=password"

kubectl \
  create secret generic schemaparser-secret \
  --from-literal="MONGO_USERNAME=library-admin" \
  --from-literal="MONGO_PASSWORD=password"

kubectl \
  create secret generic nodecleaner-secret \
  --from-literal="MONGO_USERNAME=index-admin" \
  --from-literal="MONGO_PASSWORD=password"
```

## Directory Layout

* **common** contains packages that are shared across different services
* **services** contains packages that compile to applications that are long-running processes (such as API servers)
