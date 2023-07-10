package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/index"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
)

// NodeMongo interface represents a set of methods required for node database operations.
type NodeRepository interface {
	Add(node *model.Node) error
	GetByID(nodeID string) (*model.Node, error)
	Update(node *model.Node) error
	Delete(node *model.Node) error
	SoftDelete(node *model.Node) error
}

// NewRepository function returns a new NodeRepository.
func NewNodeRepository() NodeRepository {
	return &nodeRepository{}
}

// nodeRepository struct implements NodeRepository interface.
type nodeRepository struct {
}

// Add method adds or updates a node in the database.
func (r *nodeRepository) Add(node *model.Node) error {
	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Node,
		filter,
		update,
		opt,
	)
	if err != nil {
		return index.DatabaseError{
			Message: "Error occurred during node upsert operation",
			Err:     err,
		}
	}

	var updated model.Node
	err = result.Decode(&updated)
	if err != nil {
		return index.DatabaseError{
			Message: "Error occurred during decoding of updated node",
			Err:     err,
		}
	}

	node.Version = updated.Version

	return nil
}

// GetByID method retrieves a node from the database using its id.
func (r *nodeRepository) GetByID(
	nodeID string,
) (*model.Node, error) {
	filter := bson.M{"_id": nodeID}

	result := mongo.Client.FindOne(constant.MongoIndex.Node, filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, index.NotFoundError{
				Err: err,
			}
		}
		return nil, index.DatabaseError{
			Message: "Error when trying to find a node",
			Err:     err,
		}
	}

	var node model.Node
	err := result.Decode(&node)
	if err != nil {
		return nil, index.DatabaseError{
			Message: "Error when trying to find a node",
			Err:     err,
		}
	}

	return &node, nil
}

func (r *nodeRepository) Update(node *model.Node) error {
	filter := bson.M{"_id": node.ID}

	if node.Version != nil {
		filter["__v"] = node.Version
		// Unset the version to prevent setting it.
		node.Version = nil
	}

	update := bson.M{"$set": node}

	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Node,
		filter,
		update,
	)
	if err != nil {
		// Update the document only if the version matches.
		// If the version does not match, it's an expected concurrent issue.
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return index.DatabaseError{
			Message: "Error when trying to update a node",
			Err:     err,
		}
	}

	return nil
}

func (r *nodeRepository) Delete(node *model.Node) error {
	filter := bson.M{"_id": node.ID}

	err := mongo.Client.DeleteOne(constant.MongoIndex.Node, filter)
	if err != nil {
		return index.DatabaseError{
			Err: err,
		}
	}

	return nil
}

func (r *nodeRepository) SoftDelete(node *model.Node) error {
	err := r.setDeleted(node)
	if err != nil {
		return index.DatabaseError{
			Err: err,
		}
	}
	return nil
}

func (r *nodeRepository) setDeleted(node *model.Node) error {
	node.Version = nil
	node.Status = constant.NodeStatus.Deleted
	currentTime := time.Now().Unix()
	node.LastUpdated = &currentTime

	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}

	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Node,
		filter,
		update,
	)
	if err != nil {
		logger.Error("Error when trying to update a node", err)
		return err
	}

	return nil
}
