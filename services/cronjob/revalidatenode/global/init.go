package global

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/config"
)

func Init() {
	config.Init()
	mongoInit()
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
		logger.Panic("error when trying to connect to MongoDB", err)
	}
	err = mongo.Client.Ping()
	if err != nil {
		logger.Panic("error when trying to ping the MongoDB", err)
	}
}

func natsInit() {
	err := nats.NewClient(
		config.Conf.Nats.ClusterID,
		config.Conf.Nats.ClientID,
		config.Conf.Nats.URL,
	)
	if err != nil {
		logger.Panic("error when trying to connect nats", err)
	}
}
