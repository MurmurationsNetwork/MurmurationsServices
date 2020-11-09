<div align="center">
<br/>
<h1>Murmurations Services</h1>
</div>

## Run locally


1. Make sure you have installed [Docker Desktop](https://www.docker.com/products/docker-desktop) and enable Kubernetes.

2. Make sure you have installed [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/).

3. Make sure you have installed [Skaffold](https://skaffold.dev/docs/install/).

4. Make sure you have installed [Helm](https://helm.sh/docs/intro/install/).

5. Run `make dev` to start services.

## Directory Layout

* **charts** contains Helm configuration files.
* **common** contains packages that shared across different services.
* **services** contains packages that compile to applications that are long-running processes (such as API servers).
