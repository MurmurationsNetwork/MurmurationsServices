package db

import (
	"context"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/domain/schema"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type SchemaRepo interface {
	Get(schemaName string) (interface{}, []jsonapi.Error)
	Search() (schema.Schemas, []jsonapi.Error)
}

type schemaRepo struct{}

func NewSchemaRepo() SchemaRepo {
	return &schemaRepo{}
}

func (r *schemaRepo) Get(schemaName string) (interface{}, []jsonapi.Error) {
	filter := bson.M{"name": schemaName}
	result := mongo.Client.FindOne(constant.MongoIndex.Schema, filter)

	var singleSchema map[string]interface{}
	err := result.Decode(&singleSchema)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, jsonapi.NewError([]string{"Schema Not Found"}, []string{fmt.Sprintf("Could not locate the following schema in the library: %s", schemaName)}, nil, []int{http.StatusNotFound})
		}
		return nil, jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to decode schema."}, nil, []int{http.StatusInternalServerError})
	}
	return singleSchema["full_schema"], nil
}

func (r *schemaRepo) Search() (schema.Schemas, []jsonapi.Error) {
	filter := bson.M{}
	cur, err := mongo.Client.Find(constant.MongoIndex.Schema, filter)

	if err != nil {
		logger.Error("Error when trying to find schemas", err)
		return nil, jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to find schemas."}, nil, []int{http.StatusInternalServerError})
	}

	var schemas schema.Schemas
	for cur.Next(context.TODO()) {
		var schema schema.Schema
		err := cur.Decode(&schema)
		if err != nil {
			logger.Error("Error when trying to parse a schema from db", err)
			return nil, jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to find schemas."}, nil, []int{http.StatusInternalServerError})
		}
		schemas = append(schemas, &schema)
	}

	if err := cur.Err(); err != nil {
		logger.Error("Error when trying to find schemas", err)
		return nil, jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to find schemas."}, nil, []int{http.StatusInternalServerError})
	}

	cur.Close(context.TODO())

	return schemas, nil
}
