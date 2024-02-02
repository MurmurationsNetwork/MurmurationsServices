# How to Create Alerts

Source: [Configuring PrometheusRules | Rancher](https://ranchermanager.docs.rancher.com/how-to-guides/advanced-user-guides/monitoring-v2-configuration-guides/advanced-configuration/prometheusrules)

This guide provides step-by-step instructions for creating custom alerts. Note that when you install monitoring through Rancher, it often comes with a set of default alerts. It's advisable to review these before creating your own to avoid duplications.

## General Procedure

1. Click **☰ > Cluster Management** in the upper left corner of the Rancher interface.
2. On the **Clusters** page, navigate to the desired cluster and select **Explore**.
3. In the *Cluster Dashboard* page, use the left sidebar to navigate to **Monitoring > Advanced** and choose **Prometheus Rules**.
4. Select **Create**.
5. In the bottom right corner, choose **Edit as YAML**.
6. Insert your configuration YAML and then click **Save**.

## Template for PrometheusRule

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: # [Insert a unique name for the PrometheusRule]
  namespace: cattle-monitoring-system
spec:
  groups:
    - name: # [Name of the group containing this set of rules]
      rules:
        - alert: # [Unique identifier for the alerting rule]
          annotations:
            description: # [Provide a detailed description of the alert]
            summary: # [A brief summary of the alert]
          expr: # [PromQL expression that triggers the alert]
          for: # [Time duration for the condition to be met before firing the alert, e.g., '5m' for 5 minutes]
          labels:
            severity: # [Indicate the severity of the alert, e.g., 'critical', 'warning']
```

## Examples

### High CPU Usage Alert

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: cpu-usage-alert
  namespace: cattle-monitoring-system
spec:
  groups:
    - name: system-resource-usage
      rules:
        - alert: HighCPUUsage
          annotations:
            description: This alert triggers when CPU usage exceeds 70% for more than 10 minutes.
            summary: Alert for high CPU usage on a cluster node.
          expr: 100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 70
          for: 10m
          labels:
            severity: warning
```
