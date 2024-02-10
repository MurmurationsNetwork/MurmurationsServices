# Migrating MongoDB Data

## Introduction

This guide outlines a methodical approach for migrating MongoDB data between environments, highlighting the use of MongoDB Database Tools for effective and secure data transfer. By following the steps provided, you will export data from a source MongoDB database and import it into a destination MongoDB database. This process is designed to ensure data integrity and minimize downtime.

Upon completing this guide, you will have:

- Installed MongoDB Database Tools on your machine.
- Exported data from the source MongoDB database.
- Imported data into the destination MongoDB database.

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Step 1 - Installing MongoDB Database Tools](#step-1---installing-mongodb-database-tools)
- [Step 2 - Switching Kubernetes Context to Source](#step-2---switching-kubernetes-context-to-source)
- [Step 3 - Port Forwarding the Source MongoDB Service](#step-3---port-forwarding-the-source-mongodb-service)
- [Step 4 - Dumping Data from the Source Database](#step-4---dumping-data-from-the-source-database)
- [Step 5 - Switching Kubernetes Context to Destination](#step-5---switching-kubernetes-context-to-destination)
- [Step 6 - Port Forwarding the Destination MongoDB Service](#step-6---port-forwarding-the-destination-mongodb-service)
- [Step 7 - Restoring Data to the Destination Database](#step-7---restoring-data-to-the-destination-database)
- [Conclusion](#conclusion)

## Prerequisites

Before starting, ensure you have:

1. Access to both source and destination Kubernetes clusters.
2. Administrative credentials for the MongoDB instances in both source and destination environments.

## Step 1 - Installing MongoDB Database Tools

To facilitate the data migration, begin by installing the MongoDB Database Tools:

```bash
brew tap mongodb/brew
brew install mongodb-database-tools
```

## Step 2 - Switching Kubernetes Context to Source

Change to the Kubernetes context for the source environment:

```bash
kubectl config use-context <source-context-name>
```

Replace `<source-context-name>` with the actual context name of your source Kubernetes cluster.

## Step 3 - Port Forwarding the Source MongoDB Service

For local access to the source MongoDB service, open a new terminal tab and set up a port forwarding rule. This allows you to continue executing commands in another session without interruption:

```bash
kubectl port-forward svc/index-mongo 27017:27017
```

## Step 4 - Dumping Data from the Source Database

Open another terminal tab and export the data from the source MongoDB database. Ensure to replace `<password>` with your actual MongoDB user password. The exported data will be saved to `~/Desktop/index-mongodb-backups`:

```bash
mongodump --host localhost --port 27017 --username index-admin --password <password> --authenticationDatabase admin --db murmurationsIndex --collection nodes --out ~/Desktop/index-mongodb-backups
```

**Note:** Substitute `<password>` with the real password, and the data will be stored in `~/Desktop/index-mongodb-backups`.

## Step 5 - Switching Kubernetes Context to Destination

Switch to the Kubernetes context for the destination environment:

```bash
kubectl config use-context <destination-context-name>
```

Replace `<destination-context-name>` with the context name of your destination Kubernetes cluster.

## Step 6 - Port Forwarding the Destination MongoDB Service

In a new terminal tab, establish port forwarding to the destination MongoDB service:

```bash
kubectl port-forward svc/index-mongo 27017:27017
```

## Step 7 - Restoring Data to the Destination Database

With port forwarding in place, import the exported data into the destination MongoDB database. Remember to replace `<password>` with the actual password:

```bash
mongorestore --host localhost --port 27017 --username index-admin --password <password> --authenticationDatabase admin --db murmurationsIndex --collection nodes ~/Desktop/index-mongodb-backups/murmurationsIndex/nodes.bson
```

**Note:** Substitute `<password>` with the real password.

## Conclusion

You have successfully migrated MongoDB data between environments, ensuring a secure transfer and maintaining data integrity.
