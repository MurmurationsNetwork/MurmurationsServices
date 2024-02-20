# Data Deletion Process in ElasticSearch and MongoDB

## Step 1: Define the Query in ElasticSearch

Before proceeding with data deletion, we need to find the data which we want to delete in the search queries. Here are two scenarios for defining the search query based on linked_schemas:

### Scenario 1: Deleting Nodes Containing a Specific Schema

To delete nodes that contain a specific schema (e.g., test_schema-v2.1.0), use the following search query. This query identifies all nodes linked to the specified schema.

```json
GET /nodes/_search
{
  "query": {
    "match": {
      "linked_schemas": "test_schema-v2.1.0"
    }
  },
  "_source": ["profile_url"]
}
```

### Scenario 2: Deleting Nodes Containing Only One Specific Schema

If the requirement is to delete nodes that contain only one schema, and that schema matches the specified one (e.g., test_schema-v2.1.0), use the below search query. This query filters the nodes to ensure that only those with exactly one linked schema matching the query are showed.

```json
GET /nodes/_search
{
  "query": {
    "bool": {
      "must": [
        {
          "match": {
            "linked_schemas": "test_schema-v2.1.0"
          }
        }
      ],
      "filter": [
        {
          "script": {
            "script": {
              "source": "doc['linked_schemas'].size() == 1",
              "lang": "painless"
            }
          }
        }
      ]
    }
  },
  "_source": ["profile_url"]
}
```

## Step 2: Processing the Search Results

1. After executing the search query in Elasticsearch, you will receive a JSON response containing the hits section. This section lists the nodes that match the search criteria.
    Note: **The number of hits (hits.total.value) as this number should match** in the following deletion process in both ElasticSearch and MongoDB.
    Example Result:

    ```json
    "hits" : {
        "total" : {
          "value" : 2,
          "relation" : "eq"
        },
        "hits" : [
          {
            "_id" : "57e488e7...",
            "_source" : {
              "profile_url" : "https://ic3.dev/murmurations/e2e-tests/profile-2.json"
            }
          },
          {
            "_id" : "7ccc4ca1...",
            "_source" : {
              "profile_url" : "https://ic3.dev/murmurations/e2e-tests/profile-10.json"
            }
          }
        ]
    }
    ```

2. To extract the `profile_url` from the search results, create a file named transform.js. Replace the placeholder hits array in the script below with the actual hits received from your Elasticsearch query.

    ```json
    const hits = [
      // Replace this with the actual hits from the Elasticsearch response
      {
        "_id" : "57e488e7...",
        "_source" : {
          "profile_url" : "https://ic3.dev/murmurations/e2e-tests/profile-2.json"
        }
      },
      {
        "_id" : "7ccc4ca1...",
        "_source" : {
          "profile_url" : "https://ic3.dev/murmurations/e2e-tests/profile-10.json"
        }
      }
    ]

    const profileUrls = hits.map(hit => hit._source.profile_url);
    console.log(profileUrls);
    ```

3. Execute the script.

    ```bash
    node transform.js
    ```

4. You will get the following result using the above example:

    ```json
    [
      'https://ic3.dev/murmurations/e2e-tests/profile-2.json',
      'https://ic3.dev/murmurations/e2e-tests/profile-10.json'
    ]
    ```

## Step 3: Generating MongoDB Commands for Data Deletion

1. The blank deletion command is as the following.

    ```javascript
    db.nodes.deleteMany({
      "profile_url": {
        "$in": []
      }
    });
    ```

2. Replace the placeholder array in the $in operator with the array of profile_urls obtained from the step 2. Using the above example, the command should be as the following:

    ```javascript
    db.nodes.deleteMany({
      "profile_url": {
        "$in": [
        'https://ic3.dev/murmurations/e2e-tests/profile-2.json',
        'https://ic3.dev/murmurations/e2e-tests/profile-10.json'
        ]
      }
    });
    ```

3. Execute the above command in MongoDB, the deletedCount should as same as `hits.total.value` in step 2.

## Step 4: Deleting Data from ElasticSearch

Replace <your_query> with the actual query used in Step 1 to match the documents you intend to delete.

```json
POST /nodes/_delete_by_query
{
  "query": <your_query>
}
```

For example:

```json
POST /nodes/_delete_by_query
{
  "query": {
    "match": {
      "linked_schemas": "test_schema-v2.1.0"
    }
  }
}
```

Note: After executing the delete command, ElasticSearch will return a response with a field name called "deleted". Verify that this number matches the `hits.total.value` from the ElasticSearch query results in Step 2.
