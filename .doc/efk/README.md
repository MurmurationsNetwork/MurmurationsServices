# Install Elastic Stack (EFK) Elastic, FluentD, Kibana

**Create a namespace**

```
kubectl create namespace logging
```

**Install elastic search chart**

```
helm repo add elastic https://helm.elastic.co
helm install elasticsearch elastic/elasticsearch -f values-elastic.yaml -n logging
```

**Check out the status**

```
kubectl get all -n logging
```

**Install Kibana chart**

```
helm install kibana elastic/kibana -f values-kibana.yaml -n logging
```

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
