# How to Debug Kubernetes Resources?

Use `kubectl describe` to see the error messages

```
kubetcl describe pod/<pod_name>
kubetcl describe service/<service_name>
kubetcl describe deployment.apps/<deployment_name>
kubetcl describe cronjob.batch/<cronjob_name>
```

Use `kubectl describe` to see the config values

```
k describe ConfigMap <config_name>
```
