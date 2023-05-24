package db

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/domain/schema"
	"github.com/iancoleman/orderedmap"
	"go.mongodb.org/mongo-driver/bson"
)

type SchemaRepo interface {
	Get(schemaName string) (interface{}, []jsonapi.Error)
	Search() (schema.Schemas, []jsonapi.Error)
}

type schemaRepo struct{}

func NewSchemaRepo() SchemaRepo {
	return &schemaRepo{}
}

type SingleSchema struct {
	Description string `bson:"description"`
	FullSchema  bson.D `bson:"full_schema"`
}

func (r *schemaRepo) Get(schemaName string) (interface{}, []jsonapi.Error) {
	filter := bson.M{"name": schemaName}
	result := mongo.Client.FindOne(constant.MongoIndex.Schema, filter)

	var singleSchema SingleSchema
	err := result.Decode(&singleSchema)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, jsonapi.NewError(
				[]string{"Schema Not Found"},
				[]string{
					fmt.Sprintf(
						"Could not locate the following schema in the Library: %s",
						schemaName,
					),
				},
				nil,
				[]int{http.StatusNotFound},
			)
		}
		return nil, jsonapi.NewError(
			[]string{"Database Error"},
			[]string{"Error when trying to decode schema."},
			nil,
			[]int{http.StatusInternalServerError},
		)
	}

	fullSchema := convertBsonDToMap(singleSchema.FullSchema)

	return fullSchema, nil
}

func (r *schemaRepo) Search() (schema.Schemas, []jsonapi.Error) {
	filter := bson.M{}
	cur, err := mongo.Client.Find(constant.MongoIndex.Schema, filter)

	if err != nil {
		logger.Error("Error when trying to find schemas", err)
		return nil, jsonapi.NewError(
			[]string{"Database Error"},
			[]string{"Error when trying to find schemas."},
			nil,
			[]int{http.StatusInternalServerError},
		)
	}

	var schemas schema.Schemas
	for cur.Next(context.TODO()) {
		var schema schema.Schema
		err := cur.Decode(&schema)
		if err != nil {
			logger.Error("Error when trying to parse a schema from db", err)
			return nil, jsonapi.NewError(
				[]string{"Database Error"},
				[]string{"Error when trying to find schemas."},
				nil,
				[]int{http.StatusInternalServerError},
			)
		}
		schemas = append(schemas, &schema)
	}

	if err := cur.Err(); err != nil {
		logger.Error("Error when trying to find schemas", err)
		return nil, jsonapi.NewError(
			[]string{"Database Error"},
			[]string{"Error when trying to find schemas."},
			nil,
			[]int{http.StatusInternalServerError},
		)
	}

	cur.Close(context.TODO())

	return schemas, nil
}

func convertBsonDToMap(bsonD bson.D) *orderedmap.OrderedMap {
	result := orderedmap.New()
	for _, element := range bsonD {
		key := element.Key
		value := element.Value

		if innerDoc, ok := value.(bson.D); ok {
			result.Set(key, convertBsonDToMap(innerDoc))
		} else {
			result.Set(key, value)
		}
	}
	return result
}
