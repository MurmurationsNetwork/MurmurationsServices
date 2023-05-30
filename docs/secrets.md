# Managing Secrets with kubectl

This guide will walk you through the process of managing secrets using `kubectl`. Secrets are sensitive data, such as passwords, OAuth tokens, and ssh keys, which are stored and managed in Kubernetes.

## Creating Secrets

Follow the steps below to create secrets:

### 1. Switch the current context

Use the following command to change the current context:

```bash
kubectl config use-context CONTEXT_NAME
```

_Note: You can also switch contexts using Docker Desktop._

### 2. Execute the commands below to create secrets

**Creating Secrets for MongoDB**

Run these commands to create secrets for MongoDB:

```bash
kubectl \
  create secret generic index-mongo-secret \
  --from-literal="MONGO_INITDB_ROOT_USERNAME=index-admin" \
  --from-literal="MONGO_INITDB_ROOT_PASSWORD=password"

kubectl \
  create secret generic library-mongo-secret \
  --from-literal="MONGO_INITDB_ROOT_USERNAME=library-admin" \
  --from-literal="MONGO_INITDB_ROOT_PASSWORD=password"

kubectl \
  create secret generic data-proxy-mongo-secret \
  --from-literal="MONGO_INITDB_ROOT_USERNAME=data-proxy-admin" \
  --from-literal="MONGO_INITDB_ROOT_PASSWORD=password"
```

**Creating Secrets for Each Service**

Execute these commands to create secrets for each service:

```bash
kubectl \
  create secret generic index-secret \
  --from-literal="MONGO_USERNAME=index-admin" \
  --from-literal="MONGO_PASSWORD=password"

kubectl \
  create secret generic library-secret \
  --from-literal="MONGO_USERNAME=library-admin" \
  --from-literal="MONGO_PASSWORD=password"

kubectl \
  create secret generic data-proxy-secret \
  --from-literal="MONGO_USERNAME=data-proxy-admin" \
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

# dataproxyupdater connects to DataProxy MongoDB.
kubectl \
  create secret generic dataproxyupdater-secret \
  --from-literal="MONGO_USERNAME=data-proxy-admin" \
  --from-literal="MONGO_PASSWORD=password"

# dataproxyrefresher connects to DataProxy MongoDB.
kubectl \
  create secret generic dataproxyrefresher-secret \
  --from-literal="MONGO_USERNAME=data-proxy-admin" \
  --from-literal="MONGO_PASSWORD=password"
```

Creating a secret for `schemaparser` requires a GitHub token. Follow this [instruction](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic) to create a personal access token. No additional scopes are needed for the token.

```bash
# schemaparser connects to Library MongoDB.
kubectl \
  create secret generic schemaparser-secret \
  --from-literal="MONGO_USERNAME=library-admin" \
  --from-literal="MONGO_PASSWORD=password" \
  --from-literal="GITHUB_TOKEN=<YOUR_GITHUB_ACCESS_TOKEN>"
```

## Deleting a Secret

To delete a secret that you've created, use the following command:

```bash
kubectl delete secret SECRET_NAME
```

Replace `SECRET_NAME` with the name of the secret you want to delete. For example, to delete `index-secret`, you would run:

```

```bash
kubectl delete secret index-secret
```