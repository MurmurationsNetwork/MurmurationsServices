# Murmurations Allocator

The Murmurations Allocator is A service designed to support the [Murmurations Map](https://github.com/MurmurationsNetwork/MurmurationsMap). Through the Allocator, the Map service can retrieve data via the `/profiles` route in the format of `[longitude, latitude, profile_url]`.

## Dependencies

1. Kubernetes: Access to a Kubernetes cluster for managing the service's deployment.
2. MongoDB: An existing MongoDB instance for data storage, currently hosted in our production and test Kubernetes clusters.
3. Domain Name: A registered domain name for the service endpoint.

## Deployment Guide

### Kubernetes Deployment

1. Create a Kubernetes secret named allocator-app-secret with MongoDB credentials, Replace `mongo-url`, `mongo-admin` and `mongo-password` with actual credentials to ensure secure database access.

    ```bash
    kubectl create secret generic allocator-app-secret \
        --from-literal="MONGO_HOST=mongodb+srv://mongo-url" \
        --from-literal="MONGO_USERNAME=mongo-admin" \
        --from-literal="MONGO_PASSWORD=mongo-password"
    ```

2. Setup Domain name: A dedicated domain, for example: `allocator.<DOMAIN_NAME>`. This domain will serve as the access point for the service.
3. Deployment: Run `make deploy-allocator` to deploy the deployment in the Kubernetes environment.

## Monitoring and Error Handling

### Kubernetes Error Handling

In the Kubernetes environment, the monitoring service will highlight errors or failures in the deployment, enabling quicker identification and troubleshooting of issues encountered during getting data from MongoDB.

## Migration Notice

Kubernetes Service Migration: Migrating this service does not need extra data migration, because the data is hosted in a locally-managed MongoDB instance in the Kubernetes cluster.
