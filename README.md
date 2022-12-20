# Murmurations Services

_This project is licensed under the terms of the GNU General Public License v3.0_

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/MurmurationsNetwork/MurmurationsServices/main.yaml?branch=main&style=flat-square)](https://github.com/MurmurationsNetwork/MurmurationsServices/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/MurmurationsNetwork/MurmurationsServices?style=flat-square)](https://goreportcard.com/report/github.com/MurmurationsNetwork/MurmurationsServices)

## Run locally

*Setup*

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop) and enable Kubernetes

2. Install [Helm](https://helm.sh/docs/intro/install/)

3. Install [NGINX Ingress Controller](.doc/ingress-nginx)

4. Install [Skaffold](https://skaffold.dev/docs/install/)

5. Download large docker files

```
docker pull elasticsearch:7.17.5
docker pull kibana:7.17.5
```

6. [Create secrets](.doc/secrets.md) for each service

7. Add the following to your host file `sudo vim /etc/hosts`

```
127.0.0.1   index.murmurations.dev
127.0.0.1   library.murmurations.dev
127.0.0.1   data-proxy.murmurations.dev
```

*After completing the setup*

1. Run `make dev` to start services

2. Try `index.murmurations.dev/v2/ping`, `library.murmurations.dev/v1/ping` and `data-proxy.murmurations.dev/v1/ping`

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

- [Seed the data](.doc/seed.md)
