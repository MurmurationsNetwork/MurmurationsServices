package index

import (
	"log"
	"os"

	env "github.com/caarlos0/env/v10"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

func init() {
	err := env.Parse(&config.Values)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
	setupElasticsearch()
}

// setupElasticsearch initializes Elasticsearch service and sets up necessary indices.
func setupElasticsearch() {
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
							"primary_url",
							"expires"
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
						},
						"expires": {
							"type": "date",
							"format": "epoch_second"
						}
					}
				}
			}`,
		},
	}

	// Initialize a new Elasticsearch client.
	err := elastic.NewClient(config.Values.ES.URL)
	if err != nil {
		logger.Error("Failed to create Elasticsearch client", err)
		os.Exit(1)
	}

	// Create indices in Elasticsearch.
	err = elastic.Client.CreateMappings(indices)
	if err != nil {
		logger.Error("Failed to create index mappings for Elasticsearch", err)
		os.Exit(1)
	}

	logger.Info("Elasticsearch index created successfully")
}
