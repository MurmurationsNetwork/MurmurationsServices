package elasticsearch

import (
	"os"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
)

func Init() {
	err := elastic.NewClient(os.Getenv("ELASTICSEARCH_URL"))
	if err != nil {
		logger.Panic("error when trying to ping the ElasticSearch", err)
		return
	}
	err = elastic.Client.CreateMappings(indices)
	if err != nil {
		logger.Panic("error when trying to create index for ElasticSearch", err)
		return
	}
}
