# Restore Rancher

## Introduction

This document provides a detailed guide on how to restore your Rancher deployment.

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Step 1 - Navigate to the Rancher Dashboard](#step-1---navigate-to-the-rancher-dashboard)
- [Step 2 - Executing the Restore](#step-2---executing-the-restore)
- [Step 3 - Monitoring Restore Progress](#step-3---monitoring-restore-progress)
- [Conclusion](#conclusion)

## Prerequisites

Before beginning the restoration process, ensure you have:

- Administrative access to the Rancher dashboard.
- A backup file available for restoration.
- Ensure both credentialSecret and encryptionConfigSecret are available in your k3s cluster. For instructions on creating a credentialSecret, check [here](./install-rancher-backup-tool.md#step-2---obtaining-s3-object-storage-credentials). For encryptionConfigSecret, check [here](./backup-rancher.md#step-3---creating-the-encryption-provider-configuration-file).

## Step 1 - Navigate to the Rancher Dashboard

Access the Rancher dashboard and locate the **Rancher Backups > Restores** section from the left navigation bar to manage and initiate restoration processes.

![Rancher Restore Access](./assets/images/rancher-restore-access.png)

## Step 2 - Executing the Restore

Select **Edit as YAML** and enter the following Restore YAML configuration:

![Create Rancher Restore Page](./assets/images/create-rancher-restore-page.png)

```yaml
apiVersion: resources.cattle.io/v1
kind: Restore
metadata:
  name: restore-migration
spec:
  backupFilename: <backup file name>
  encryptionConfigSecretName: encryptionconfig
  storageLocation:
    s3:
      credentialSecretName: contabo-s3-creds
      credentialSecretNamespace: default
      bucketName: rancher-backups
      region: EU
      endpoint: eu2.contabostorage.com
```

Replace the placeholder values with your actual backup details and storage location information. Click **Create** to proceed with the restoration. Note: `contabo-s3-creds` is created during the [obtaining S3 object storage credentials](./install-rancher-backup-tool.md#step-2---obtaining-s3-object-storage-credentials) step, and `encryptionconfig` is created during the [creating the encryption provider configuration file](./backup-rancher.md#step-3---creating-the-encryption-provider-configuration-file) step.

## Step 3 - Monitoring Restore Progress

To monitor the restoration progress, run the following command to check the logs of the backup operator:

```bash
kubectl config use-context k3s-murm-rancher
kubectl logs -n cattle-resources-system -l app.kubernetes.io/name=rancher-backup -f
```

This command provides real-time insights into the restoration process.

## Conclusion

Following this guide ensures a structured and efficient approach to restoring your Rancher-managed Kubernetes environments.
