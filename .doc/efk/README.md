# Logging in Kubernetes with EFK Stack

## Install Elastic Stack (EFK) Elastic, FluentD, Kibana

**install elastic search chart**

```
helm repo add elastic https://Helm.elastic.co
helm repo update
helm upgrade elasticsearch elastic/elasticsearch --version="7.9.0" -f ./.doc/efk/values-elasticsearch.yaml -n kube-logging --install
```

**install Kibana chart**

```
helm upgrade kibana elastic/kibana --version="7.9.0" -f ./.doc/efk/values-kibana.yaml -n kube-logging --install
```

**install Fluentd**

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm upgrade fluentd bitnami/fluentd --version="2.0.1" -n kube-logging --install
```

**apply Fluentd config**

```
kubectl apply -f ./.doc/efk/fluentd-config.yaml
kubectl rollout restart daemonset/fluentd -n kube-logging
```

**create index in Kibana**

1. Login to Kibana
2. Choose Kibana > Discover from menu
3. Create index with `*app*`

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
