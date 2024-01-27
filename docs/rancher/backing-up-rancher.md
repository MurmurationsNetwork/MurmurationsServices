# Backing up Rancher

Source: [Backing up Rancher | Rancher](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/backup-restore-and-disaster-recovery/back-up-rancher)

## 1. Create a Contabo Bucket

1. Visit <https://new.contabo.com/storage/object-storage/buckets>

2. Click **Create Bucket**

   1. Fill out the Bucket Name: rancher-backup

3. Visit <https://new.contabo.com/account/security>

4. You will see your bucket along with its access key and secret key.

## 2. Create Secrets in k3s

1. From the previous step, you have the access key and secret key.

2. Update the access key and secret key, and create a secret in the cluster running Rancher:

```sh
kubectl create secret generic contabo-s3-creds \
  --from-literal=accessKey=<access key> \
  --from-literal=secretKey=<secret key>
```

## 3. Install the Rancher Backup Operator

1. In the upper left corner, click **☰ > Cluster Management**.

2. On the **Clusters** page, go to the `local` cluster and click **Explore**. The `local` cluster runs the Rancher server.

3. Click **Apps > Charts**.

4. Click **Rancher Backups**.

5. Click **Install**.

6. Choose **Customize Helm options before install**.

7. Click **Next**.

8. In Chart options, choose **Use an S3-compatible object store**.

9. Fill out the following:

   1. Credential secret: contabo-s3-creds

   2. Bucket name: rancher-backup

   3. Region: EU

   4. Endpoint: [eu2.contabostorage.com](http://eu2.contabostorage.com)

10. Click **Next**.

11. Click **Install**.

## 4. Perform a Backup

To perform a backup, create a custom resource of type Backup.

1. In the upper left corner, click **☰ > Cluster Management**.

2. On the **Clusters** page, go to the `local` cluster and click **Explore**.

3. In the left navigation bar, click **Rancher Backups > Backups**.

4. Click **Create**.

   1. Name: daily-backup

   2. Click: Recurring Backups

   3. Schedule: 0 0 \* \* \*

   4. Retention count: 7

   5. Storage Location: Use the default
