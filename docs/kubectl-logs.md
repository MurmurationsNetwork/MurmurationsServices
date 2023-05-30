# Kubectl Logs

```
# Display only the most recent 20 lines of output in pod nginx
kubectl logs --tail=100 pod/validation-5f448b9f9c-jd7np

# Show all logs from pod nginx written in the last hour
kubectl logs --since=1h nginx
```
