# Rancher-Managed Kubernetes for Developers

This guide offers developers a comprehensive walkthrough on initiating an operations-ready Kubernetes cluster managed by Rancher. Unlike traditional Kubernetes setups, Rancher simplifies the orchestration, ensuring developers can deploy consistent tooling and configurations across diverse cloud environments. This document outlines the necessary steps and tools to achieve an operationally ready state with Rancher-managed Kubernetes.

## Table of Contents

1. [Set Up Ubuntu Server](01-setup-ubuntu/README.md)
2. [Set Up K3s for Rancher](02-setup-k3s/README.md)
    - [Upgrade K3s](02-setup-k3s/upgrade-k3s.md)
3. [Set Up Rancher](03-setup-rancher/README.md)
    - [Backup Rancher](./03-setup-rancher/backup-rancher.md)
    - [Disaster Recovery](03-setup-rancher/disaster-recovery-rancher.md)
    - [Restore Rancher](03-setup-rancher/restore-rancher.md)
    - [Upgrade Rancher](03-setup-rancher/upgrade-rancher.md)
4. [Set Up RKE2 Cluster for Murmuration Services](04-setup-rke2-cluster/README.md)
    - [Manage Cluster Nodes](04-setup-rke2-cluster/manage-cluster-nodes.md)
    - [Upgrade Cluster](04-setup-rke2-cluster/upgrade-rk2-cluster.md)
5. [Set Up a Load Balancer](05-setup-lb/README.md)
    - [Failover Procedure for Load Balancers](05-setup-lb/failover-procedure-for-load-balancer.md)
6. [Set Up Longhorn](06-setup-longhorn/README.md)
    - [Set Up Recurring Backup](06-setup-longhorn/recurring-backup.md)
    - [Recover from Backup](06-setup-longhorn/recover-from-backup.md)
7. [Run Murmuration Services](07-run-murmuration-services/README.md)
    - [Migrate MongoDB](07-run-murmuration-services/migrate-mongodb.md)
    - [Migrating Elasticsearch](07-run-murmuration-services/migrate-es.md)
8. [Set Up Monitoring](08-setup-monitoring/README.md)
    - [Receive Alerts](08-setup-monitoring/how-to-receive-alerts.md)
9. [Set Up Logging](./09-setup-logging/README.md)
