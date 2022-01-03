package db

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SchemaRepository interface {
	Update(schema *domain.Schema) error
}

func NewSchemaRepository() SchemaRepository {
	return &schemaRepository{}
}

type schemaRepository struct {
}

func (r *schemaRepository) Update(schema *domain.Schema) error {
	filter := bson.M{"name": schema.Name}
	update := bson.M{"$set": schema}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	_, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Schema, filter, update, opt)
	if err != nil {
		return err
	}

	return nil
}
