# How to Install Nginx Controller

## Installation

Follow [Official Installation Guide](https://kubernetes.github.io/ingress-nginx/deploy/) for more information.

**Docker for Mac**

```
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install ingress-nginx ingress-nginx/ingress-nginx -n ingress-nginx
```
