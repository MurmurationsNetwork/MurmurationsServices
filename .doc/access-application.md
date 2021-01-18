# How to Access Application in Kubernetes

You can use [Port Forwarding](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) and then just connect your client to it.

If you want to stop, just hit `Ctrl+C` on the same cmd window to stop the process.

```
kubectl port-forward svc/index-mongo 27017:27017
kubectl port-forward svc/library-mongo 27018:27017
kubectl port-forward svc/index-kibana 5601:5601

kubectl port-forward svc/elasticsearch 9200:9200 -n kube-logging
kubectl port-forward svc/kibana 5601:5601 -n kube-logging
```
