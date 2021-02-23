# Logging in Kubernetes with EFK Stack

## Install Elastic Stack (EFK) Elastic, FluentD, Kibana

**install elastic search chart**

```
helm repo add elastic https://Helm.elastic.co
helm repo update
helm install elasticsearch elastic/elasticsearch -f values-elasticsearch.yaml -n kube-logging
```

**install Kibana chart**

```
helm install kibana elastic/kibana -f values-kibana.yaml -n kube-logging
```

**install Fluentd**

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install fluentd bitnami/fluentd -n kube-logging
```

**apply Fluentd config**

```
kubectl apply -f fluentd-config.yaml
```

**restart Fluentd deamonSet**

```
kubectl rollout restart daemonset/fluentd -n kube-logging
```

## Other useful commands

**access Kibana locally**

```
kubectl port-forward deployment/kibana-kibana 5601 -n kube-logging
access: localhost:5601
```

**restart elastic search statefulSet**

```
kubectl rollout restart statefulset/elasticsearch-master -n kube-logging
```
