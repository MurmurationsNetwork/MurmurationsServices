package db

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type BatchRepository interface {
	SaveUser(userId string, batchId string) error
	SaveProfile(profile map[string]interface{}) error
	SaveNodeId(profileId string, profile map[string]interface{}) error
	CheckUser(userId string, batchId string) (bool, error)
	GetProfileByCuid(cuid string) (map[string]interface{}, error)
	GetProfilesByBatchId(batchId string) ([]map[string]interface{}, error)
	GetProfileOidsAndHashesByBatchId(batchId string) (map[string][2]string, error)
	UpdateProfile(profileId string, profile map[string]interface{}) error
	DeleteProfileByCuid(cuid string) error
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

func (r *batchRepository) GetProfileByCuid(cuid string) (map[string]interface{}, error) {
	filter := bson.M{"cuid": cuid}
	doc := mongo.Client.FindOne(constant.MongoIndex.Profile, filter)
	if doc.Err() != nil {
		return nil, doc.Err()
	}

	var profile map[string]interface{}
	if err := doc.Decode(&profile); err != nil {
		return nil, err
	}

	return profile, nil
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

func (r *batchRepository) GetProfileOidsAndHashesByBatchId(batchId string) (map[string][2]string, error) {
	filter := bson.M{"batch_id": batchId}
	opts := options.Find().SetProjection(bson.D{{"_id", 0}, {"oid", 1}, {"cuid", 1}, {"source_data_hash", 1}})
	cursor, err := mongo.Client.Find(constant.MongoIndex.Profile, filter, opts)
	if err != nil {
		return nil, err
	}

	var profiles []map[string]string
	if err = cursor.All(context.Background(), &profiles); err != nil {
		return nil, err
	}

	profileOidsAndHashes := make(map[string][2]string)
	for _, profile := range profiles {
		profileOidsAndHashes[profile["oid"]] = [2]string{profile["cuid"], profile["source_data_hash"]}
	}

	return profileOidsAndHashes, nil
}

func (r *batchRepository) UpdateProfile(profileId string, profile map[string]interface{}) error {
	filter := bson.M{"cuid": profileId}
	update := bson.M{"$set": profile}
	_, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Profile, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *batchRepository) DeleteProfileByCuid(cuid string) error {
	filter := bson.M{"cuid": cuid}
	err := mongo.Client.DeleteOne(constant.MongoIndex.Profile, filter)
	if err != nil {
		return err
	}

	return nil
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