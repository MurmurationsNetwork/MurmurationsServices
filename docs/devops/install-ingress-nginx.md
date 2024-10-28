# How to Install Nginx Controller

## Production

Read the official [Installation Guide](https://kubernetes.github.io/ingress-nginx/deploy/) to install in a production environment.

## Development

```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
```

Note: k8s version v1.24.0 is not compatible with the current latest ingress-nginx v4.1.4. If the above command is not working, you can add `--version 4.0.19` after the command.
