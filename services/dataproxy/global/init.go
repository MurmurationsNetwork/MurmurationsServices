package global

import (
	"os"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
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
		logger.Error("error when trying to connect to MongoDB", err)
		os.Exit(1)
	}
	err = mongo.Client.Ping()
	if err != nil {
		logger.Error("error when trying to ping the MongoDB", err)
		os.Exit(1)
	}

	// Create an index on the `cuid` field in the `profiles` collection
	err = mongo.Client.CreateUniqueIndex("profiles", "cuid")
	if err != nil {
		logger.Error("error when trying to create index on `cuid` field in `profiles` collection", err)
		os.Exit(1)
	}
}
