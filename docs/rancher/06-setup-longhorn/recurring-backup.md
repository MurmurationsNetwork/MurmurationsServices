# Set Up Recurring Backup

## Introduction

This document provides a step-by-step guide on setting up recurring backups for your data using Contabo object storage. By following these steps, you will create a backup bucket in Contabo, generate access credentials, create a Kubernetes secret, and configure the backup settings in Longhorn.

## Table of Contents

- [Introduction](#introduction)
- [Step 1 - Creating a Bucket in Contabo Object Storage](#step-1---creating-a-bucket-in-contabo-object-storage)
- [Step 2 - Obtaining Access Key and Secret Key](#step-2---obtaining-access-key-and-secret-key)
- [Step 3 - Creating a Kubernetes Secret](#step-3---creating-a-kubernetes-secret)
- [Step 4 - Configuring Backup Settings in Longhorn](#step-4---configuring-backup-settings-in-longhorn)
- [Step 5 - Testing the Backup](#step-5---testing-the-backup)
- [Step 6 - Verifying the Backup](#step-6---verifying-the-backup)
- [Step 7 - Creating Recurring Backups](#step-7---creating-recurring-backups)
- [Conclusion](#conclusion)

## Step 1 - Creating a Space in DigitalOcean Object Storage

Navigate to the DigitalOcean Object Storage page to create a space.

![DigitalOcean Space Creation 1](./assets/images/do-space-creation-1.png)

Choose "Amsterdam" for the datacenter region, provide a unique space name, and then click "Create".

![DigitalOcean Space Creation 2](./assets/images/do-space-creation-2.png)

## Step 2 - Obtaining Access Key and Secret Key

Navigate to the "Applications & API" section on DigitalOcean and create a Spaces access key.

![Do Secret Keys](./assets/images/do-secret-keys.png)

## Step 3 - Creating a Kubernetes Secret

Open a terminal and switch your Kubernetes context to the desired one.

```bash
kubectl config use-context <context>
```

Replace `<your_access_key>` and `<your_secret_key>` with your actual credentials. Execute the following command to create the secret in the `longhorn-system` namespace:

```bash
kubectl create secret generic aws-secrets \
  --from-literal=AWS_ACCESS_KEY_ID=<your_access_key> \
  --from-literal=AWS_SECRET_ACCESS_KEY=<your_secret_key> \
  --from-literal=AWS_ENDPOINTS='https://ams3.digitaloceanspaces.com' \
  --from-literal=AWS_REGION=us-east-1 \
  --from-literal=VIRTUAL_HOSTED_STYLE=false \
  -n longhorn-system
```

## Step 4 - Configuring Backup Settings in Longhorn

Access the Longhorn user interface. Click on "Settings" in the top navigation bar. Within the "Backup" section, specify your Backup Target using the following format:

```bash
s3://[SPACENAME]@[REGION]/[FOLDER-NAME]
# For example: s3://mybackup@ams3/myfolder-backups
```

This configuration directs Longhorn to store backups in the specified S3-compatible storage location. Replace `[SPACENAME]` with the name of your space, `[REGION]` with the region of your storage, and `[FOLDER-NAME]` with your desired folder name for the backups.

For the Backup Target Credential Secret, use:

```bash
aws-secrets
```

Scroll down to the bottom and click on "Save".

![Longhorn Backup Settings](./assets/images/longhorn-backup-settings.png)

## Step 5 - Testing the Backup

Go to the "Volume" menu, select the volume you wish to back up, click on "Create Backup," and then click "OK."

![Longhorn Test Backup](./assets/images/longhorn-test-backup.png)

## Step 6 - Verifying the Backup

Ensure that your backup has been successfully created by navigating to the backup section and verifying the presence of your backup.

![Longhorn Check Backups](./assets/images/longhorn-check-backups.png)
![Longhorn Check Backups 2](./assets/images/longhorn-check-backups-2.png)

## Step 7 - Creating Recurring Backups

To set up recurring backups, go to the Recurring Job menu and click "Create Recurring Job".

![Longhorn Create Recurring Job](./assets/images/longhorn-create-recurring-job.png)

Change the task to "Backup", click "Add to Default Group", and then click "OK".

![Longhorn Recurring Job Config](./assets/images/longhorn-recurring-job-config.png)

## Conclusion

With the successful configuration of recurring backups using Contabo object storage and Longhorn, you have now ensured the safety and availability of your data.

Go back to [Home](../README.md).
