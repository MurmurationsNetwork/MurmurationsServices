# Restoring Rancher

Source: [Restoring Rancher](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/backup-restore-and-disaster-recovery/restore-rancher)

## 1. Create the Restore Custom Resource

1. In the upper left corner, click **☰ > Cluster Management**.

2. On the **Clusters** page, go to the `local` cluster and click **Explore**. The `local` cluster runs the Rancher server.

3. In the left navigation bar, click **Rancher Backups > Restores**.

4. Click **Create**.

5. For using the YAML editor, we can click **Create > Create from YAML**. Enter the Restore YAML.

   ```yaml
   apiVersion: resources.cattle.io/v1
   kind: Restore
   metadata:
     name: restore-migration
   spec:
     backupFilename: [backup file name]
     storageLocation:
       s3:
         credentialSecretName: contabo-s3-creds
         credentialSecretNamespace: default
         bucketName: rancher-backup
         region: EU
         endpoint: eu2.contabostorage.com
   ```

6. Click **Create**.

## 2. Logs[​](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/backup-restore-and-disaster-recovery/restore-rancher#logs "Direct link to Logs")

To check how the restore is progressing, you can check the logs of the operator. Run this command to follow the logs:

```bash
kubectl logs -n cattle-resources-system -l app.kubernetes.io/name=rancher-backup -f
```
