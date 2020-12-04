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

7. Add `127.0.0.1 index.murmurations.network.dev` & `127.0.0.1 library.murmurations.network.dev` to your host file

*After completing the setup*

1. Run `make dev` to start services

2. Try `index.murmurations.network.dev/ping` or `library.murmurations.network.dev/ping`

## Directory Layout

* **common** contains packages that are shared across different services
* **services** contains packages that compile to applications that are long-running processes (such as API servers)
