# Updating MongoDB and Elasticsearch

This was necessitated by [issue #184](https://github.com/MurmurationsNetwork/MurmurationsServices/issues/184).

## Update ElasticSearch

Copy the commands into ElasticSearch Dev Tool and execute one by one. 

```
# Create a temporary index called nodes2 with the new mappings

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
            "last_updated": {
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

# Link the nodes index to nodes2 and replace `last_validated` with `last_updated`

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

# Delete the old nodes index

DELETE /nodes

# Create nodes with the new mappings

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
            "last_updated": {
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

# Copy the records from nodes2 to nodes

POST _reindex
{
   "source": {
      "index": "nodes2"
   },
   "dest": {
      "index": "nodes"
   }
}

# Delete the temporary nodes2 index

DELETE /nodes2
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
