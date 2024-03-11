# Failover Procedures for Load Balancers

## Introduction

This document provides a detailed walkthrough for implementing failover procedures for load balancers.

## Table of Contents

- [Introduction](#introduction)
- [Step 1 - Attempt to Reboot the Server](#step-1---attempt-to-reboot-the-server)
- [Step 2 - Deploying a New VPS](#step-2---deploying-a-new-vps)
- [Step 3 - Configuring Nginx on the New VPS](#step-3---configuring-nginx-on-the-new-vps)
- [Conclusion](#conclusion)

## Step 1 - Attempt to Reboot the Server

Log into the server using the following command:

```bash
ssh root@your_load_balancer_node_ip
```

Once logged in, use the reboot command to restart the server:

```bash
sudo reboot
```

After the reboot, monitor the server's health and responsiveness to ensure it has returned to a healthy state.

## Step 2 - Deploying a New VPS

If the server does not return to a healthy state after a reboot, the next step is to prepare a new Virtual Private Server (VPS) for deployment.

## Step 3 - Configuring Nginx on the New VPS

After setting up the new VPS, configure Nginx to handle load balancing as outlined in the ["Set Up a Load Balancer for Murmuration Services"](./README.md) guide.

## Conclusion

By following these steps, you can swiftly respond to server failures and ensure your services remain available with minimal downtime.

Go back to [Home](../README.md).
