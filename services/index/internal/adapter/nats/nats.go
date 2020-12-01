package nats

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

func Init() {
	err := nats.NewClient(config.Conf.Nats.ClusterID, config.Conf.Nats.ClientID, config.Conf.Nats.URL)
	if err != nil {
		logger.Panic("error when trying to connect nats", err)
	}
}
