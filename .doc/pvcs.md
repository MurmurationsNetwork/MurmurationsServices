# How to Creat Persistent Volume Claim

*PVCs should be created only once and never remove them unless it's needed.*

## Index MongoDB

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: index-mongo-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: do-block-storage
```

## Library MongoDB

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: library-mongo-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: do-block-storage
```
