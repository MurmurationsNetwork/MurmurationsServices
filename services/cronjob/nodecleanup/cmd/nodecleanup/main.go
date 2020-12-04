package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	var client *mongo.Client

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("INDEX_MONGO_URL")))
	if err != nil {
		logger.Panic("error when trying to connect to MongoDB", err)
		return
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		logger.Panic("trying to re-connect MongoDB %s \n", err)
		return
	}

	filter := bson.M{"status": constant.NodeStatus.ValidationFailed}
	// TODO: Abstract MongoDB operations.
	result, err := client.Database("murmurations").Collection("nodes").DeleteMany(context.Background(), filter)
	if err != nil {
		logger.Panic("error when trying to delete nodes", err)
		return
	}

	logger.Info(fmt.Sprintf("Delete %d nodes with %s status", result.DeletedCount, constant.NodeStatus.ValidationFailed))

	if err := client.Disconnect(context.Background()); err != nil {
		logger.Panic("error when trying to disconnect from MongoDB", err)
		return
	}
}
