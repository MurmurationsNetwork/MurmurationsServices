# Back Up and Restore a Kubernetes Cluster on DigitalOcean Using Velero

## 1. Install Velero

```
brew install velero
```

## 2. Create Space on Digital Ocean

Navigate to spaces from left panel

![image](https://user-images.githubusercontent.com/11765228/114300594-ade30980-9af3-11eb-862b-196381951d8b.png)

## 3. Create Space Keys on Digital Ocean

Name can be anything of your choice

![image](https://user-images.githubusercontent.com/11765228/114300197-a458a200-9af1-11eb-884a-7ef32d60d3f8.png)

Copy the Key and paste in cloud-credentials file over `<SPACE_API_KEY>`

Copy the Secret and paste in cloud-credentials file over `<SPACE_API_SECRET>`

## 4. Create Personal Access token

Navigate to API section from left panel

![image](https://user-images.githubusercontent.com/11765228/114300326-5a23f080-9af2-11eb-87cd-c44f8c283831.png)

Generate the token

Copy the token and paste in `velero-secret.patch.yaml` over `<DIGITALOCEAN_TOKEN>`

## 5. Installing the Velero Server

```
velero install \
  --provider velero.io/aws \
  --bucket <SPACE_NAME> \
  --plugins velero/velero-plugin-for-aws:v1.0.0,digitalocean/velero-plugin:v1.0.0 \
  --backup-location-config s3Url=https://<REGION>.digitaloceanspaces.com,region=<REGION> \
  --use-volume-snapshots=false \
  --secret-file .doc/backup/cloud-credentials
```

Check for logs in Velero

```
kubectl logs deployment/velero -n velero
kubectl get deployment/velero --namespace velero
```

## 6. Configuring snapshots

```
velero snapshot-location create default --provider digitalocean.com/velero
```

## 7. Adding API Token

Execute the following command for updating the secret in velero namespace

```
kubectl patch secret/cloud-credentials -p "$(cat .doc/backup/velero-secret.patch.yaml)" --namespace velero
```

Execute the following command for updating the deployment to use the secret

```
kubectl patch deployment/velero -p "$(cat .doc/backup/velero-deployment.patch.yaml)" --namespace velero
```

## 8. Setup backup

```
velero create schedule index-mongo-backup \
  --schedule="@every 24h" \
  --include-namespaces default \
  --include-resources persistentvolumeclaims,persistentvolumes \
  --ttl 168h0m0s \
  --selector app=index-mongo

velero create schedule index-es-backup \
  --schedule="@every 24h" \
  --include-namespaces default \
  --include-resources persistentvolumeclaims,persistentvolumes \
  --ttl 168h0m0s \
  --selector app=index-es

velero create schedule library-mongo-backup \
  --schedule="@every 24h" \
  --include-namespaces default \
  --include-resources persistentvolumeclaims,persistentvolumes \
  --ttl 168h0m0s \
  --selector app=library-mongo
  
velero create schedule data-proxy-mongo-backup \
  --schedule="@every 24h" \
  --include-namespaces default \
  --include-resources persistentvolumeclaims,persistentvolumes \
  --ttl 168h0m0s \
  --selector app=data-proxy-mongo
```

# Creating backup using Velero

**Apply a label to the PV object**

```
kubectl get pvc && kubectl get pv
kubectl label pvc <PVC_NAME> app=<LABEL_NAME>
kubectl label pv <PV_NAME> app=<LABEL_NAME>
```

example

```
kubectl get pvc && kubectl get pv
kubectl label pvc index-mongo-pvc app=index-mongo
kubectl label pv pvc-36ab0171-6e7b-406d-a6ca-de038781d24f app=index-mongo
```

**On demand backup**

```
velero backup create <BACKUP_NAME> \
  --include-namespaces default \
  --include-resources persistentvolumeclaims,persistentvolumes \
  --ttl 168h0m0s \
  --selector app=<LABEL_NAME>
```

**Schedule backup (Cronjob)**

```
velero create schedule <BACKUP_NAME> \
  --schedule="@every 24h" \
  --include-namespaces default \
  --include-resources persistentvolumeclaims,persistentvolumes \
  --ttl 168h0m0s \
  --selector app=<LABEL_NAME>
```

**After creating backup check in velero for status**

```
velero backup describe <BACKUP_NAME>
```

**Get the list of backups**

```
velero get backups
```

**Delete a backup**

```
velero delete backup <BACKUP_NAME>
```

**Restore from backup**

```
velero restore create <RESTORE_NAME> --from-backup <BACKUP_NAME>
```

**Describe restores**

```
velero restore describe <RESTORE_NAME>
```

**Retrieve restores**

```
velero restore get
```

**Delete a restore**

```
velero restore delete index-mongo-backup
```
