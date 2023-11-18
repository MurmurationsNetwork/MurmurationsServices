# Run Murmurations Services on Contabo

This guide is designed to assist you in deploying Murmurations Services on Contabo Virtual Private Server (VPS) instances.

## Terms and Definitions

Before we begin, let's familiarize ourselves with some key terms:

- **k8s (Kubernetes)**: An open-source system for automating the deployment, scaling, and management of containerized applications.
- **Rancher**: A comprehensive enterprise management platform for Kubernetes.
- **k3s**: A lightweight and optimized version of Kubernetes, tailored for Rancher.
- **rke2 (Rancher Kubernetes Engine 2)**: A Certified Kubernetes distribution tailored for enterprise use, accredited by the Cloud Native Computing Foundation (CNCF).
- **Longhorn**: A robust distributed block storage system for Kubernetes.

## Prerequisites

Ensure you have the following prerequisites ready:

- **VPS Instances**: Minimum of two with Ubuntu Server - one for hosting Rancher, and another as a Kubernetes cluster node using rke2. For extended capabilities, additional VPS instances may be used.
- **[kubectl](https://kubernetes.io/docs/tasks/tools/)**: Installed on your local machine.
- **[Helm](https://helm.sh)**: Also installed on your local machine.

## 1. Setting Up Your Ubuntu Server

### Establishing SSH Access

Begin by generating an SSH key pair, if you haven't already:

```bash
ssh-keygen -t rsa -b 4096
```

To transfer your public SSH key to your server, first display it:

```bash
cat ~/.ssh/id_rsa.pub
```

After copying the public key, connect to your VPS:

```bash
ssh root@<ip_address>
```

Upon login, ensure the `.ssh` directory exists with the correct permissions, and append your public key to the `authorized_keys` file:

```bash
mkdir -p ~/.ssh
chmod 700 ~/.ssh
vim ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
```

Next, modify the SSH daemon configuration to enhance security. Open the configuration file and make the following changes:

```bash
vim /etc/ssh/sshd_config

# Modify these lines or add them if they're not present
PasswordAuthentication no
ChallengeResponseAuthentication no
UsePAM no
PermitRootLogin prohibit-password
```

Restart the SSH service to apply these settings:

```bash
systemctl restart ssh
```

You can now use SSH key-based authentication for server access:

```bash
# Please open a new terminal to test connection first.
ssh root@<ip_address>
```

### Uploading and Executing Setup Scripts

Transfer the setup script to your server:

```bash
cd /path/to/MurmurationsServices
scp scripts/ubuntu_setup.sh root@<ip_address>:
ssh root@<ip_address>
```

Run the setup script:

```bash
chmod +x ubuntu_setup.sh && ./ubuntu_setup.sh
```

## 2. Installing Rancher on a k3s Cluster

This section guides you through installing k3s on an Ubuntu server and then deploying Rancher to manage your Kubernetes clusters.

### Installing k3s on Ubuntu

To install k3s on your Ubuntu server, execute the following command. Make sure to replace `<VERSION>` with the desired version of k3s, such as `v1.26.10+k3s1`.

```bash
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=<VERSION> sh -s - server --cluster-init
```

**Important Note**: Rancher necessitates a compatible Kubernetes version. Refer to the [Rancher Support Matrix](https://rancher.com/support-maintenance-terms/) to verify version compatibility.

**Update (2023-11-07)**: As of this writing, version `v1.26.10+k3s1` is recommended. Consult the [official k3s GitHub releases](https://github.com/k3s-io/k3s/tags) for the latest version updates.

### Transferring the kubeconfig File to Your Workstation

To remotely manage your cluster, transfer the `kubeconfig` file to your local machine:

```bash
# Open a new tab and execute in your local machine.
scp root@<ip_address>:/etc/rancher/k3s/k3s.yaml ~/.kube/<config_name>
```

Substitute `<ip_address>` with your server's IP and `<config_name>` with a preferred name for your configuration file.

The updated section for "Merge config together" with a consistent writing style is as follows:

### Merging Configuration Files

To effectively manage multiple Kubernetes clusters, it's essential to merge the kubeconfig files.

1. **Set the KUBECONFIG environment variable:** This step combines your current kubeconfig file with the new configuration file you've acquired from your k3s cluster. Replace `<config_name>` with the name of your new configuration file.

    ```bash
    export KUBECONFIG=~/.kube/config:~/.kube/<config_name>.yaml
    ```

2. **Create a unified kubeconfig file:** This command merges the configurations into a single file while ensuring all data is intact and correctly formatted.

    ```bash
    kubectl config view --merge --flatten > ~/.kube/merged_kubeconfig
    ```

3. **Backup the original config file:** It's a good practice to keep a backup of your original kubeconfig file in case you need to revert to the previous settings.

    ```bash
    mv ~/.kube/config ~/.kube/config_backup
    ```

4. **Replace the current config with the merged file:** This step finalizes the process by replacing the existing kubeconfig file with the newly merged file.

    ```bash
    mv ~/.kube/merged_kubeconfig ~/.kube/config
    ```

### Updating the Rancher Server URL in kubeconfig

Modify your kubeconfig file to reflect the correct server URL:

```bash
vim ~/.kube/<config_name>
```

In the file, update the `server` field to `https://<ip_address>:6443`.

### Deploying Rancher with Helm

Proceed with adding the Rancher Helm chart repository:

```bash
# Make sure you have switched to the correct context.
kubectl config use-context <context_name>

helm repo add rancher-latest https://releases.rancher.com/server-charts/latest
```

Create the `cattle-system` namespace for Rancher:

```bash
kubectl create namespace cattle-system
```

Before installing Rancher, set up Cert-Manager's Custom Resource Definitions (CRDs):

You can find the the latest version by checking their [GitHub tags page](https://github.com/cert-manager/cert-manager/tags).

```bash
VERSION=<version>

kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/${VERSION}/cert-manager.crds.yaml
```

Add the Jetstack Helm repository, which provides Cert-Manager:

```bash
helm repo add jetstack https://charts.jetstack.io
helm repo update
```

Install Cert-Manager in its dedicated namespace:

```bash
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace
```

Now, deploy Rancher using Helm in the `cattle-system` namespace. Replace `<ip_address>` with your server's IP and `<PASSWORD_FOR_RANCHER_ADMIN>` with your chosen password:

```bash
helm install rancher rancher-latest/rancher \
  --namespace cattle-system \
  --set hostname=<ip_address>.sslip.io \
  --set replicas=1 \
  --set bootstrapPassword=<PASSWORD_FOR_RANCHER_ADMIN>
```

**Note**: For comprehensive instructions and updated guides, refer to the [official Rancher documentation](https://ranchermanager.docs.rancher.com/getting-started/quick-start-guides/deploy-rancher-manager/helm-cli).

## 3. Creating a Kubernetes Cluster

Once Rancher is installed, you can access its dashboard and start building your Kubernetes cluster.

### Accessing the Rancher Dashboard

1. Launch a web browser and navigate to `http://<ip_address>.sslip.io`. Replace `<ip_address>` with the IP address of your Linux node where Rancher is hosted.
2. Log in using the credentials set during the Rancher installation.

### Initiating Cluster Creation

1. Within the Rancher dashboard, go to the **Cluster Management** tab.
2. Click **Create** and select the "Custom" option to commence the creation of a new cluster.
3. Provide a descriptive name for your cluster in the designated field.
4. Proceed by clicking **Create** to move to the next step.

![](https://github.com/MurmurationsNetwork/MurmurationsServices/assets/11765228/f7e02619-48ed-4436-b690-0541920ee73f)
![](https://github.com/MurmurationsNetwork/MurmurationsServices/assets/11765228/a7dabc2b-cbab-4e7b-b6e4-59b9cf1178ee)
![](https://github.com/MurmurationsNetwork/MurmurationsServices/assets/11765228/b5871018-6c10-4a14-813b-6568dcf6b1a1)

### Registering Nodes to the Cluster

1. In the new cluster's configuration page, switch to the **Registration** tab.
2. If you're not using certificate-based authentication between your nodes and the Rancher server, choose the "Insecure" option in Step 2 of the registration process.
3. The command displayed is what you’ll need to execute on your other Ubuntu servers that you intend to include in this Kubernetes cluster.

![](https://github.com/MurmurationsNetwork/MurmurationsServices/assets/11765228/76d4c966-c7f1-4712-bcf7-ab22721fcd01)

### Adding Nodes to Your Cluster

1. Copy the command provided in the Rancher dashboard.
2. Connect to your other server intended to be a Kubernetes node:

   ```bash
   ssh root@<ANOTHER_ip_address>
   ```

3. Paste and execute the copied command in the server's terminal.

Following these steps will integrate your new Ubuntu server into the Kubernetes cluster managed by Rancher.

### Downloading the Cluster kubeconfig File

To interact with your cluster using `kubectl` from your local machine, you’ll need the cluster's kubeconfig file.

1. In the Rancher dashboard, go to the **Cluster Management** tab.
2. Find and click on the `...` button next to the Explore button.
3. From the menu, select **Download Kubeconfig** to obtain the kubeconfig file.

![](https://github.com/MurmurationsNetwork/MurmurationsServices/assets/11765228/b6a167ee-f750-4de4-b968-fd01cfe95f28)

Follow [Merging Configuration Files](#merging-configuration-files) to merge configs together.

Here’s an example of how to use the kubeconfig file:

```bash
# Make sure you have switched to the correct context.
kubectl config use-context <context_name>

kubectl get nodes
```

## 4. Installing Longhorn

Longhorn is a critical component for managing persistent volumes in Kubernetes, enabling automatic volume provisioning, replication, and backup.

### Logging into Rancher

1. Access your Rancher dashboard.
2. Select the desired cluster in the "Explore Cluster" section.

![](https://github.com/MurmurationsNetwork/MurmurationsServices/assets/11765228/a7457f1e-fc98-4470-94d1-aeff9e3d0019)

### Installing Longhorn

1. Inside the cluster view, click on "Cluster Tools."
2. Locate Longhorn in the tools list and initiate its installation.

![](https://github.com/MurmurationsNetwork/MurmurationsServices/assets/11765228/08d080b6-5a72-4c8b-b5e5-bf924bd49ecc)


### Customizing Longhorn Installation

1. Optionally, before installing, you can customize settings by clicking on "Customize Helm." For instance, you might want to set the storage class's retain policy to 'Retain'.
2. Proceed with the installation of Longhorn.

### Accessing Longhorn

Once installed, Longhorn can be accessed from the left sidebar in the Rancher dashboard.

![](https://github.com/MurmurationsNetwork/MurmurationsServices/assets/11765228/cfb6d78a-9b90-421d-b731-ece597905afc)

## 5. Setting Up HTTPS

Configuring HTTPS is crucial for securing your services. This involves deploying an ingress controller and cert-manager for automatic TLS certificate management. Before proceeding, ensure your DNS records are correctly set up.

### Configuring DNS Records

1. Log in to your DNS provider's management console.
2. Locate the DNS record management section.
3. Create A records for your services, pointing them to the external IP address of your Linux node hosting the ingress controller. For example:

   ```plaintext
   index.murmurations.network A <ip_address>
   library.murmurations.network A <ip_address>
   data-proxy.murmurations.network A <ip_address>
   ```

   Replace `<ip_address>` with your node's IP and adjust domain names according to your setup.

### Installing cert-manager with Helm

After setting up DNS, proceed to install cert-manager in your Kubernetes cluster:

```shell
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --set installCRDs=true
```

### Deploying Ingress

Navigate to the `charts/murmurations/charts/ingress` directory and update the `certificate` and `ingress` charts as needed.

Currently, we have 4 environments: `production`, `staging`, `proto` and `development`.

```shell
# Make sure you have switched to the correct context.
kubectl config use-context <context_name>

# Replace '<environment>' with your desired environment.
make manually-deploy-ingress ENV=<environment>
```

## 6. Deploying Murmurations Services

This section details the process of deploying Murmurations services in your Kubernetes environment.

### Creating Required Secrets

Before deploying the services, you need to set up necessary secrets:

1. Follow the guidelines provided in [Creating Secrets for Murmurations Services](../secrets.md) to properly configure each service.
2. Ensure that all secrets are created and stored as per the instructions.

### Deploying the Services

With the secrets in place, you can now proceed to deploy Murmurations services.

```shell
# Make sure you have switched to the correct context.
kubectl config use-context <context_name>

# Replace '<environment>' with your desired deployment environment.
make deploy-all-services ENV=<environment>
```

This step will initiate the deployment process for Murmurations services, ensuring they are correctly configured and deployed in your chosen environment.
