# Simulating a Service Being Offline

Sometimes it is useful to simulate a service being offline to test the fallback mechanisms in your code.

To simulate a service being offline, you can use the `kubectl scale deployments` command.

```bash
kubectl scale deployments index-app --replicas=0
kubectl scale deployments library-app --replicas=0
kubectl scale deployments dataproxy-app --replicas=0
```

To scale the services back to their original number of replicas, you can use the following commands:

```bash
kubectl scale deployments index-app --replicas=1
kubectl scale deployments library-app --replicas=1
kubectl scale deployments dataproxy-app --replicas=1
```

All the data will remain in the database, so you can continue to test the fallback mechanisms.
