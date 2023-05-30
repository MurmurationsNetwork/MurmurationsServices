# Set Up Cert-Manager on DigitalOcean Kubernetes

## Create A records

You’ll need to ensure that your domains are pointed to the Load Balancer via `A` records. This is done through your DNS provider.

```
Create A records for:
index.murmurations.network
library.murmurations.network
data-proxy.murmurations.network
monitoring.murmurations.network

All pointing to the IP address assigned to the load balancer of the K8s cluster.
```

## Install cert-manager with Helm

```
kubectl create namespace cert-manager

helm repo add jetstack https://charts.jetstack.io

helm repo update

helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.8.0 \
  --set installCRDs=true
```

## Install Issuer and Ingress

```
make manually-deploy-ingress
```

## Debugging

```
k get order
k describe order murmurations-network-tls-z2ndq-508344205
k describe challenge murmurations-network-tls-z2ndq-508344205-59023570
```
