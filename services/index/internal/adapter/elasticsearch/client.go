package elasticsearch

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

func Init() {
	err := elastic.NewClient(config.Conf.Es.URL)
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
