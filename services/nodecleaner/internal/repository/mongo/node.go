package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/nodecleaner/config"
)

const (
	StatusField      = "status"
	CreatedAtField   = "createdAt"
	LastUpdatedField = "last_updated"
	ExpiresField     = "expires"
)

// NodeRepository defines the operations available for manipulating nodes in a MongoDB repository.
type NodeRepository interface {
	RemoveByCreatedAt(
		ctx context.Context,
		status string,
		timeBefore int64,
	) error
	RemoveByLastUpdated(
		ctx context.Context,
		status string,
		timeBefore int64,
	) error
	UpdateStatusByExpiration(
		ctx context.Context,
		status string,
		timeBefore int64,
	) error
}

type nodeRepository struct {
	client *mongo.Client
}

// NewNodeRepository initializes and returns a new NodeRepository with the provided MongoDB client.
func NewNodeRepository(client *mongo.Client) NodeRepository {
	return &nodeRepository{client: client}
}

// RemoveByCreatedAt removes nodes with the specified status created before the given time.
func (r *nodeRepository) RemoveByCreatedAt(
	ctx context.Context,
	status string,
	timeBefore int64,
) error {
	return r.removeNodes(ctx, status, CreatedAtField, timeBefore)
}

// RemoveByLastUpdated removes nodes with the specified status that were last updated
// before the given time.
func (r *nodeRepository) RemoveByLastUpdated(
	ctx context.Context,
	status string,
	timeBefore int64,
) error {
	return r.removeNodes(ctx, status, LastUpdatedField, timeBefore)
}

// removeNodes is a helper function encapsulating the logic for removing nodes
// based on a time field, status, and timeBefore.
func (r *nodeRepository) removeNodes(
	ctx context.Context,
	status, timeField string,
	timeBefore int64,
) error {
	filter := bson.M{
		StatusField: status,
		timeField: bson.M{
			"$lt": timeBefore,
		},
	}

	result, err := r.client.Database(config.Values.Mongo.DBName).
		Collection(constant.MongoIndex.Node).
		DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("error removing nodes: %v", err)
	}

	if result.DeletedCount > 0 {
		fmt.Printf("Deleted %d nodes with %s status that were %s before %d\n",
			result.DeletedCount, status, timeField, timeBefore)
	}

	return nil
}

// UpdateStatusByExpiration updates the status of nodes with expired status before the given time.
func (r *nodeRepository) UpdateStatusByExpiration(
	ctx context.Context,
	status string,
	timeBefore int64,
) error {
	filter := bson.M{
		StatusField: status,
		ExpiresField: bson.M{
			"$lt": timeBefore,
		},
	}

	update := bson.M{
		"$set": bson.M{
			StatusField: constant.NodeStatus.Deleted,
		},
	}

	result, err := r.client.Database(config.Values.Mongo.DBName).
		Collection(constant.MongoIndex.Node).
		UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating nodes status: %v", err)
	}

	if result.ModifiedCount > 0 {
		fmt.Printf("Updated %d nodes with %s status to %s which expired before %d\n",
			result.ModifiedCount, status, constant.NodeStatus.Deleted, timeBefore)
	}

	return nil
}
