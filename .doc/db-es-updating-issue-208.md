# Updating MongoDB and Elasticsearch

This was necessitated by [issue #208](https://github.com/MurmurationsNetwork/MurmurationsServices/issues/208).

## Update ElasticSearch
1. The **query condition need to be changed** according to the condition. (Here uses wildcard as an example)
   ```
   POST nodes/_update_by_query?conflicts=proceed
    {
        "query": {
            "wildcard": {
            "profile_url": {
                "value": "*murmkvm1*"
            }
            }
        },
        "script": {
            "source": """
            def str = ctx._source['geolocation'];
            def first = str.indexOf(",");
            ctx._source.geolocation = ['lat': Float.parseFloat(str.substring(0, first)), 'lon': Float.parseFloat(str.substring(first + 1, str.length()))]
            """,
            "lang": "painless"
        }
    }
   ```