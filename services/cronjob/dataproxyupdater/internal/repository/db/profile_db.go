package db

import (
	"context"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfileRepository interface {
	Count(profileId string) (int64, error)
	Add(profileJson map[string]interface{}) error
	Update(schemaName string, profileJson map[string]interface{})
}

func NewProfileRepository(client *mongo.Client) ProfileRepository {
	return &profileRepository{
		client: client,
	}
}

type profileRepository struct {
	client *mongo.Client
}

func (r *profileRepository) Count(profileId string) (int64, error) {
	filter := bson.M{"oid": profileId}

	count, err := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Update).CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *profileRepository) Add(profileJson map[string]interface{}) error {
	_, err := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Profile).InsertOne(context.Background(), profileJson)

	if err != nil {
		return err
	}

	return nil
}

func (r *profileRepository) Update(profileId string, profileJson map[string]interface{}) {
	filter := bson.M{"oid": profileId}
	update := bson.M{"$set": profileJson}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Profile).FindOneAndUpdate(context.Background(), filter, update, opt)
}
