package mongodb

import (
	"os"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
)

func Init() {
	err := mongo.NewClient(os.Getenv("MONGO_URI"), "murmurations")
	if err != nil {
		logger.Panic("error when trying to connect to MongoDB", err)
	}
	err = mongo.Client.Ping()
	if err != nil {
		logger.Panic("error when trying to ping the MongoDB", err)
	}
}
