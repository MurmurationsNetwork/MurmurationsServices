package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/internal/entity"
)

type NodeRepository interface {
	FindByStatuses(statuses []string) (entity.Nodes, error)
}

func NewNodeRepository(client *mongo.Client) NodeRepository {
	return &nodeRepository{
		client: client,
	}
}

type nodeRepository struct {
	client *mongo.Client
}

func (r *nodeRepository) FindByStatuses(
	statuses []string,
) (entity.Nodes, error) {
	filter := bson.M{"status": bson.M{"$in": statuses}}

	cur, err := r.client.Database(config.Conf.Mongo.DBName).
		Collection(constant.MongoIndex.Node).
		Find(context.Background(), filter)
	if err != nil {
		logger.Error(
			fmt.Sprintf("Error trying to find nodes with %v status", statuses),
			err,
		)
		return nil, err
	}

	var nodes entity.Nodes
	for cur.Next(context.TODO()) {
		var node entity.Node
		err := cur.Decode(&node)
		if err != nil {
			logger.Error(
				fmt.Sprintf(
					"Error trying to find nodes with %v status",
					statuses,
				),
				err,
			)
			return nil, err
		}
		nodes = append(nodes, &node)
	}

	if err := cur.Err(); err != nil {
		logger.Error(
			fmt.Sprintf("Error trying to find nodes with %v status", statuses),
			err,
		)
		return nil, err
	}

	cur.Close(context.TODO())

	return nodes, nil
}
