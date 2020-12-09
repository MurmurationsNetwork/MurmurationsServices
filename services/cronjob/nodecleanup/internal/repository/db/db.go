package db

import (
	"context"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleanup/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NodeRepository interface {
	Remove(status string) error
}

func NewNodeRepository(client *mongo.Client) NodeRepository {
	return &nodeRepository{
		client: client,
	}
}

type nodeRepository struct {
	client *mongo.Client
}

func (r *nodeRepository) Remove(status string) error {
	filter := bson.M{"status": status}
	result, err := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Node).DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Delete %d nodes with %s status", result.DeletedCount, status))
	return nil
}
