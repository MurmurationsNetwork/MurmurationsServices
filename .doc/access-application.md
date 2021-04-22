# How to Access Application in Kubernetes

You can use [Port Forwarding](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) and then just connect your client to it.

If you want to stop, just hit `Ctrl+C` on the same cmd window to stop the process.

```
kubectl port-forward svc/index-mongo 27017:27017
kubectl port-forward svc/library-mongo 27018:27017
kubectl port-forward svc/index-kibana 5601:5601
kubectl port-forward svc/schemaparser-redis 6379:6379

kubectl port-forward svc/grafana 3000:80 -n kube-monitoring
```
