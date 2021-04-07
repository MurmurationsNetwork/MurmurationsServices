# Create Volume Snapshots and Restore Volumes from Snapshots

## Create a Snapshot of a Volume

**List all avaliable persistent volume claim**

```
kubectl get pvc
```

**Create a snapshot**

```
kubectl create -f .doc/volume-snapshots/take-volume-snapshot.yaml
```

**Observe the state of your volumes and snapshots**

```
kubectl get pvc && kubectl get pv && kubectl get volumesnapshot
```

## Restore A Volume from a Snapshot

```
kubectl create -f .doc/volume-snapshots/pvc-from-snapshot.yaml
```

**Observe the state of your volumes and snapshots**

```
kubectl get pvc && kubectl get pv && kubectl get volumesnapshot
```

## Ask DB to Pick Up The Volume

**Change `claimName` in deployment yaml file for the database**

```yaml
# ...
volumes:
- name: mongo-storage
    persistentVolumeClaim:
    claimName: index-mongo-pvc-restore
# ...
```

**Deploy manually**

Make sure you specify the correct image tag and environment.

```bash
make manually-deploy-index
```
