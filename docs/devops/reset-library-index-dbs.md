# How to Reset the Library and Index Databases

This guide provides instructions on how to reset the library and index databases, which can be useful when wanting to clear the data for development or testing purposes in a deployed Kubernetescluster.

## Reset Library Service

### Port Forwards for Library Service

```bash
# Terminal 1
kubectl port-forward svc/library-mongo 27018:27017
# Terminal 2
kubectl port-forward svc/schemaparser-redis 6379:6379
```

### Reset Library MongoDB

#### Access Library MongoDB

![image](https://user-images.githubusercontent.com/11765228/126724352-21b79a56-baa0-430d-a92f-abcfd576bd9f.png)
![image](https://user-images.githubusercontent.com/11765228/126724287-968992ad-e0b9-43e5-9e90-18b2f22f79f4.png)

#### Delete Documents

![image](https://user-images.githubusercontent.com/11765228/126724549-33601a52-e731-454a-b03d-b69f28db6f54.png)

### Update Redis

#### Access Redis

![image](https://user-images.githubusercontent.com/11765228/126724409-2ed92781-9a74-4a19-a93b-dbf28501e947.png)

#### Update Last Commit Timestamp

![image](https://user-images.githubusercontent.com/11765228/126724486-13de6c7a-9859-45b5-b9ac-29465ba6f5f3.png)

## Reset Index Service

### Port Forwards for Index Service

```bash
# Terminal 1
kubectl port-forward svc/index-mongo 27017:27017
# Terminal 2
kubectl port-forward svc/index-kibana 5601:5601
```

### Reset Elasticsearch

#### Access Elasticsearch

<http://localhost:5601/app/dev_tools#/console>

#### Delete ES Documents

```js
POST nodes/_delete_by_query
{
  "query": { 
    "match_all": {}
  }
}
```

![image](https://user-images.githubusercontent.com/11765228/126725134-f2aed913-3149-4e90-90dc-f99a336b9941.png)

### Reset Index MongoDB

#### Access Index MongoDB

![image](https://user-images.githubusercontent.com/11765228/126724738-f5ddd133-a85c-438b-a1e0-f6296b268d3a.png)
![image](https://user-images.githubusercontent.com/11765228/126724745-fd9a625f-172c-4b9e-a8b5-32c79f24fe6a.png)

#### Delete MongoDB Documents

![Screen Shot 2021-07-23 at 8 35 37 AM](https://user-images.githubusercontent.com/11765228/126725249-e2c869c8-3451-4be3-9406-8d483151abc7.png)

```js
db.getCollection('nodes').deleteMany({})
```
