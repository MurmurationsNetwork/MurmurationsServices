package global

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

// Init is used to initialize essential services.
func Init() {
	esInit()
	natsInit()
}

// esInit initializes Elasticsearch service and sets up necessary indices.
func esInit() {
	// Define indices for Elasticsearch
	var indices = []elastic.Index{
		{
			Name: constant.ESIndex.Node,
			Body: `{
				"mappings": {
					"dynamic": "false",
					"_source": {
						"includes": [
							"name",
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
						"name": {
							"type": "text"
						},
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

	// Initialize a new Elasticsearch client.
	err := elastic.NewClient(config.Values.ES.URL)
	if err != nil {
		logger.Panic("Failed to create Elasticsearch client", err)
		return
	}

	// Create indices in Elasticsearch.
	err = elastic.Client.CreateMappings(indices)
	if err != nil {
		logger.Panic("Failed to create index mappings for Elasticsearch", err)
		return
	}

	logger.Info("Elasticsearch index created successfully")
}

// natsInit initializes Nats service.
func natsInit() {
	// Initialize a new Nats client.
	err := nats.NewClient(
		config.Values.Nats.ClusterID,
		config.Values.Nats.ClientID,
		config.Values.Nats.URL,
	)
	if err != nil {
		logger.Panic("Failed to create Nats client", err)
	}
}
