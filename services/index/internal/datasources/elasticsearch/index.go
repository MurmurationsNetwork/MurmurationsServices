package elasticsearch

import (
	"context"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/olivere/elastic"
)

type index struct {
	name constant.ESIndexType
	body string
}

var indices = []index{
	{
		name: constant.ESIndex().Node,
		body: `{
			"mappings":{
			   "dynamic":"false",
			   "_source":{
				  "includes":[
					 "geolocation",
					 "last_validated",
					 "linked_schemas",
					 "location",
					 "profile_url"
				  ]
			   },
			   "properties":{
				  "geolocation":{
					 "type":"geo_point"
				  },
				  "last_validated":{
					 "type":"date",
					 "format":"epoch_second"
				  },
				  "linked_schemas":{
					 "type":"keyword"
				  },
				  "location":{
					 "properties":{
						"country":{
						   "type":"keyword"
						},
						"locality":{
						   "type":"keyword"
						},
						"region":{
						   "type":"keyword"
						}
					 }
				  },
				  "profile_url":{
					 "type":"keyword"
				  }
			   }
			}
		 }`,
	},
}

func createMappings(client *elastic.Client) error {
	for _, index := range indices {
		exists, err := client.IndexExists(string(index.name)).Do(context.Background())
		if err != nil {
			return err
		}
		if !exists {
			createIndex, err := client.CreateIndex(string(index.name)).BodyString(index.body).Do(context.Background())
			if err != nil {
				return err
			}
			if !createIndex.Acknowledged {
				return err
			}
		}
	}
	return nil
}
