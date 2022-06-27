# Updating Elasticsearch
1. Execute in ES dev tools.
   ```
   PUT nodes/_settings
   {
       "max_result_window" : 500000
   }
   ```
2. [Reference](https://www.elastic.co/guide/en/elasticsearch/guide/current/pagination.html)
