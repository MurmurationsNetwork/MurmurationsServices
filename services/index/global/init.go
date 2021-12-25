package global

import (
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

func Init() {
	config.Init()

	fmt.Println("config Init ok")

	mongoInit()

	fmt.Println("mongoInit ok")

	esInit()

	fmt.Println("esInit ok")

	natsInit()

	fmt.Println("natsInit ok")
}

func mongoInit() {
	uri := mongo.GetURI(config.Conf.Mongo.USERNAME, config.Conf.Mongo.PASSWORD, config.Conf.Mongo.HOST)

	fmt.Println("mongo uri", uri)

	err := mongo.NewClient(uri, config.Conf.Mongo.DBName)
	if err != nil {
		logger.Panic("Error when trying to connect to MongoDB", err)
	}

	fmt.Println("NewClient OK", uri)

	err = mongo.Client.Ping()
	if err != nil {
		logger.Panic("Error when trying to ping the MongoDB", err)
	}

	fmt.Println("Client Ping", uri)
}

func esInit() {
	var indices = []elastic.Index{
		{
			Name: constant.ESIndex.Node,
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

	fmt.Println("config.Conf.ES.URL ", config.Conf.ES.URL)

	err := elastic.NewClient(config.Conf.ES.URL)
	if err != nil {
		logger.Panic("Error when trying to ping the ElasticSearch", err)
		return
	}
	err = elastic.Client.CreateMappings(indices)
	if err != nil {
		logger.Panic("Error when trying to create index for ElasticSearch", err)
		return
	}

	logger.Info("ElasticSearch index created")
}

func natsInit() {
	err := nats.NewClient(config.Conf.Nats.ClusterID, config.Conf.Nats.ClientID, config.Conf.Nats.URL)
	if err != nil {
		logger.Panic("Error when trying to connect nats", err)
	}
}
