# Resolving MongoDB Replica Set Issues After RKE2 Upgrade

## Introduction

Upgrading RKE2 clusters may cause MongoDB replica sets to go offline or enter an invalid state. This document describes a step-by-step process to identify and resolve common MongoDB replica set issues after a cluster upgrade, ensuring minimal disruption and quick recovery of your database services.

## Table of Contents

- [Introduction](#introduction)
- [Symptoms](#symptoms)
- [Step 1 - Extracting the Replica Set Configuration](#step-1---extracting-the-replica-set-configuration)
- [Step 3 - Forcing Reconfiguration With the Primary Node](#step-3---forcing-reconfiguration-with-the-primary-node)
- [Step 4 - Re-adding the Remaining Members](#step-4---re-adding-the-remaining-members)
- [Example Commands for Common Deployments](#example-commands-for-common-deployments)
- [Conclusion](#conclusion)

## Symptoms

After upgrading your RKE2 cluster, you may notice that the MongoDB cluster is offline. When checking the replica set status by:

```bash
kubectl exec -it index-mongo-0 -- mongosh --eval "rs.status()"
```

you may receive an error like:

```bash
MongoServerError: Our replica set config is invalid or we are not a member of it
```

## Step 1 - Extracting the Replica Set Configuration

On the primary node (e.g., `index-mongo-0`), extract the current replica set configuration:

```bash
const local = db.getSiblingDB("local");
const cfg = local.system.replset.findOne();
cfg
```

The output will look similar to:

```json
{
  "_id": "rs0",
  "version": 254031,
  "term": 15,
  "members": [
    { "_id": 0, "host": "index-mongo-0.index-mongo:27017", ... },
    { "_id": 1, "host": "index-mongo-1.index-mongo:27017", ... },
    { "_id": 2, "host": "index-mongo-2.index-mongo:27017", ... }
  ],
  ...
}
```

## Step 3 - Forcing Reconfiguration With the Primary Node

1. Identify the primary node: Typically, the pod with `-0` in its name (e.g., `index-mongo-0.index-mongo:27017`) is the primary.
2. Connect to that primary pod (if you're not already connected):

    ```bash
    kubectl exec -it index-mongo-0 -- mongosh
    ```

3. Update the replica set configuration to only include the primary node:

    ```bash
    cfg.version = (cfg.version || 1) + 1;
    cfg.members = [
      { _id: 0, host: "index-mongo-0.index-mongo:27017" }
    ];
    rs.reconfig(cfg, { force: true });
    ```

    > **Note:** Adjust the `host` and `_id` fields according to the actual primary node if it differs.

## Step 4 - Re-adding the Remaining Members

Once the replica set is successfully reconfigured with the primary node, re-add the remaining members:

```bash
rs.add("index-mongo-1.index-mongo:27017")
rs.add("index-mongo-2.index-mongo:27017")
```

Wait for the replica set to reach a healthy state, then check the status again:

```bash
rs.status()
```

## Example Commands for Common Deployments

The following commands can be used for other MongoDB StatefulSets in your cluster. Replace `index-mongo` with `library-mongo` or `data-proxy-mongo` as needed.

1. For `index-mongo` (Default Example)

    ```bash
    const local = db.getSiblingDB("local");
    const cfg = local.system.replset.findOne();
    cfg

    cfg.version = (cfg.version || 1) + 1;
    cfg.members = [ { _id: 0, host: "index-mongo-0.index-mongo:27017" } ];
    rs.reconfig(cfg, { force: true });

    rs.add("index-mongo-1.index-mongo:27017")
    rs.add("index-mongo-2.index-mongo:27017")
    ```

2. For `library-mongo`

    ```bash
    const local = db.getSiblingDB("local");
    const cfg = local.system.replset.findOne();
    cfg

    cfg.version = (cfg.version || 1) + 1;
    cfg.members = [ { _id: 0, host: "library-mongo-0.library-mongo:27017" } ];
    rs.reconfig(cfg, { force: true });

    rs.add("library-mongo-1.library-mongo:27017")
    rs.add("library-mongo-2.library-mongo:27017")
    ```

3. For `data-proxy-mongo`

    ```bash
    const local = db.getSiblingDB("local");
    const cfg = local.system.replset.findOne();
    cfg

    cfg.version = (cfg.version || 1) + 1;
    cfg.members = [ { _id: 0, host: "data-proxy-mongo-0.data-proxy-mongo:27017" } ];
    rs.reconfig(cfg, { force: true });

    rs.add("data-proxy-mongo-1.data-proxy-mongo:27017")
    rs.add("data-proxy-mongo-2.data-proxy-mongo:27017")
    ```

> In most cases, the `-0` pod is the primary node.  
> If your primary is different, **replace the host and `_id` values accordingly**.

## Conclusion

By following these steps, you can quickly recover a MongoDB replica set that has entered an invalid or offline state after an RKE2 upgrade. The process involves forcing a reconfiguration from the primary node and incrementally re-adding other members to restore normal operations. For different MongoDB deployments, simply adjust the pod names in the commands above.

Go back to [Home](../README.md).
