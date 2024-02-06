# Setting Up a Load Balancer for Murmuration Services

## Introduction

This guide outlines the steps for setting up a load balancer for Murmuration services, leveraging Docker and Nginx. A load balancer efficiently distributes incoming network traffic across multiple servers, enhancing the reliability and availability of your services. Upon completing this guide, you will have a fully operational load balancing setup for Murmuration services.

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Step 1 - Preparing the Load Balancer Server](#step-1---preparing-the-load-balancer-server)
- [Step 2 - Installing Docker](#step-2---installing-docker)
- [Step 3 - Setting Up Nginx as a Load Balancer](#step-3---setting-up-nginx-as-a-load-balancer)
- [Step 4 - Launching the Nginx Load Balancer](#step-4---launching-the-nginx-load-balancer)
- [Step 5 - Integrating the Load Balancer with Kubernetes](#step-5---integrating-the-load-balancer-with-kubernetes)
- [Conclusion](#conclusion)

## Prerequisites

Before beginning, ensure you have:

- A server designated for the load balancer.
- Basic understanding of Nginx configurations.
- Access to a Kubernetes cluster via `kubectl`.

## Step 1 - Preparing the Load Balancer Server

Identify a server to serve as your load balancer. For instructions on setting up a server, see the [Ubuntu Server Setup Guide](../01-setup-ubuntu/README.md).

## Step 2 - Installing Docker

Docker streamlines the deployment and management of containerized applications. To install Docker on your load balancer server, follow these steps:

First, update your package list:

```bash
sudo apt update
```

Then, install the necessary packages for HTTPS repository access:

```bash
sudo apt install apt-transport-https ca-certificates curl software-properties-common
```

Next, add the Docker repository GPG key:

```bash
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
```

Include the Docker repository:

```bash
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

Update the package list again:

```bash
sudo apt update
```

Finally, install Docker:

```bash
sudo apt install docker-ce
```

Verify Docker installation by checking its status:

```bash
sudo systemctl status docker
```

Docker will now be installed and configured to start on boot.

## Step 3 - Setting Up Nginx as a Load Balancer

To configure Nginx as a load balancer, create a configuration file in `/etc/nginx.conf` to handle your traffic:

First, open the configuration file for editing:

```bash
vim /etc/nginx.conf
```

Then, insert the following configuration, replacing `<server-address>` with the actual addresses of your services:

```nginx
worker_processes auto;
worker_rlimit_nofile 40000;

events {
    worker_connections 8192;
}

stream {
    upstream murmur_servers_http {
        least_conn;
        server <server-address-1>:80 max_fails=3 fail_timeout=5s;
        server <server-address-2>:80 max_fails=3 fail_timeout=5s;
        server <server-address-3>:80 max_fails=3 fail_timeout=5s;
    }

    upstream murmur_servers_https {
        least_conn;
        server <server-address-1>:443 max_fails=3 fail_timeout=5s;
        server <server-address-2>:443 max_fails=3 fail_timeout=5s;
        server <server-address-3>:443 max_fails=3 fail_timeout=5s;
    }

    server {
        listen 80;
        proxy_pass murmur_servers_http;
        proxy_protocol on;
    }

    server {
        listen 443;
        proxy_pass murmur_servers_https;
        proxy_protocol on;
    }
}
```

Note, if you only have one node for your RKE2 cluster, then you only need to have one server in murmur_servers_http and murmur_servers_https.

## Step 4 - Launching the Nginx Load Balancer

Deploy Nginx within Docker to initiate the load balancing functionality:

First, execute the following command to start Nginx as a load balancer:

```bash
docker run -d --restart=unless-stopped \
-p 80:80 -p 443:443 \
-v /etc/nginx.conf:/etc/nginx/nginx.conf \
nginx:1.14
```

## Step 5 - Integrating the Load Balancer with Kubernetes

Configure Kubernetes to direct traffic through your new load balancer:

First, switch to the appropriate Kubernetes context:

```bash
kubectl config use-context <your-cluster-name>
```

Then, modify the ingress-nginx controller's configuration to enable proxy protocol support:

```bash
kubectl edit configmaps rke2-ingress-nginx-controller -n kube-system
```

Add `use-proxy-protocol: "true"` directly under the data section of the config map. Ensure you do not modify any other parts of the configuration.

```yaml
data:
  use-proxy-protocol: "true"
```

## Conclusion

You have successfully established a load balancer for Murmuration services using Docker and Nginx. This setup enhances the distribution and reliability of your service traffic, ensuring high availability and performance.

Remember to replace placeholders with actual data relevant to your configuration and validate the setup by accessing your services through the load balancer's IP or domain name.
