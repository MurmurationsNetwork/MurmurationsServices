# GoatCounter

GoatCounter is an open-source web analytics platform designed to provide insights into your website's traffic without compromising on user privacy. We use it to track the number of visitors to the Murmurations websites.

This guide covers how to deploy GoatCounter on a Kubernetes cluster and how to migrate your data from an old service to a new one.

## Kubernetes Deployment

1. Create a Kubernetes secret named goatcounter-secret.

    ```bash
    kubectl create secret generic goatcounter-secret \
    --from-literal=GOATCOUNTER_DOMAIN="example.com" \
    --from-literal=GOATCOUNTER_EMAIL="user@example.com" \
    --from-literal=GOATCOUNTER_PASSWORD="securepassword" \
    --from-literal=PGDATABASE="goatcounter"
    --from-literal=PGHOST="goatcounter-pg-service"
    --from-literal=PGPORT="5432"
    --from-literal=PGUSER="goatcounter"
    --from-literal=PGPASSWORD="password"
    --from-literal=POSTGRES_DB="goatcounter"
    --from-literal=POSTGRES_USER="goatcounter"
    --from-literal=POSTGRES_PASSWORD="password"
    ```

2. Deployment: Run `make deploy` to deploy the deployment in the Kubernetes environment.

## Migration

### Backup from Old Service

1. Identify the Postgresql Pod running GoatCounter in the old service.

    ```bash
    kubectl get pods | grep 'goatcounter-pg'
    ```

2. Connect to the pod. Replace <POD_NAME> with the name of GoatCounter Pod.

    ```bash
    kubectl exec -it <POD_NAME> -- bash
    ```

3. Backup the goatcounter.db file.

   ```bash
   pg_dump -U goatcounter goatcounter > goatcounter.sql
   ```

4. Copy back sql from the old Pod.

   ```bash
   kubectl cp <POD_NAME>:/goatcounter.sql ~/Desktop/goatcounter.sql
   ```

### Restore to New Service

1. Identify the Pod running GoatCounter in the new service.

    ```bash
    kubectl get pods | grep 'goatcounter-pg'
    ```

2. Copy the goatcounter.db to new service's Pod. Replace <POD_NAME> with the name of GoatCounter Postgresql Pod.

    ```bash
    kubectl cp ~/Desktop/goatcounter.sql <POD_NAME>:/goatcounter.sql
    ```

3. Connect to the identified Pod and delete the goatcounter.db file.

    ```bash
    kubectl exec -it <POD_NAME> -- bash
    ```

4. Connect to PostgreSQL to execute an SQL command.

    ```bash
    psql -U goatcounter -d postgres
    ```

5. Delete the auto-generated database and create a new, empty database for backup restoration.

    ```sql
    DROP DATABASE "goatcounter";
    CREATE DATABASE "goatcounter";
    \q
    ```

6. Restore the goatcounter.sql file to default goatcounter database.

    ```bash
    psql -U goatcounter -d goatcounter -f /goatcounter.sql
    ```
