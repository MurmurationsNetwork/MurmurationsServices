package global

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/config"
)

func Init() {
	config.Init()
	mongoInit()
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
