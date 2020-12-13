package mongodb

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/config"
)

func Init() {
	err := mongo.NewClient(config.Conf.Mongo.URL, config.Conf.Mongo.DBName)
	if err != nil {
		logger.Panic("error when trying to connect to MongoDB", err)
	}
	err = mongo.Client.Ping()
	if err != nil {
		logger.Panic("error when trying to ping the MongoDB", err)
	}
}
