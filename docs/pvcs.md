# How to Creat Persistent Volume Claim

_PVCs should be created only once and never remove them unless it's needed._

## Index MongoDB

```
kubectl apply -f docs/pvc/index.yaml
kubectl apply -f docs/pvc/library.yaml
kubectl apply -f docs/pvc/dataproxy.yaml
```
