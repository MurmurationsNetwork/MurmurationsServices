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
- [Step 3 - Customizing the kubeconfig File](#step-3---customizing-the-kubeconfig-file)
- [Step 4 - Merging Configuration Files for Cluster Management](#step-4---merging-configuration-files-for-cluster-management)
- [Conclusion](#conclusion)

## Prerequisites

Before starting, ensure you have:

1. An Ubuntu server set up and secured. For instructions, refer to "[How to Set Up and Secure Your Ubuntu Server](../01-setup-ubuntu/README.md)"
2. Terminal or SSH client access to your local machine.

## Step 1 - Installing and Verifying k3s

Install k3s on Ubuntu to kickstart your Rancher setup:

```bash
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.24.14+k3s1 sh -s - server --cluster-init
```

Verify the installation is successful:

```bash
k3s kubectl get nodes
k3s kubectl get pods --all-namespaces
```

Nodes should appear as `Ready`, with pods in `Running` or `Completed` states.

## Step 2 - Transferring the kubeconfig File

Transfer the kubeconfig file to manage the cluster remotely:

```bash
scp root@<ip_address>:/etc/rancher/k3s/k3s.yaml ~/.kube/k3s-murm-rancher.yaml
```

Replace `<ip_address>` with your server's actual IP. This step is crucial for remote cluster administration.

## Step 3 - Customizing the kubeconfig File

To better identify your cluster, let's customize your Kubernetes configuration file, `~/.kube/k3s-murm-rancher.yaml`. You can use the following steps:

Open the file using `vim`:

```bash
vim ~/.kube/k3s-murm-rancher.yaml
```

Modify the cluster name to "rancher" and update the server IP address to your Ubuntu server's IP address. For example:

```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ***
    server: https://<ip_address>:6443 # Remember to replace <ip_address> with your actual IP address
  name: rancher # Change the cluster name to rancher
contexts:
- context:
    cluster: rancher # Change the context's cluster to rancher
    user: default
  name: rancher # Change the context name to rancher
current-context: default
kind: Config
preferences: {}
users:
- name: default
  user:
    client-certificate-data: ***
    client-key-data: ***
```

Once you've made these changes, save and exit `vim`. Your kubeconfig file is now customized, and you can use the updated settings to manage your Kubernetes cluster.

## Step 4 - Merging Configuration Files for Cluster Management

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
```

Next, confirm your ability to connect to the k3s cluster while ensuring there's no disruption with the configurations of any existing clusters. Begin by listing the current contexts to view all available clusters:

```bash
kubectl config get-contexts
```

Switch to your k3s cluster's context to verify connectivity:

```bash
kubectl config use-context rancher
```

Check the resources across all namespaces to ensure the cluster is responsive:

```bash
kubectl get all --all-namespaces
```

Wrap up the process by deleting the now unnecessary k3s kubeconfig file and the backup of the original configuration:

```bash
rm ~/.kube/k3s-murm-rancher.yaml
rm ~/.kube/config_backup
```

## Conclusion

Your Ubuntu server is now equipped with k3s, creating a solid foundation for deploying Rancher.
