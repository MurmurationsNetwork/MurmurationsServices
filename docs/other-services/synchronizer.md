# Murmurations Synchronizer

Murmurations Synchronizer is a service designed to support the [Murmurations Map](https://github.com/MurmurationsNetwork/MurmurationsMap). It addresses the challenge posed by Elasticsearch, which limits data retrieval to 10,000 records per query. By enabling effective data synchronization from the index service to [`mapdata`](/docs/other-services/front-end-mongodb.md) profiles, the Synchronizer guarantees extensive and current data access beyond the Elasticsearch constraint.

## Features

1. Data Synchronization: Execute the `/export` route from the index service to synchronize and update profiles data within `mapdata` profiles.
2. `sort` Field: Utilizes a `sort` field within the settings collection to track and manage the order of data synchronization, guaranteeing that the most recent updates are always reflected.

## Dependencies

1. Vercel: Requires a Vercel account for deployment.
2. Kubernetes: Access to a Kubernetes cluster for scheduling cronjobs.
3. MongoDB: An existing MongoDB instance for data storage, currently hosted in our production and test Kubernetes clusters.

## Architecture

### Vercel

The API component of the Murmurations Synchronizer is deployed on Vercel, providing a serverless environment that supports dynamic scaling and high availability. This setup is ideal for handling API requests efficiently and securely.

### Kubernetes

In addition to the Vercel deployment, the project utilizes Kubernetes for scheduling and executing cronjobs every minute to ensure regular data updates.

## Deployment Guide

### Vercel Deployment

1. Prerequisites: Ensure you have a GitHub account and repository containing the Synchronizer service code.
2. Steps:
    - Connect your GitHub account to Vercel.
    - Select the repository for deployment.
    - Configure the API_SECRET_KEY in Vercel's environment variables section to secure API requests.
    - Click deploy to initiate the service hosting on Vercel.

### Kubernetes Deployment

1. Creating the Secret: Use the command `kubectl create secret generic synchronizer-job-secret --from-literal="API_SECRET_KEY=YOUR_KEY"` to store the API_SECRET_KEY securely in Kubernetes.
2. Deployment: Run `make deploy` to deploy the cronjob service in the Kubernetes environment, ensuring regular data updates.

## Monitoring and Error Handling

### Vercel Error Handling

For the Vercel-hosted components of the Murmurations Synchronizer, error monitoring primarily relies on manual inspection of Vercel logs. Currently, **automated alerting is not configured**. To identify any issues, it is necessary to periodically check the Vercel logs.

### Kubernetes Error Handling

In the Kubernetes environment, the monitoring service will highlight errors or failures in cronjob execution, enabling quicker identification and troubleshooting of issues encountered during data synchronization tasks.

## Migration Notice

1. Vercel Service Migration: If there's a need to migrate the API hosted on Vercel to another platform or service, the migration process will not impact the data stored in MongoDB. Since the data layer remains unaffected, the transition can be made smoothly with proper configuration updates to ensure the new hosting environment communicates correctly with MongoDB.
2. Kubernetes Service Migration: Similar to the Vercel service, migrating the Kubernetes cronjobs to another cluster or scheduling service does not need extra data migration. The key is to update the Kubernetes configurations to maintain secure and authenticated API calls to the Vercel-hosted service.
3. MongoDB: Unlike the migrations related to the Vercel or Kubernetes services, migrating MongoDB to a new environment or service requires `mongodump` to dump to new MongoDB service to ensure data integrity and continuity.
