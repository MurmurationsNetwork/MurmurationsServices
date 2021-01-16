# Install Elastic Stack (EFK) Elastic, FluentD, Kibana

**Install Fluentd**

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install fluentd bitnami/fluentd -f values-fluentd.yaml -n logging
```

## Other useful commands

**Restart Fluentd deamonSet**

```
kubectl rollout restart daemonset/fluentd -n logging
```

**Restart elastic search statefulSet**

```
kubectl rollout restart statefulset/elasticsearch-master -n logging
```
