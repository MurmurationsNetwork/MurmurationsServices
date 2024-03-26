# Maintenance Mode

## Turn on maintenance mode

1. Deploy the Maintenance Service: Use the following command to deploy the maintenance service.

    ```bash
    make manually-deploy-maintenance
    ```

2. Update the Ingress Configuration: Modify the [ingress.yaml](./charts/murmurations/charts/ingress/templates/ingress/ingress.yaml) file to route traffic to the maintenance service. Change backend.service.name to maintenance-service:

   ```yaml
   backend:
     service:
       name: maintenance-service
   ```
   
3. Deploy the Updated Ingress: Apply the changes to the cluster with the following command, ensuring to replace <YOUR_DEPLOY_ENV> with your environment:

    ```bash
    make manually-deploy-ingress DEPLOY_ENV=<YOUR_DEPLOY_ENV>
    ```
   
## Turn off maintenance mode

1. Restore the Original Service in the Ingress Configuration: Revert the changes in [ingress.yaml](./charts/murmurations/charts/ingress/templates/ingress/ingress.yaml) by setting backend.service.name back to your original service (e.g., index-app):

   ```yaml
   backend:
     service:
       name: index-app
   ```
   
2. Deploy the Reverted Ingress: Apply the ingress changes to your environment with the following command, replacing <YOUR_DEPLOY_ENV> with your environment:

    ```bash
    make manually-deploy-ingress DEPLOY_ENV=<YOUR_DEPLOY_ENV>
    ```
   
3. Remove the Maintenance Service: Delete the temporary maintenance service from your cluster:

    ```bash
    helm delete murmurations-maintenance
    ```
