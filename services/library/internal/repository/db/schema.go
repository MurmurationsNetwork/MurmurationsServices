package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/library"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/model"
)

// SchemaRepo defines the methods a SchemaRepo can perform.
type SchemaRepo interface {
	Get(schemaName string) (interface{}, error)
	Search() (*model.Schemas, error)
}

type schemaRepo struct{}

// NewSchemaRepo returns a new schema repository.
func NewSchemaRepo() SchemaRepo {
	return &schemaRepo{}
}

// Get retrieves a specific schema from the DB based on its name.
func (r *schemaRepo) Get(schemaName string) (interface{}, error) {
	filter := bson.M{"name": schemaName}
	result := mongo.Client.FindOne(constant.MongoIndex.Schema, filter)

	var singleSchema model.SingleSchema
	err := result.Decode(&singleSchema)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, library.SchemaNotFoundError{SchemaName: schemaName}
		}
		return nil, library.DatabaseError{Operation: "decode schema"}
	}

	return singleSchema.ToMap(), nil
}

// Search retrieves all schemas from the DB.
func (r *schemaRepo) Search() (*model.Schemas, error) {
	filter := bson.M{}
	cur, err := mongo.Client.Find(constant.MongoIndex.Schema, filter)

	if err != nil {
		logger.Error("Error when trying to find schemas", err)
		return nil, err
	}

	var schemas model.Schemas
	for cur.Next(context.TODO()) {
		var schema model.Schema
		err := cur.Decode(&schema)
		if err != nil {
			logger.Error("Error when trying to parse a schema from db", err)
			return nil, err
		}
		schemas = append(schemas, &schema)
	}

	if err := cur.Err(); err != nil {
		logger.Error("Error when trying to find schemas", err)
		return nil, err
	}

	cur.Close(context.TODO())

	return &schemas, nil
}
