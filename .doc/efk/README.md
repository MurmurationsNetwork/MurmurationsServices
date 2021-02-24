# Logging in Kubernetes with EFK Stack

## Install Elastic Stack (EFK) Elastic, FluentD, Kibana

**install elastic search chart**

```
helm repo add elastic https://Helm.elastic.co
helm repo update
helm upgrade elasticsearch elastic/elasticsearch --version="7.9.0" -f values-elasticsearch.yaml -n kube-logging --install
```

**install Kibana chart**

```
helm upgrade kibana elastic/kibana --version="7.9.0" -f values-kibana.yaml -n kube-logging --install
```

**install Fluentd**

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm upgrade fluentd bitnami/fluentd --version="2.0.1" -n kube-logging --install
```

**apply Fluentd config**

```
kubectl apply -f fluentd-config.yaml
kubectl rollout restart daemonset/fluentd -n kube-logging
```

## Other useful commands

**access Kibana locally**

```
kubectl port-forward deployment/kibana-kibana 5601 -n kube-logging
access: localhost:5601
```

**restart Fluentd deamonSet**

```
kubectl rollout restart daemonset/fluentd -n kube-logging
```

**restart elastic search statefulSet**

```
kubectl rollout restart statefulset/elasticsearch-master -n kube-logging
```
