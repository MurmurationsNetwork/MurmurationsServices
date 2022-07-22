package db

import (
	"context"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MappingRepository interface {
	Get(schemaName string) map[string]interface{}
}

func NewMappingRepository(client *mongo.Client) MappingRepository {
	return &mappingRepository{
		client: client,
	}
}

type mappingRepository struct {
	client *mongo.Client
}

func (r *mappingRepository) Get(schemaName string) map[string]interface{} {
	filter := bson.M{"schema": schemaName}

	result := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Mapping).FindOne(context.Background(), filter)

	var res map[string]interface{}
	result.Decode(&res)

	return res
}
