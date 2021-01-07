# Remove Completed Kubernetes Jobs

```
kubectl delete jobs --field-selector status.successful=1
kubectl delete jobs --field-selector status.successful=0
```
