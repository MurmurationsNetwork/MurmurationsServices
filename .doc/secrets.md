# Managing Secret using kubectl

## Create Secrets

### 1. Change the current-context

```
kubectl config use-context CONTEXT_NAME
```

*Note: You can also change contexts on Docker Desktop*

### 2. Run the following commands to create secrets

**Create secrets for MongoDB**

```bash
kubectl \
  create secret generic index-mongo-secret \
  --from-literal="MONGO_INITDB_ROOT_USERNAME=index-admin" \
  --from-literal="MONGO_INITDB_ROOT_PASSWORD=password"

kubectl \
  create secret generic library-mongo-secret \
  --from-literal="MONGO_INITDB_ROOT_USERNAME=library-admin" \
  --from-literal="MONGO_INITDB_ROOT_PASSWORD=password"
```

**Create secrets for services**

```bash
kubectl \
  create secret generic index-secret \
  --from-literal="MONGO_USERNAME=index-admin" \
  --from-literal="MONGO_PASSWORD=password"

kubectl \
  create secret generic library-secret \
  --from-literal="MONGO_USERNAME=library-admin" \
  --from-literal="MONGO_PASSWORD=password"

# schemaparser connects to Library MongoDB.
kubectl \
  create secret generic schemaparser-secret \
  --from-literal="MONGO_USERNAME=library-admin" \
  --from-literal="MONGO_PASSWORD=password"

# nodecleaner connects to Index MongoDB.
kubectl \
  create secret generic nodecleaner-secret \
  --from-literal="MONGO_USERNAME=index-admin" \
  --from-literal="MONGO_PASSWORD=password"

# revalidatenode connects to Index MongoDB.
kubectl \
  create secret generic revalidatenode-secret \
  --from-literal="MONGO_USERNAME=index-admin" \
  --from-literal="MONGO_PASSWORD=password"
```

**Create secrets for kibana (the logging interface)**

```
htpasswd -c ./auth <your-user>

kubectl -n kube-logging create secret generic kibana-basic-auth --from-file auth

rm auth
```

## Delete a Secret

To delete the Secret you have just created:

```
kubectl delete secret index-secret
```
