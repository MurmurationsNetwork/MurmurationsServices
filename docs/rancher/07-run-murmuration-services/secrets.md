# Managing Secrets with kubectl

This guide will walk you through the process of managing secrets using `kubectl`. Secrets are sensitive data, such as passwords, OAuth tokens, and ssh keys, which are stored and managed in Kubernetes.

## Creating Secrets

Remember to replace the `{{*_PASSWORD}}` placeholders below with strong passwords. Also make sure you use the same strong password for usernames that are the same (e.g., all instances of `INDEX_ADMIN_PASSWORD` below should be using the same strong password).

You will also need to create a GitHub access token for the `schemaparser` service ([see the next section](#creating-a-github-personal-access-token) for instructions).

Execute these commands to create secrets for each service:

```bash
kubectl \
  create secret generic index-mongo-secret \
  --from-literal="MONGO_INITDB_ROOT_USERNAME=index-admin" \
  --from-literal="MONGO_INITDB_ROOT_PASSWORD={{INDEX_ADMIN_PASSWORD}}"

kubectl \
  create secret generic library-mongo-secret \
  --from-literal="MONGO_INITDB_ROOT_USERNAME=library-admin" \
  --from-literal="MONGO_INITDB_ROOT_PASSWORD={{LIBRARY_ADMIN_PASSWORD}}"

kubectl \
  create secret generic data-proxy-mongo-secret \
  --from-literal="MONGO_INITDB_ROOT_USERNAME=data-proxy-admin" \
  --from-literal="MONGO_INITDB_ROOT_PASSWORD={{DATA_PROXY_ADMIN_PASSWORD}}"

kubectl \
  create secret generic index-secret \
  --from-literal="MONGO_USERNAME=index-admin" \
  --from-literal="MONGO_PASSWORD={{INDEX_ADMIN_PASSWORD}}"

kubectl \
  create secret generic library-secret \
  --from-literal="MONGO_USERNAME=library-admin" \
  --from-literal="MONGO_PASSWORD={{LIBRARY_ADMIN_PASSWORD}}"

kubectl \
  create secret generic data-proxy-secret \
  --from-literal="MONGO_USERNAME=data-proxy-admin" \
  --from-literal="MONGO_PASSWORD={{DATA_PROXY_ADMIN_PASSWORD}}"

kubectl \
  create secret generic nodecleaner-secret \
  --from-literal="MONGO_USERNAME=index-admin" \
  --from-literal="MONGO_PASSWORD={{INDEX_ADMIN_PASSWORD}}"

kubectl \
  create secret generic revalidatenode-secret \
  --from-literal="MONGO_USERNAME=index-admin" \
  --from-literal="MONGO_PASSWORD={{INDEX_ADMIN_PASSWORD}}"

kubectl \
  create secret generic dataproxyupdater-secret \
  --from-literal="MONGO_USERNAME=data-proxy-admin" \
  --from-literal="MONGO_PASSWORD={{DATA_PROXY_ADMIN_PASSWORD}}"

kubectl \
  create secret generic dataproxyrefresher-secret \
  --from-literal="MONGO_USERNAME=data-proxy-admin" \
  --from-literal="MONGO_PASSWORD={{DATA_PROXY_ADMIN_PASSWORD}}"

kubectl \
  create secret generic schemaparser-secret \
  --from-literal="MONGO_USERNAME=library-admin" \
  --from-literal="MONGO_PASSWORD={{LIBRARY_ADMIN_PASSWORD}}" \
  --from-literal="GITHUB_TOKEN={{GITHUB_TOKEN}}"
```

## Creating a GitHub Personal Access Token

Creating a secret for `schemaparser` requires a GitHub token. Please refer to [GitHub's documentation](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic) for creating a personal access token. No additional scopes are needed.

![Personal Access Token](./assets/images/personal-access-token.png)
