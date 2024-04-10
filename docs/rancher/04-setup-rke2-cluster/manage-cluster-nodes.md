# Manage Cluster Nodes

## Introduction

This guide outlines the management of cluster nodes within Rancher. By following the document, you can effectively add and remove nodes.

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Step 1 - Select Cluster](#step-1---select-cluster)
- [Step 2 - Go to Nodes Page](#step-2---go-to-nodes-page)
- [Step 3 - Select and Drain Nodes](#step-3---select-and-drain-nodes)
- [Step 4 - Remove Drained Nodes](#step-4---remove-drained-nodes)
- [Step 5 - Remove Nodes from Cluster Management Page](#step-5---remove-nodes-from-cluster-management-page)
- [Step 6 - Rename Servers Accordingly](#step-6---rename-servers-accordingly)
- [Conclusion](#conclusion)

## Prerequisites

- **Important**: Before starting, you must add an equal number of new nodes to your cluster as the ones you plan to drain. This step is critical to keep the total number of nodes unchanged, ensuring the stability and capacity of your cluster are maintained. To add nodes to your cluster, please follow this [guide](./README.md#step-3---registering-nodes-to-the-cluster).

## Step 1 - Select Cluster

Navigate to the Rancher homepage and choose the cluster you wish to manage.

![Selecting a Cluster in Rancher](./assets/images/rancher-cluster-selection.png)

## Step 2 - Go to Nodes Page

From the cluster homepage, access the nodes section via the left sidebar.

![Accessing Nodes Page](./assets/images/rancher-go-to-nodes-page.png)

## Step 3 - Select and Drain Nodes

On the Nodes page, identify and select the nodes for drainage, then click "Drain".

![Selecting Nodes for Drainage](./assets/images/rancher-select-and-drain-nodes.png)

Choose "Delete Empty Dir Data" and confirm by clicking "Drain".

![Configuring Drainage Options](./assets/images/rancher-drain-config.png)

## Step 4 - Remove Drained Nodes

Following successful drainage, eliminate the drained nodes from the Nodes page.

![Removing Drained Nodes from Management](./assets/images/rancher-removed-drained-node.png)

## Step 5 - Remove Nodes from Cluster Management Page

Eliminate nodes from the cluster management interface to prevent Rancher from repeatedly trying to connect.

First, select the cluster from the cluster management page:

![Rancher Cluster Management Page](./assets/images/rancher-cluster-management-page.png)

Secondly, select the node you removed in step 4 and delete it:

![Rancher Delete Drained Nodes](./assets/images/rancher-delete-drained-nodes.png)

## Step 6 - Rename Servers Accordingly

Make sure to rename the servers at your VPC provider to ensure they can be reused or deleted later on.

## Conclusion

This guide provides a streamlined approach to managing cluster nodes within Rancher, ensuring operational efficiency and reliability.

Go back to [Home](../README.md).
