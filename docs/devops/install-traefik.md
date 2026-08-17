# How to Install Traefik Controller

Traefik replaces the NGINX Ingress Controller as the ingress controller for
Murmurations Services. The chart in `charts/murmurations/charts/ingress` renders
`ingressClassName: traefik` for every environment.

No controller-level configuration is required. Everything the ingress needs
(the 8 MB request body limit and the HTTP -> HTTPS redirection) ships with the
chart as Traefik `Middleware` resources.

## Production

RKE2 v1.36 and later deploy Traefik as a bundled chart (`rke2-traefik` in the
`kube-system` namespace) and create the `traefik` IngressClass, so no separate
install is needed. If customization is ever required, it is done with a
`HelmChartConfig` named `rke2-traefik` in `kube-system` rather than by running
Helm directly.

Two cases need attention:

- **Clusters upgraded from an earlier version.** RKE2 keeps ingress-nginx as
  the default on upgrade — it does not switch automatically. Follow the
  official [Ingress NGINX to Traefik Migration Guide](https://docs.rke2.io/reference/ingress_migration).
- **Clusters older than v1.36.** Migrate before upgrading. The ingress-nginx
  chart receives no further updates and is removed entirely in v1.37.

## Development

Docker Desktop's Kubernetes does not bundle an ingress controller, so install
Traefik with Helm. The release name determines the IngressClass name, so keep
it as `traefik`:

```bash
helm repo add traefik https://traefik.github.io/charts
helm repo update
helm upgrade --install traefik traefik/traefik \
  --namespace traefik --create-namespace
```
