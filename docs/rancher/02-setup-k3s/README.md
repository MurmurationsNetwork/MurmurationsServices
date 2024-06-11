# Setting Up k3s for Rancher Integration

## Introduction

Prepare your Ubuntu server for Rancher deployment by installing k3s, a lightweight Kubernetes distribution ideal for simplified cluster management. This setup ensures minimal resource usage while providing a robust environment for Rancher.

After completing this guide, you will have:

- Installed k3s on your Ubuntu server, ready for Rancher.
- Configured your local environment for seamless k3s cluster management.

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Step 1 - Installing and Verifying k3s](#step-1---installing-and-verifying-k3s)
- [Step 2 - Transferring the kubeconfig File](#step-2---transferring-the-kubeconfig-file)
- [Step 3 - Setting Up DNS for the Server](#step-3---setting-up-dns-for-the-server)
- [Step 4 - Customizing the kubeconfig File](#step-4---customizing-the-kubeconfig-file)
- [Step 5 - Merging Configuration Files for Cluster Management](#step-5---merging-configuration-files-for-cluster-management)
- [Conclusion](#conclusion)

## Prerequisites

Before starting, ensure you have:

1. An Ubuntu server set up and secured. For instructions, refer to "[How to Set Up and Secure Your Ubuntu Server](../01-setup-ubuntu/README.md)"
2. Terminal or SSH client access to your local machine.

## Step 1 - Installing and Verifying k3s

Before initiating the installation of k3s, it's imperative to establish a secure connection to your Ubuntu server:

```bash
ssh root@{{ip_address}}
```

Ensure to replace `{{ip_address}}` with the actual IP address of your server.

Once a secure connection is established, you can proceed with the k3s installation. It's crucial to consult the [SUSE Rancher support matrix](https://www.suse.com/suse-rancher/support-matrix/all-supported-versions/rancher-v2-8-1/) to verify compatibility with your intended k3s version. The highest supported version at this time is v1.27. For a detailed list of available versions, visit the [list of k3s versions](https://github.com/k3s-io/k3s/tags).

To install k3s, execute the following command, ensuring you adjust the version number as needed to match the version compatible with your setup:

```bash
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.27.14+k3s1 sh -s - server --cluster-init
```

After the installation completes, verify that your cluster is operational:

```bash
k3s kubectl get nodes
k3s kubectl get pods --all-namespaces
```

You should see nodes marked as `Ready`, and pods should be in `Running` or `Completed` states, indicating that k3s is installed and functioning correctly on your server.

## Step 2 - Transferring the kubeconfig File

Transfer the kubeconfig file to manage the cluster remotely:

```bash
scp root@{{ip_address}}:/etc/rancher/k3s/k3s.yaml ~/.kube/k3s-murm-rancher.yaml
```

Replace `{{ip_address}}` with your server's actual IP. This step is crucial for remote cluster administration.

## Step 3 - Setting Up DNS for the Server

Before customizing your kubeconfig file, ensure your server is accessible via a DNS name. This involves adding an A record in your DNS management interface pointing to your server’s IP address.

- Log in to your DNS provider's management console.
- Navigate to the DNS settings section for your domain.
- Add an A record with the following details:
  - **Host:** The subdomain or prefix you wish to use (e.g., `k3s` or `rancher`).
  - **Points to:** The IP address of your Ubuntu server.
  - **TTL:** For A Record, 3600 seconds is recommended.

This DNS setup facilitates access to your server using a memorable URL instead of an IP address. Most importantly, if DNS is not used, the backup and restore processes will not work if the IP address changes.

## Step 4 - Customizing the kubeconfig File

In this step we will customize the kubeconfig file for easier identification and management of the cluster.

First, update the server URL in the kubeconfig file to point to your server's DNS name. Ensure you replace `{{server_dns_name}}` with your server's DNS name:

```bash
perl -pi -e "s/127.0.0.1/{{server_dns_name}}/g" ~/.kube/k3s-murm-rancher.yaml
```

Next, change all instances of the default context, cluster, and user names to `k3s-murm-rancher` to reflect your specific Rancher setup. This naming convention makes managing multiple clusters more intuitive:

```bash
perl -pi -e "s/default/k3s-murm-rancher/g" ~/.kube/k3s-murm-rancher.yaml
```

## Step 5 - Merging Configuration Files for Cluster Management

To streamline the management of multiple Kubernetes clusters, integrate your kubeconfig files:

Set the `KUBECONFIG` environment variable to encompass both your original and the new k3s kubeconfig files:

```bash
export KUBECONFIG=~/.kube/config:~/.kube/k3s-murm-rancher.yaml
```

Then, merge these configurations into one unified file:

```bash
kubectl config view --merge --flatten > ~/.kube/merged_kubeconfig
```

Proceed by backing up your original config file. Afterward, replace it with the newly merged configuration:

```bash
mv ~/.kube/config ~/.kube/config_backup
mv ~/.kube/merged_kubeconfig ~/.kube/config
chmod 600 ~/.kube/config
```

Next, confirm your ability to connect to the k3s cluster while ensuring there's no disruption with the configurations of any existing clusters. Begin by listing the current contexts to view all available clusters:

```bash
kubectl config get-contexts
```

Switch to your k3s cluster's context to verify connectivity:

```bash
kubectl config use-context k3s-murm-rancher
```

Check the resources across all namespaces to ensure the cluster is responsive:

```bash
kubectl get all --all-namespaces
```

Wrap up the process by deleting the unnecessary k3s kubeconfig file:

```bash
rm ~/.kube/k3s-murm-rancher.yaml
```

## Conclusion

Your Ubuntu server is now equipped with k3s, creating a solid foundation for deploying Rancher.

Go to Section 3 - [Setting up Rancher](../03-setup-rancher/README.md).
