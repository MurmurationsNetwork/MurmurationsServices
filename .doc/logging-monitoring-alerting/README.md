# Logging, Monitoring and Alerting

# Prerequisite

## 1. Add Helm Charts

```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add kube-state-metrics https://kubernetes.github.io/kube-state-metrics
helm repo add grafana https://grafana.github.io/helm-charts
helm repo add loki https://grafana.github.io/loki/charts
helm repo update
```

## 2. Create a Name Space

```
kubectl create namespace kube-monitoring
```

# Logging

## 1. Install Loki

```
helm upgrade loki grafana/loki -f .doc/logging-monitoring-alerting/loki.values.yaml -n kube-monitoring --install
```

## 2. Install Promtail

```
helm upgrade promtail grafana/promtail -f .doc/logging-monitoring-alerting/promtail.values.yaml -n kube-monitoring --install
```

After deployment do port-forward for promtail

```
kubectl port-forward daemonset/promtail 3101 --namespace kube-monitoring
```

on http://localhost:3101/targets

You will see all the pods which are getting scraped by promtail for logs.


# Monitoring

## 1. Setup Webhook API Slack

Goto this link https://api.slack.com/messaging/webhooks and create a slack app.

![image](https://user-images.githubusercontent.com/11765228/114982939-9ffe0180-9ec2-11eb-9d45-4da79125951f.png)

After creating the slack app, click on incoming webhooks.

![image](https://user-images.githubusercontent.com/11765228/114983126-ce7bdc80-9ec2-11eb-8ac0-0240e045d164.png)

![image](https://user-images.githubusercontent.com/11765228/114983254-f5d2a980-9ec2-11eb-80d2-8fcaf2e53ad8.png)

After clicking the add new webhook to workspace.

Select the channel.

![image](https://user-images.githubusercontent.com/11765228/114983572-5d88f480-9ec3-11eb-9733-43e92c734eab.png)

After selecting and allowing it.

![image](https://user-images.githubusercontent.com/11765228/114983936-c7090300-9ec3-11eb-95ca-8b557fa7ac2a.png)

Test it with curl whether you are receiving the slack notification or not.

Copy the `Webhook URL` and then use it wherever it is needed below.

## 2. Deploy Prometheus

**Update slack-notification.yaml**

Replace `<WEBHOOK_URL>` and `<CHANNEL_NAME>` in .doc/logging-monitoring-alerting/slack-notification.yaml

**Deploy Prometheus**

```
helm upgrade prometheus prometheus-community/prometheus -f .doc/logging-monitoring-alerting/prometheus.values.yaml -f .doc/logging-monitoring-alerting/slack-notification.yaml -n kube-monitoring --install
```

## 3. Deploy Grafana

**Update Grafana Password**

Replace `<ADMIN_PASSWORD>` in .doc/logging-monitoring-alerting/prometheus.values.yaml

**Deploy Grafana**

```
helm upgrade grafana grafana/grafana -f .doc/logging-monitoring-alerting/grafana.values.yaml -n kube-monitoring  --install
```

## 4. Setup the Monitoring and Dashboards

```
kubectl port-forward svc/grafana 3000:80 -n kube-monitoring
```

**Setup Grafana**

Navigate to http://localhost:3000

**Add Alert Webhook**

Adding alert for certain things from dashboard.

![image](https://user-images.githubusercontent.com/11765228/115104231-9ab7ba00-9f89-11eb-9d06-cf4d592b1b03.png)

Add Webhook Url

![image](https://user-images.githubusercontent.com/11765228/115104198-6a701b80-9f89-11eb-8d5e-4d69e1446d03.png)

Check these boxes as well

![image](https://user-images.githubusercontent.com/11765228/115104205-79ef6480-9f89-11eb-804c-c8e1828ccca1.png)

## 5. Configure the Grafana Dashboard


![image](https://user-images.githubusercontent.com/11765228/115194754-780bd980-a120-11eb-9284-c01458983f6b.png)

Type 1860 for Node Exporter Full Dashboard

![image](https://user-images.githubusercontent.com/11765228/115194824-907bf400-a120-11eb-86f3-68d06aa5ffcd.png)

![image](https://user-images.githubusercontent.com/11765228/115195026-dafd7080-a120-11eb-89ac-2af4e5120ea1.png)

Now again add one more dashboard: 8685

![image](https://user-images.githubusercontent.com/11765228/115195120-f49eb800-a120-11eb-971a-993c668e6af4.png)

# Alerting

If want to add alert on pre-built specific property.

![](https://i.imgur.com/aXYWiPy.png)

After clicking on edit, you will reach to this page. 

![](https://i.imgur.com/wo4GiAM.png)

*Note Alert option will only appear in Graph visualization only*

- For 1, You may use any [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/)

Note: PromQL should not hold any literal variable. Meaning anything which is having `$` sign, otherwise you will not be allowed to create alerts on it. If `$` is used for variable get the equivalent expression of it and put it in your query

- For 2, Alert

![](https://i.imgur.com/WbfdOcY.png)

Specify the number for which the alert should be triggered.

![](https://user-images.githubusercontent.com/11765228/115198719-f79ba780-a124-11eb-9e43-508a4659c06a.png)

You may add any detail in the message for particular message and save it.

![](https://i.imgur.com/34OHpjS.png)
