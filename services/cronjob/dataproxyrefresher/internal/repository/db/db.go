package db

import (
	"context"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ProfileRepository interface {
	Count(profileId string) (int64, error)
	Add(profileJson map[string]interface{}) error
	Update(profileId string, profileJson map[string]interface{}) (map[string]interface{}, error)
	UpdateNodeId(profileId string, nodeId string) error
	FindLessThan(schemaName string, timestamp int64) ([]entity.Profile, error)
	UpdateAccessTime(profileId string) error
	Delete(profileId string) error
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

	count, err := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Profile).CountDocuments(context.Background(), filter)
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

func (r *profileRepository) Update(profileId string, profileJson map[string]interface{}) (map[string]interface{}, error) {
	filter := bson.M{"oid": profileId}
	update := bson.M{"$set": profileJson}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Profile).FindOneAndUpdate(context.Background(), filter, update, opt)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var profile map[string]interface{}
	err := result.Decode(&profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *profileRepository) UpdateNodeId(profileId string, nodeId string) error {
	filter := bson.M{"oid": profileId}
	update := bson.M{"$set": bson.M{"node_id": nodeId, "is_posted": false}}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Profile).FindOneAndUpdate(context.Background(), filter, update, opt)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *profileRepository) FindLessThan(schemaName string, timestamp int64) ([]entity.Profile, error) {
	schemaArray := [1]string{schemaName}

	filter := bson.M{
		"linked_schemas": bson.M{"$in": schemaArray},
		"metadata.sources.access_time": bson.M{
			"$lte": timestamp,
		},
	}

	var profiles []entity.Profile
	cursor, err := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Profile).Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &profiles)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (r *profileRepository) UpdateAccessTime(profileId string) error {
	timestamp := time.Now().Unix()
	filter := bson.M{"oid": profileId}
	update := bson.M{"$set": bson.M{"metadata.sources.$[].access_time": timestamp}}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Profile).FindOneAndUpdate(context.Background(), filter, update, opt)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *profileRepository) Delete(profileId string) error {
	filter := bson.M{"cuid": profileId}

	_, err := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Profile).DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
