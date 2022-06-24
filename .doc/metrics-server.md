# Kubernetes Metrics Server

```
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

kubectl top pod POD_NAME --containers               # Show metrics for a given pod and its containers
kubectl top pod POD_NAME --sort-by=cpu              # Show metrics for a given pod and sort it by 'cpu' or 'memory'
```
