# Set Up Cert-Manager on DigitalOcean Kubernetes

## Create A records

You’ll need to ensure that your domains are pointed to the Load Balancer via `A` records. This is done through your DNS provider.

```
Create A records for index.murmurations.network and library.murmurations.network, both pointing to 139.59.201.176
```

## Install cert-manager with Helm

```
kubectl create namespace cert-manager

helm repo add jetstack https://charts.jetstack.io

helm repo update

helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v1.2.0 \
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
