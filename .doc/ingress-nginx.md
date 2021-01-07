# How to Install Nginx Controller

## Installation

Follow [Official Installation Guide](https://kubernetes.github.io/ingress-nginx/deploy/) for more information.

**Docker for Mac**

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.43.0/deploy/static/provider/cloud/deploy.yaml
```

**Check Resources**

```
k get all -n ingress-nginx
```

## Uninstallation

```
kubectl delete all --all -n ingress-nginx
```
