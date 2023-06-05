package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/entity/query"
)

type NodeRepository interface {
	RemoveValidationFailed(status string, timeBefore int64) error
	RemoveDeleted(status string, timeBefore int64) error
	RemoveES(status string, timeBefore int64) error
}

func NewNodeRepository(client *mongo.Client) NodeRepository {
	return &nodeRepository{
		client: client,
	}
}

type nodeRepository struct {
	client *mongo.Client
}

func (r *nodeRepository) RemoveValidationFailed(status string, timeBefore int64) error {
	filter := bson.M{
		"status": status,
		"createdAt": bson.M{
			"$lt": timeBefore,
		},
	}

	result, err := r.client.Database(config.Conf.Mongo.DBName).
		Collection(constant.MongoIndex.Node).
		DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount != 0 {
		logger.Info(
			fmt.Sprintf(
				"Delete %d nodes with %s status",
				result.DeletedCount,
				status,
			),
		)
	}

	return nil
}

func (r *nodeRepository) RemoveDeleted(status string, timeBefore int64) error {
	filter := bson.M{
		"status": status,
		"last_updated": bson.M{
			"$lte": timeBefore,
		},
	}

	result, err := r.client.Database(config.Conf.Mongo.DBName).
		Collection(constant.MongoIndex.Node).
		DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount != 0 {
		logger.Info(
			fmt.Sprintf(
				"Delete %d nodes with %s status",
				result.DeletedCount,
				status,
			),
		)
	}

	return nil
}

func (r *nodeRepository) RemoveES(status string, timeBefore int64) error {
	query := query.EsQuery{Status: &status, TimeBefore: &timeBefore}

	err := elastic.Client.DeleteMany(constant.ESIndex.Node, query.Build())
	if err != nil {
		return err
	}

	return nil
}
