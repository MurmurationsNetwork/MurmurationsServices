package main

import (
	"context"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleanup/internal/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleanup/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleanup/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	config.Init()
}

func main() {
	var client *mongo.Client

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Conf.Mongo.URL))
	if err != nil {
		logger.Panic("error when trying to connect to MongoDB", err)
		return
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		logger.Panic("trying to re-connect MongoDB %s \n", err)
		return
	}

	svc := service.NewNodeService(db.NewNodeRepository(client))
	err = svc.Remove()
	if err != nil {
		logger.Panic("error when trying to delete nodes", err)
		return
	}

	if err := client.Disconnect(context.Background()); err != nil {
		logger.Panic("error when trying to disconnect from MongoDB", err)
		return
	}
}
