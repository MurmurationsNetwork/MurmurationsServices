# How to Install Nginx Controller

## Installation

Follow [Official Installation Guide](https://kubernetes.github.io/ingress-nginx/deploy/) for more information.

```
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
kubectl create namespace ingress-nginx
helm upgrade ingress-nginx ingress-nginx/ingress-nginx -f .doc/ingress-nginx/values-ingress-nginx.yaml -n ingress-nginx --install
```
