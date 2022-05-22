# How to Creat Persistent Volume Claim

_PVCs should be created only once and never remove them unless it's needed._

## Index MongoDB

```
kubectl apply -f .doc/pvc/index.yaml
kubectl apply -f .doc/pvc/library.yaml
kubectl apply -f .doc/pvc/profile.yaml
```
