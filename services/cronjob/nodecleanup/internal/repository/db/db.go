package db

import (
	"context"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
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
	result, err := r.client.Database("murmurations").Collection("nodes").DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Delete %d nodes with %s status", result.DeletedCount, status))
	return nil
}
