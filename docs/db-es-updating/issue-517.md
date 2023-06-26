## Update ElasticSearch

Copy the commands into ElasticSearch Dev Tool and execute one by one.

```
PUT /nodes2
{
  "mappings": {
      "dynamic": "false",
      "_source": {
          "includes": [
              "name",
              "geolocation",
              "last_updated",
              "linked_schemas",
              "country",
              "locality",
              "region",
              "profile_url",
              "status",
              "tags",
              "primary_url"
          ]
      },
      "properties": {
          "name": {
              "type": "text"
          },
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
              "type": "text"
          },
          "locality": {
              "type": "text"
          },
          "region": {
              "type": "text"
          },
          "profile_url": {
              "type": "keyword"
          },
          "status": {
              "type": "keyword"
          },
          "tags": {
              "type": "text"
          },
          "primary_url": {
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

DELETE /nodes

PUT /nodes
{
  "mappings": {
      "dynamic": "false",
      "_source": {
          "includes": [
              "name",
              "geolocation",
              "last_updated",
              "linked_schemas",
              "country",
              "locality",
              "region",
              "profile_url",
              "status",
              "tags",
              "primary_url"
          ]
      },
      "properties": {
          "name": {
              "type": "text"
          },
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
              "type": "text"
          },
          "locality": {
              "type": "text"
          },
          "region": {
              "type": "text"
          },
          "profile_url": {
              "type": "keyword"
          },
          "status": {
              "type": "keyword"
          },
          "tags": {
              "type": "text"
          },
          "primary_url": {
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

DELETE /nodes2
```

## Update all node's Name
### localhost
1. Connect to Index pod. 
```bash
cd services/index/cmd/esseeder
go run main.go
```
### Production
1. Connect to Index pod. 
```bash
cd app
./seeder
```
