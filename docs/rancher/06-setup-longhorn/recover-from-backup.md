# Recover from Backup

## Introduction

This document provides a detailed guide on how to recover your data from backups using Longhorn. It is designed to help users restore their services and volumes efficiently.

## Table of Contents

- [Introduction](#introduction)
- [Step 1 - Access the Longhorn Backup Page](#step-1---access-the-longhorn-backup-page)
- [Step 2 - Recover Backups](#step-2---recover-backups)
- [Step 3 - Bind Volumes](#step-3---bind-volumes)
- [Step 4 - Deploy Services](#step-4---deploy-services)

## Step 1 - Access the Longhorn Backup Page

Navigate to the Longhorn backup page to see a comprehensive list of all available backups.

![Available Backups in Longhorn](./assets/images/longhorn-backups.png)

## Step 2 - Recover Backups

1. Navigate to the Backup tab.

2. Select the backups you wish to recover.

Note: If you want to recover an entire service, you must choose all the related backups.

![Recover Backup - Selection](./assets/images/longhorn-recover-backup-selection.png)

3. For "Recovery Options", select "Read-Write" mode, then click "OK".

![Recover Backup - Options](./assets/images/longhorn-recover-backup-options.png)

## Step 3 - Bind Volumes

First, navigate to the Volumes tab.

Initially, the volume's state will be "Detached", and PV/PVC will be empty.

![Volume State - Detached](./assets/images/longhorn-volumn-state.png)

Create PV and PVC for the backup by selecting the backup and clicking "Create PV/PVC".

![Create PV/PVC - Step 1](./assets/images/longhorn-volumn-create-pv-pvc.png)

Click "OK".

![Create PV/PVC - Confirmation](./assets/images/longhorn-volumn-create-pv-pvc-2.png)

The PV/PVC section of the volume will then become "Bound".

![Volume State - Bound](./assets/images/longhorn-volumn-create-pv-pvc-3.png)

If you have multiple volumes to recover, please repeat the process for each one.

## Step 4 - Deploy Services

Switch to the correct Kubernetes context:

```bash
kubectl config use-context {{k8s-context}}
```

Ensure the service you wish to recover is not running. You can uninstall the service using Helm:

```bash
helm uninstall {{service-name}}
```

Redeploy the service, and the volume will be mounted automatically.

```bash
make manually-deploy-{{service-name}} DEPLOY_ENV={{env}}
```

![Service Restored Successfully](./assets/images/longhorn-restore-successully.png)
