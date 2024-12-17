package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/revalidatenode/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/revalidatenode/internal/model"
)

// NodeRepository defines methods to interact with node data in MongoDB.
type NodeRepository interface {
	FindByStatuses(
		ctx context.Context,
		statuses []string,
		page, pageSize int,
	) ([]*model.Node, error)
}

// NewNodeRepository initializes and returns an instance of NodeRepository.
func NewNodeRepository(client *mongo.Client) NodeRepository {
	return &nodeRepository{client: client}
}

type nodeRepository struct {
	client *mongo.Client
}

// FindByStatuses retrieves paginated nodes with the given statuses.
func (r *nodeRepository) FindByStatuses(
	ctx context.Context,
	statuses []string,
	page, pageSize int,
) ([]*model.Node, error) {
	filter := bson.M{"status": bson.M{"$in": statuses}}
	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize))

	cur, err := r.client.Database(config.Values.Mongo.DBName).
		Collection(constant.MongoIndex.Node).
		Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var nodes []*model.Node
	for cur.Next(ctx) {
		var node model.Node
		err := cur.Decode(&node)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &node)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}
