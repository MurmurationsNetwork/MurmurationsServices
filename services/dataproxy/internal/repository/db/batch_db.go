package db

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

type BatchRepository interface {
	SaveUser(userId string, batchId string) error
	SaveProfile(profile map[string]interface{}) error
	SaveNodeId(profileId string, profile map[string]interface{}) error
	CheckUser(userId string, batchId string) (bool, error)
	GetProfilesByBatchId(batchId string) ([]map[string]interface{}, error)
	DeleteProfilesByBatchId(batchId string) error
	DeleteBatchId(batchId string) error
}

type batchRepository struct{}

func NewBatchRepository() BatchRepository {
	return &batchRepository{}
}

func (r *batchRepository) SaveUser(userId string, batchId string) error {
	doc := bson.M{
		"user_id":  userId,
		"batch_id": batchId,
	}
	_, err := mongo.Client.InsertOne(constant.MongoIndex.Batch, doc)
	if err != nil {
	}

	return nil
}

func (r *batchRepository) SaveProfile(profile map[string]interface{}) error {
	_, err := mongo.Client.InsertOne(constant.MongoIndex.Profile, profile)
	if err != nil {
	}

	return nil
}

func (r *batchRepository) SaveNodeId(profileId string, profile map[string]interface{}) error {
	filter := bson.M{"cuid": profileId}
	update := bson.M{"$set": profile}
	_, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Profile, filter, update)
	if err != nil {
	}

	return nil
}

func (r *batchRepository) CheckUser(userId string, batchId string) (bool, error) {
	filter := bson.M{"user_id": userId, "batch_id": batchId}
	count, err := mongo.Client.Count(constant.MongoIndex.Batch, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *batchRepository) GetProfilesByBatchId(batchId string) ([]map[string]interface{}, error) {
	filter := bson.M{"batch_id": batchId}
	cursor, err := mongo.Client.Find(constant.MongoIndex.Profile, filter)
	if err != nil {
		return nil, err
	}

	var profiles []map[string]interface{}
	if err = cursor.All(context.Background(), &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (r *batchRepository) DeleteProfilesByBatchId(batchId string) error {
	filter := bson.M{"batch_id": batchId}
	err := mongo.Client.DeleteMany(constant.MongoIndex.Profile, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *batchRepository) DeleteBatchId(batchId string) error {
	filter := bson.M{"batch_id": batchId}
	err := mongo.Client.DeleteOne(constant.MongoIndex.Batch, filter)
	if err != nil {
		return err
	}

	return nil
}
