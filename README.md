# Murmurations Services

_This project is licensed under the terms of the GNU General Public License v3.0_

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/MurmurationsNetwork/MurmurationsServices/CI?style=flat-square)](https://github.com/MurmurationsNetwork/MurmurationsServices/actions?query=workflow:CI)
[![Go Report Card](https://goreportcard.com/badge/github.com/MurmurationsNetwork/MurmurationsServices)](https://goreportcard.com/report/github.com/MurmurationsNetwork/MurmurationsServices)

## Run locally

*Setup*

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop) and enable Kubernetes

2. Install [Helm](https://helm.sh/docs/intro/install/)

3. Install [NGINX Ingress Controller](.doc/ingress-nginx)

4. Install [Skaffold](https://skaffold.dev/docs/install/)

5. Download large docker files

```
docker pull elasticsearch:7.16.2
docker pull kibana:7.16.2
```

7. [Create secrets](.doc/secrets.md) for each service

8. Add `127.0.0.1 index.murmurations.dev` & `127.0.0.1 library.murmurations.dev` to your host file

*After completing the setup*

1. Run `make dev` to start services

2. Try `index.murmurations.dev/v2/ping` or `library.murmurations.dev/v1/ping`

## Run in DigitalOcean

1. Install [Helm](https://helm.sh/docs/intro/install/) and [doctl](https://github.com/digitalocean/doctl#installing-doctl)

2. Create Kubernetes Clusters in DigitalOcean and install [metrics-server](https://github.com/kubernetes-sigs/metrics-server#installation) for monitoring CPU/MEM

3. Install [NGINX Ingress Controller](.doc/ingress-nginx)

4. [Create secrets](.doc/secrets.md) for each service

5. [Create PVCs](.doc/pvcs.md) for each service

6. [Deploy services](.doc/deploy-services.md)

7. [Installing and Configuring Cert-Manager](.doc/cert-manager.md)

8. Try `index.murmurations.network/v2/ping` or `library.murmurations.network/v1/ping`

**Optional**

- [Logging, Monitoring and Alerting](.doc/logging-monitoring-alerting/README.md)
