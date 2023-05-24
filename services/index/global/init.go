package global

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

func Init() {
	config.Init()
	mongoInit()
	esInit()
	natsInit()
}

func mongoInit() {
	uri := mongo.GetURI(
		config.Conf.Mongo.USERNAME,
		config.Conf.Mongo.PASSWORD,
		config.Conf.Mongo.HOST,
	)

	err := mongo.NewClient(uri, config.Conf.Mongo.DBName)
	if err != nil {
		logger.Panic("Error when trying to connect to MongoDB", err)
	}

	err = mongo.Client.Ping()
	if err != nil {
		logger.Panic("Error when trying to ping the MongoDB", err)
	}
}

func esInit() {
	var indices = []elastic.Index{
		{
			Name: constant.ESIndex.Node,
			Body: `{
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
							"profile_url",
							"status",
							"tags",
							"primary_url"
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
			}`,
		},
	}

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
	err := nats.NewClient(
		config.Conf.Nats.ClusterID,
		config.Conf.Nats.ClientID,
		config.Conf.Nats.URL,
	)
	if err != nil {
		logger.Panic("Error when trying to connect nats", err)
	}
}
