# Logging, Monitoring and Alerting

## Prerequisite

### 1. Create a Name Space

```
kubectl create namespace kube-monitoring
```

### 2. Add Helm Charts

```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add kube-state-metrics https://kubernetes.github.io/kube-state-metrics
helm repo add grafana https://grafana.github.io/helm-charts
helm repo add loki https://grafana.github.io/loki/charts
helm repo update
```

## Deployment (Updated at 2022/07/01)
[Ref](https://github.com/digitalocean/Kubernetes-Starter-Kit-Developers/tree/main/04-setup-prometheus-stack)

### 1. Deploy Prometheus with Grafana

- Revise `adminPassword` in .doc/logging-monitoring-alerting/prom-stack-values.yaml

```
helm upgrade kube-prom-stack prometheus-community/kube-prometheus-stack --version 40.3.1 -n kube-monitoring -f .doc/logging-monitoring-alerting/prom-stack-values.yaml --install
```

- Check the available resources
```
kubectl get all -n kube-monitoring
```

- Connect to the services
```
kubectl port-forward svc/kube-prom-stack-grafana 3000:80 -n kube-monitoring
kubectl port-forward svc/kube-prom-stack-kube-prome-prometheus 9090:9090 -n kube-monitoring
kubectl port-forward alertmanager-kube-prom-stack-kube-prome-alertmanager-0 9093 -n kube-monitoring
```

### 2. Install Loki & Promtail

```
helm upgrade loki grafana/loki-stack --version 2.8.3  -n kube-monitoring -f .doc/logging-monitoring-alerting/loki-stack-values.yaml --install
```

- Navigate to Grafana: https://localhost:3000
```
kubectl port-forward svc/kube-prom-stack-grafana 3000:80 -n kube-monitoring
```

- Configuration -> Data Source -> Add data source -> Select Loki -> set url as `http://loki:3100` -> Save & Test -> The page shows 'Data source connected and labels found.' (Successful!)

```
kubectl port-forward daemonset/promtail 3101 -n kube-monitoring
```

on http://localhost:3101/targets

You will see all the pods which are getting scraped by promtail for logs.

### 3. Setup Grafana

Navigate to http://localhost:3000

<!--**Add Alert Webhook**

Adding alert for certain things from dashboard.

![image](https://user-images.githubusercontent.com/11765228/115104231-9ab7ba00-9f89-11eb-9d06-cf4d592b1b03.png)

Add Webhook Url

![image](https://user-images.githubusercontent.com/11765228/115104198-6a701b80-9f89-11eb-8d5e-4d69e1446d03.png)

Check these boxes as well

![image](https://user-images.githubusercontent.com/11765228/115104205-79ef6480-9f89-11eb-804c-c8e1828ccca1.png)-->

### 6. Configure the Grafana Dashboard

![image](https://user-images.githubusercontent.com/11765228/115194754-780bd980-a120-11eb-9284-c01458983f6b.png)

Type 1860 for Node Exporter Full Dashboard

![image](https://user-images.githubusercontent.com/11765228/115194824-907bf400-a120-11eb-86f3-68d06aa5ffcd.png)

![image](https://user-images.githubusercontent.com/11765228/115195026-dafd7080-a120-11eb-89ac-2af4e5120ea1.png)

Now again add one more dashboard: 8685

![image](https://user-images.githubusercontent.com/11765228/115195120-f49eb800-a120-11eb-971a-993c668e6af4.png)

- Find loki's uid in setting > datasources
- Replace all of loki with new uid (`"uid": "pvuSqSenk"`) in [(.doc/logging-monitoring-alerting/grafana-logging.json)](.doc/logging-monitoring-alerting/grafana-logging.json)
- Import via panel json [(.doc/logging-monitoring-alerting/grafana-logging.json)](.doc/logging-monitoring-alerting/grafana-logging.json)

## Alerting

If you want to add an alert on a pre-built specific property.

![](https://i.imgur.com/aXYWiPy.png)

After clicking on edit, you will reach this page. 

![](https://i.imgur.com/wo4GiAM.png)

*Note: Alert option will only appear in Graph visualization*

- For 1, You may use any [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/)

Note: PromQL should not hold any literal variable. Meaning anything which is having `$` sign, otherwise you will not be allowed to create alerts on it. If `$` is used for variable get the equivalent expression of it and put it in your query

- For 2, Alert

![](https://i.imgur.com/WbfdOcY.png)

Specify the number for which the alert should be triggered.

![](https://user-images.githubusercontent.com/11765228/115198719-f79ba780-a124-11eb-9e43-508a4659c06a.png)

You may add any detail in the message for a particular alert and save it.

![](https://i.imgur.com/34OHpjS.png)
