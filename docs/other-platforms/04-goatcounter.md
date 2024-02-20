# GoatCounter

GoatCounter is an open-source web analytics platform designed to provide insights into your website's traffic without compromising on privacy. This guide covers how to deploy GoatCounter on a Kubernetes cluster and how to migrate your data from an old service to a new one.

## Kubernetes Deployment

1. Create a Kubernetes secret named goatcounter-secret.

    ```bash
    kubectl create secret generic goatcounter-secret \
    --from-literal=GOATCOUNTER_DOMAIN="example.com" \
    --from-literal=GOATCOUNTER_EMAIL="user@example.com" \
    --from-literal=GOATCOUNTER_PASSWORD="securepassword"
    ```

2. Deployment: Run `make deploy` to deploy the deployment in the Kubernetes environment.

## Migration

### Backup from Old Service

1. Identify the Pod running GoatCounter in the old service.

    ```bash
    kubectl get pods | grep 'goatcounter'
    ```

2. Backup the goatcounter.db file from the identified Pod. Replace <POD_NAME> with the name of GoatCounter Pod.

    ```bash
    kubectl cp <POD_NAME>:/var/lib/sqlite/goatcounter.db ./goatcounter.db
    ```

### Restore to New Service

1. Identify the Pod running GoatCounter in the old service.

    ```bash
    kubectl get pods | grep 'goatcounter'
    ```

2. Restore the goatcounter.db file to the new service's Pod. Replace <POD_NAME> with the name of GoatCounter Pod.

    ```bash
    kubectl cp ./goatcounter.db <POD_NAME>:/var/lib/sqlite/goatcounter.db
    ```
