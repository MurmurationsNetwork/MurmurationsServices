# Setting Up Logging for Your Cluster

## Introduction

This guide will take you through the process of setting up logging for your environment using Elasticsearch and Kibana. By the end of this tutorial, you will have a powerful logging system in place, enabling efficient log management and analysis.

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Step 1 - Installing Elasticsearch and Kibana](#step-1---installing-elasticsearch-and-kibana)
- [Step 2 - Checking Installation Status](#step-2---checking-installation-status)
- [Step 3 - Testing Access to Kibana](#step-3---testing-access-to-kibana)
- [Conclusion](#conclusion)

## Prerequisites

Before starting, make sure you have:

1. Command-line access to your system.
2. Administrative privileges on the system.
3. The `kubectl` command-line tool configured for your cluster.

## Step 1 - Installing Elasticsearch and Kibana

Begin by deploying Elasticsearch and Kibana under the `murm-logging` namespace with the following command. This sets up the essential components for logging within your system.

```bash
make manually-deploy-murm-logging DEPLOY_ENV=production
```

This command kickstarts the deployment of Elasticsearch and Kibana, crucial for log aggregation, analysis, and visualization.

## Step 2 - Checking Installation Status

Following the installation, it's important to verify that everything is functioning as expected. Check the status of the pods in the `murm-logging` namespace by executing:

```bash
kubectl get pods -n murm-logging
```

This command will list all pods within the `murm-logging` namespace, allowing you to confirm the successful deployment of Elasticsearch and Kibana.

## Step 3 - Testing Access to Kibana

To access Kibana and start exploring your logs, use the `kubectl port-forward` command to forward a local port to the Kibana service within your cluster:

```bash
kubectl port-forward -n murm-logging svc/murm-logging-kibana 5601:5601
```

After running this command, you can access Kibana by navigating to `http://localhost:5601` in your web browser. This allows you to visualize and analyze logs gathered by Elasticsearch.

## Conclusion

You've now successfully set up a logging system with Elasticsearch and Kibana in your cluster. This system will aid in the aggregation, analysis, and visualization of logs, providing valuable insights into your system's operations.
