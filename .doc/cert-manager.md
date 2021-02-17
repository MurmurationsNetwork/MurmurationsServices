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

## Create one issues Let’s Encrypt certificates

```
make helm-production-core
```
