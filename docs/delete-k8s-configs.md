# How to delete Kubernetes configs

## List all the config

```
kubectl config get-contexts
```

## Delete a config

```
kubectl config unset contexts.<context_name>
```
