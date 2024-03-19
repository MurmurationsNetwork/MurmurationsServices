# Updating MongoDB and Elasticsearch

This was necessitated by [issue #210](https://github.com/MurmurationsNetwork/MurmurationsServices/issues/210).

## Close/open K8s ingress

1. Change makefile
    ```
    ENV ?= live-test
    ```
2. change ingress.yaml
   - test-index.murmurations.network -> test-index1.murmurations.network
   - test-library.murmurations.network -> test-library1.murmurations.network
3. Make manually-deploy-ingress
4. Use Postman make sure service is done
5. Upgrade Elasticsearch
6. Deploy new branch
7. Change ingress.yaml
   - test-index1.murmurations.network -> test-index.murmurations.network
   - test-library1.murmurations.network -> test-library.murmurations.network
8. Make manually-deploy-ingress
9. Restore makefile
10. kubectl config delete-context do-lon1-murmtest

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
                "country",
                "locality",
                "region",
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
            "country": {
                "type": "keyword"
            },
            "locality": {
                "type": "keyword"
            },
            "region": {
                "type": "keyword"
            },
            "profile_url": {
                "type": "keyword"
            }
        }
    }
}

PUT _ingest/pipeline/move-location
{
   "description" : "move location",
   "processors" : [
      {
         "rename": {
         "field": "location.country",
         "target_field": "country"
         }
      },
      {
         "rename": {
         "field": "location.locality",
         "target_field": "locality"
         }
      },
      {
         "rename": {
         "field": "location.region",
         "target_field": "region"
         }
      }
   ]
}

POST _reindex
{
   "source": {
      "index": "nodes",
      "query": {
        "exists": {
          "field": "location"
        }
      }
   },
   "dest": {
      "index": "nodes2",
      "pipeline": "move-location"
   }
}

POST _reindex
{
   "source": {
      "index": "nodes",
      "query": {
        "bool": { 
          "must_not": {
            "exists": {
              "field": "location"
            }
          }
        }
      }
    },
   "dest": {
      "index": "nodes2"
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
                "country",
                "locality",
                "region",
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
            "country": {
                "type": "keyword"
            },
            "locality": {
                "type": "keyword"
            },
            "region": {
                "type": "keyword"
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
