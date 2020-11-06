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
			"mappings": {
				"properties": {
					"profileUrl": {
						"type": "keyword"
					},
					"linkedSchemas": {
						"type": "keyword"
					},
					"name": {
						"type": "keyword"
					},
					"url": {
						"type": "keyword"
					},
					"mission": {
						"type": "text"
					},
					"keywords": {
						"type": "keyword"
					},
					"geolocation": {
						"type": "geo_point"
					},
					"maplocation": {
						"properties": {
							"locality": {
								"type": "keyword"
							},
							"region": {
								"type": "keyword"
							},
							"country": {
								"type": "keyword"
							}
						}
					},
					"lastChecked": {
						"type": "date",
						"format": "epoch_second"
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
