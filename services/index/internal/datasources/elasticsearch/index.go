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
					 "profileUrl",
					 "linkedSchemas",
					 "name",
					 "url",
					 "mission",
					 "keywords",
					 "geolocation",
					 "maplocation",
					 "lastChecked"
				  ]
			   },
			   "properties":{
				  "geolocation":{
					 "type":"geo_point"
				  },
				  "keywords":{
					 "type":"keyword"
				  },
				  "lastChecked":{
					 "type":"date",
					 "format":"epoch_second"
				  },
				  "linkedSchemas":{
					 "type":"keyword"
				  },
				  "maplocation":{
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
				  "mission":{
					 "type":"text"
				  },
				  "name":{
					 "type":"keyword"
				  },
				  "profileUrl":{
					 "type":"keyword"
				  },
				  "url":{
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
