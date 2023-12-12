package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/internal/model"
)

type SchemaRepository interface {
	Update(schema *model.Schema) error
}

func NewSchemaRepository() SchemaRepository {
	return &schemaRepository{}
}

type schemaRepository struct {
}

func (r *schemaRepository) Update(schema *model.Schema) error {
	filter := bson.M{"name": schema.Name}
	update := bson.M{"$set": schema}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Schema,
		filter,
		update,
		opt,
	)
	if err != nil {
		return err
	}

	return nil
}
