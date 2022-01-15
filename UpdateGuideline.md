# Issue-184 Update Guideline
1. Connect to the cluster.
2. Get the index-mongo-name.
   ```
   kubectl get pods
   ```
3. Connect to index-mongo directly. (Replace the index-mongo-name with the name you get in previous step.)
   ```
   kubectl exec -it "index-mongo-name" -- /bin/bash
   ```
4. Execute the commands line by line to update index-mongo.
   ```
   mongo -u index-admin -p password
   use murmurationsIndex
   db.nodes.updateMany({}, {$rename: { "last_validated": "last_updated" }})
   exit
   exit
   ```