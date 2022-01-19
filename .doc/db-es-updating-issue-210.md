# Updating MongoDB and Elasticsearch

This was necessitated by [issue #210](https://github.com/MurmurationsNetwork/MurmurationsServices/issues/210).

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

POST _reindex
{
   "source": {
      "index": "nodes"
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
