package global

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

func Init() {
	config.Init()
	esInit()
	natsInit()
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
		logger.Panic("Error when trying to ping Elasticsearch", err)
		return
	}
	err = elastic.Client.CreateMappings(indices)
	if err != nil {
		logger.Panic("Error when trying to create index for Elasticsearch", err)
		return
	}

	logger.Info("Elasticsearch index created")
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
