package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyrefresher/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyrefresher/internal/model"
)

type ProfileRepository interface {
	Count(profileID string) (int64, error)
	Add(profileJSON map[string]interface{}) error
	Update(
		profileID string,
		profileJSON map[string]interface{},
	) (map[string]interface{}, error)
	UpdateNodeID(profileID string, nodeID string) error
	FindLessThan(schemaName string, timestamp int64) ([]model.Profile, error)
	UpdateAccessTime(profileID string) error
	Delete(profileID string) error
}

func NewProfileRepository(client *mongo.Client) ProfileRepository {
	return &profileRepository{
		client: client,
	}
}

type profileRepository struct {
	client *mongo.Client
}

func (r *profileRepository) Count(profileID string) (int64, error) {
	filter := bson.M{"oid": profileID}

	count, err := r.client.Database(config.Values.Mongo.DBName).
		Collection(constant.MongoIndex.Profile).
		CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *profileRepository) Add(profileJSON map[string]interface{}) error {
	_, err := r.client.Database(config.Values.Mongo.DBName).
		Collection(constant.MongoIndex.Profile).
		InsertOne(context.Background(), profileJSON)

	if err != nil {
		return err
	}

	return nil
}

func (r *profileRepository) Update(
	profileID string,
	profileJSON map[string]interface{},
) (map[string]interface{}, error) {
	filter := bson.M{"oid": profileID}
	update := bson.M{"$set": profileJSON}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result := r.client.Database(config.Values.Mongo.DBName).
		Collection(constant.MongoIndex.Profile).
		FindOneAndUpdate(context.Background(), filter, update, opt)

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

func (r *profileRepository) UpdateNodeID(
	profileID string,
	nodeID string,
) error {
	filter := bson.M{"oid": profileID}
	update := bson.M{"$set": bson.M{"node_id": nodeID, "is_posted": false}}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result := r.client.Database(config.Values.Mongo.DBName).
		Collection(constant.MongoIndex.Profile).
		FindOneAndUpdate(context.Background(), filter, update, opt)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *profileRepository) FindLessThan(
	schemaName string,
	timestamp int64,
) ([]model.Profile, error) {
	schemaArray := [1]string{schemaName}

	filter := bson.M{
		"linked_schemas": bson.M{"$in": schemaArray},
		"metadata.sources.access_time": bson.M{
			"$lte": timestamp,
		},
	}

	var profiles []model.Profile
	cursor, err := r.client.Database(config.Values.Mongo.DBName).
		Collection(constant.MongoIndex.Profile).
		Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &profiles)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (r *profileRepository) UpdateAccessTime(profileID string) error {
	timestamp := time.Now().Unix()
	filter := bson.M{"oid": profileID}
	update := bson.M{
		"$set": bson.M{"metadata.sources.$[].access_time": timestamp},
	}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result := r.client.Database(config.Values.Mongo.DBName).
		Collection(constant.MongoIndex.Profile).
		FindOneAndUpdate(context.Background(), filter, update, opt)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *profileRepository) Delete(profileID string) error {
	filter := bson.M{"cuid": profileID}

	_, err := r.client.Database(config.Values.Mongo.DBName).
		Collection(constant.MongoIndex.Profile).
		DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
