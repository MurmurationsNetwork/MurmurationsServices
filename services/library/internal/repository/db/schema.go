package db

import (
	"context"
	"errors"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/domain/schema"
	"go.mongodb.org/mongo-driver/bson"
)

type SchemaRepo interface {
	Get(schemaName string) (interface{}, resterr.RestErr)
	Search() (schema.Schemas, resterr.RestErr)
}

type schemaRepo struct{}

func NewSchemaRepo() SchemaRepo {
	return &schemaRepo{}
}

func (r *schemaRepo) Get(schemaName string) (interface{}, resterr.RestErr) {
	filter := bson.M{"name": schemaName}
	result := mongo.Client.FindOne(constant.MongoIndex.Schema, filter)

	var singleSchema map[string]interface{}
	err := result.Decode(&singleSchema)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, resterr.NewNotFoundError("schema not found.")
		}
		return nil, resterr.NewInternalServerError("Error when trying to decode schema.", err)
	}
	return singleSchema["full_schema"], nil
}

func (r *schemaRepo) Search() (schema.Schemas, resterr.RestErr) {
	filter := bson.M{}

	cur, err := mongo.Client.Find(constant.MongoIndex.Schema, filter)
	if err != nil {
		logger.Error("Error when trying to find schemas", err)
		return nil, resterr.NewInternalServerError("Error when trying to find schemas.", errors.New("database error"))
	}

	var schemas schema.Schemas
	for cur.Next(context.TODO()) {
		var schema schema.Schema
		err := cur.Decode(&schema)
		if err != nil {
			logger.Error("Error when trying to parse a schema from db", err)
			return nil, resterr.NewInternalServerError("Error when trying to find schemas.", errors.New("database error"))
		}
		schemas = append(schemas, &schema)
	}

	if err := cur.Err(); err != nil {
		logger.Error("Error when trying to find schemas", err)
		return nil, resterr.NewInternalServerError("Error when trying to find schemas.", errors.New("database error"))
	}

	cur.Close(context.TODO())

	return schemas, nil
}
