# How to Creat Persistent Volume Claim

*PVCs should be created only once and never remove them unless it's needed.*

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: <index-mongo-pvc>
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  storageClassName: do-block-storage
```
