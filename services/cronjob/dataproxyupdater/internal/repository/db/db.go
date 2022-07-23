package db

import (
	"context"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateRepository interface {
	Get(schemaName string) *entity.Update
	Save(schemaName string, lastUpdated int64, apiEntry string) error
	Update(schemaName string, lastUpdated int64) error
	SaveError(schemaName string, errorMessage string) error
}

func NewUpdateRepository(client *mongo.Client) UpdateRepository {
	return &updateRepository{
		client: client,
	}
}

type updateRepository struct {
	client *mongo.Client
}

func (r *updateRepository) Get(schemaName string) *entity.Update {
	filter := bson.M{"schema": schemaName}

	result := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Update).FindOne(context.Background(), filter)

	var res *entity.Update
	result.Decode(&res)

	return res
}

func (r *updateRepository) Save(schemaName string, lastUpdated int64, apiEntry string) error {
	filter := bson.M{"schema": schemaName, "last_updated": lastUpdated, "has_error": false, "api_entry": apiEntry}

	_, err := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Update).InsertOne(context.Background(), filter)

	if err != nil {
		return err
	}

	return nil
}

func (r *updateRepository) Update(schemaName string, lastUpdated int64) error {
	filter := bson.M{"schema": schemaName}
	update := bson.M{"$set": bson.M{"last_updated": lastUpdated}}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Update).FindOneAndUpdate(context.Background(), filter, update, opt)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *updateRepository) SaveError(schemaName string, errorMessage string) error {
	filter := bson.M{"schema": schemaName}
	update := bson.M{"$set": bson.M{"has_error": true, "error_message": errorMessage}}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Update).FindOneAndUpdate(context.Background(), filter, update, opt)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
