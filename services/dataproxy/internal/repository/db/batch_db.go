package db

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type BatchRepository interface {
	SaveUser(userCuid string, batchCuid string) error
	SaveProfile(profile map[string]interface{}) error
	SaveNodeId(profileId string, profile map[string]interface{}) error
}

type batchRepository struct{}

func NewBatchRepository() BatchRepository {
	return &batchRepository{}
}

func (r *batchRepository) SaveUser(userCuid string, batchCuid string) error {
	doc := bson.M{
		"user_id":  userCuid,
		"batch_id": batchCuid,
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
