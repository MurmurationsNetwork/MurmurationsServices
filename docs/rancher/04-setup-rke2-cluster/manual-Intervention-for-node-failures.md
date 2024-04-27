# Manual Intervention for Node Failures

## Introduction

This document provides a step-by-step guide for manually handling node failures within Kubernetes clusters managed by Rancher. Node failures can disrupt normal operations, and timely intervention is crucial to maintain the availability and performance of your applications.

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Step 1 - Forcing Pod Recreation on Other Nodes](#step-1---forcing-pod-recreation-on-other-nodes)
- [Step 2 - Removing the Unavailable Machine](#step-2---removing-the-unavailable-machine)
- [Step 3 - Adding a New Node Using a Temporary Server](#step-3---adding-a-new-node-using-a-temporary-server)
- [Conclusion](#conclusion)

## Prerequisites

Before proceeding with this guide, ensure that you have administrative privileges on the Rancher platform and that you have set up RKE2 and access to the Rancher cluster with command-line tools.

## Step 1 - Forcing Pod Recreation on Other Nodes

When a node fails, some pods may remain in a terminating state indefinitely. To address this, force the scheduler to recreate these pods on other available nodes using the following command:

```bash
kubectl get pods --all-namespaces | grep Terminating | awk '{print $1 " " $2}' | while read ns pod; do kubectl delete pod $pod -n $ns --grace-period=0 --force; done
```

This command identifies all pods that are stuck in a Terminating state across all namespaces and forcefully deletes them, prompting Kubernetes to recreate them on other nodes.

## Step 2 - Removing the Unavailable Machine

Next, navigate to the Rancher UI to remove the node that has become unavailable. Here’s how:

1. Go to the ☰ menu and select Cluster Management.
2. Locate the cluster containing the failed node.
3. Select the node in question and use the option to delete it, removing the unavailable machine from your cluster.

This step ensures that the cluster's resources are updated and that the failed node is no longer considered part of the cluster.

## Step 3 - Adding a New Node Using a Temporary Server

To maintain the desired capacity of your cluster, you can quickly add a new node using one of the available temporary servers. Following the instruction provided in the [documentation](./README.md#step-3---registering-nodes-to-the-cluster), you can add a new node to your cluster and ensure that the applications continue to run smoothly.

## Conclusion

By following these steps, you can manually intervene in the event of node failures to maintain the resilience and efficiency of your Kubernetes clusters. Rancher's tools facilitate quick recovery actions, ensuring that your environments remain robust and your applications continue to operate smoothly.
