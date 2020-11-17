package elasticsearch

import "github.com/MurmurationsNetwork/MurmurationsServices/common/elastic"

var indices = []elastic.Index{
	{
		Name: "nodes",
		Body: `{
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
						   "type":"text"
						},
						"locality":{
						   "type":"text"
						},
						"region":{
						   "type":"text"
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
