# Set Up Cert-Manager on DigitalOcean Kubernetes

Please follow this article for more information [How to Set Up an Nginx Ingress with Cert-Manager on DigitalOcean Kubernetes](https://www.digitalocean.com/community/tutorials/how-to-set-up-an-nginx-ingress-with-cert-manager-on-digitalocean-kubernetes)

## Install cert-manager and its Custom Resource Definitions (CRDs)

```
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.2.0/cert-manager.yaml
```

## Verify Installation

```
kubectl get pods --namespace cert-manager
```

## Create A records

You’ll need to ensure that your domains are pointed to the Load Balancer via `A` records. This is done through your DNS provider.

```
Create A records for index.murmurations.network and library.murmurations.network, both pointing to 139.59.201.176
```

## Create one issues Let’s Encrypt certificates

```
make helm-production-ingress
```
