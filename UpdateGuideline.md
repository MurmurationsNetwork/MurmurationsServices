# Issue-184 Update Guideline
## Update ElasticSearch
1. Copy the commands into ElasticSearch Dev Tool and execute one by one. The following code accomplish this functions.
   - Create a new node called nodes2 with the new mappings.
   - Paste the nodes indices to nodes2 and replace last_validated with last_updated.
   - Delete nodes.
   - Create nodes with the new mappings.
   - Copy indices from nodes2 to nodes.
   - Delete nodes2.
   ```
   PUT /nodes2
   {
      "mappings": {
         "dynamic": "false",
         "_source": {
               "includes": [
                  "geolocation",
                  "last_updated",
                  "linked_schemas",
                  "location",
                  "profile_url"
               ]
         },
         "properties": {
               "geolocation": {
                  "type": "geo_point"
               },
               "last_validated": {
                  "type": "date",
                  "format": "epoch_second"
               },
               "linked_schemas": {
                  "type": "keyword"
               },
               "location": {
                  "properties": {
                     "country": {
                           "type": "text"
                     },
                     "locality": {
                           "type": "text"
                     },
                     "region": {
                           "type": "text"
                     }
                  }
               },
               "profile_url": {
                  "type": "keyword"
               }
         }
      }
   }

   PUT _ingest/pipeline/rename-last_validated
   {
      "description" : "rename last_validated",
      "processors" : [
         {
            "rename": {
            "field": "last_validated",
            "target_field": "last_updated"
            }
         }
      ]
   }

   POST _reindex
   {
      "source": {
         "index": "nodes"
      },
      "dest": {
         "index": "nodes2",
         "pipeline": "rename-last_validated"
      }
   }

   DELETE /nodes

   PUT /nodes
   {
      "mappings": {
         "dynamic": "false",
         "_source": {
               "includes": [
                  "geolocation",
                  "last_updated",
                  "linked_schemas",
                  "location",
                  "profile_url"
               ]
         },
         "properties": {
               "geolocation": {
                  "type": "geo_point"
               },
               "last_validated": {
                  "type": "date",
                  "format": "epoch_second"
               },
               "linked_schemas": {
                  "type": "keyword"
               },
               "location": {
                  "properties": {
                     "country": {
                           "type": "text"
                     },
                     "locality": {
                           "type": "text"
                     },
                     "region": {
                           "type": "text"
                     }
                  }
               },
               "profile_url": {
                  "type": "keyword"
               }
         }
      }
   }

   POST _reindex
   {
      "source": {
         "index": "nodes2"
      },
      "dest": {
         "index": "nodes"
      }
   }
   ```
## Update MongoDB
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