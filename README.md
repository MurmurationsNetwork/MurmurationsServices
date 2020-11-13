<div align="center">
<br/>
<h1>Murmurations Services</h1>
</div>

## Run locally

*Prerequisite*

1. Installed [Docker Desktop](https://www.docker.com/products/docker-desktop) and enable Kubernetes.

2. Istalled [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/).

3. Installed [Skaffold](https://skaffold.dev/docs/install/).

4. Installed [Helm](https://helm.sh/docs/intro/install/).

7. Add `127.0.0.1 index.murmurations.network` to your host file.

*After finishing the prerequisite*

1. Run `make dev` to start services.

2. Try `index.murmurations.network/ping`

## Directory Layout

* **common** contains packages that shared across different services.
* **services** contains packages that compile to applications that are long-running processes (such as API servers).
