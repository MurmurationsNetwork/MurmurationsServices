package global

import (
	"os"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyrefresher/config"
)

func Init() {
	config.Init()
	mongoInit()
}

func mongoInit() {
	uri := mongo.GetURI(
		config.Values.Mongo.USERNAME,
		config.Values.Mongo.PASSWORD,
		config.Values.Mongo.HOST,
	)

	err := mongo.NewClient(uri, config.Values.Mongo.DBName)
	if err != nil {
		logger.Error("error when trying to connect to MongoDB", err)
		os.Exit(1)
	}
	err = mongo.Client.Ping()
	if err != nil {
		logger.Error("error when trying to ping the MongoDB", err)
		os.Exit(1)
	}
}
