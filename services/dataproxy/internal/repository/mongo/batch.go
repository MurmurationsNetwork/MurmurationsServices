package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/model"
)

type BatchRepository interface {
	GetBatchesByUserID(userID string) ([]model.Batch, error)
	SaveUser(
		userID string,
		batchTitle string,
		batchID string,
		schemas []string,
		totalNodes int,
	) error
	UpdateBatchTitle(batchTitle string, batchID string) ([]string, error)
	SaveProfile(profile map[string]interface{}) error
	SaveNodeID(profileID string, profile map[string]interface{}) error
	CheckUser(userID string, batchID string) (bool, error)
	GetProfileByCuid(cuid string) (map[string]interface{}, error)
	GetProfilesByBatchID(batchID string) ([]map[string]interface{}, error)
	GetProfileOidsAndHashesByBatchID(
		batchID string,
	) (map[string][2]string, error)
	UpdateProfile(profileID string, profile map[string]interface{}) error
	DeleteProfileByCuid(cuid string) error
	DeleteProfilesByBatchID(batchID string) error
	DeleteBatchID(batchID string) error
	UpdateBatchError(batchID string, errorMessage string) error
	UpdateBatchProgress(batchID string, progress int) error
	UpdateBatchStatus(batchID string, status string) error
}

type batchRepository struct{}

func NewBatchRepository() BatchRepository {
	return &batchRepository{}
}

func (r *batchRepository) GetBatchesByUserID(
	userID string,
) ([]model.Batch, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetProjection(bson.M{"user_id": 0})
	cursor, err := mongo.Client.Find(constant.MongoIndex.Batch, filter, opts)
	if err != nil {
		return nil, err
	}

	batches := make([]model.Batch, 0)
	for cursor.Next(context.Background()) {
		var batch model.Batch
		if err := cursor.Decode(&batch); err != nil {
			return nil, err
		}
		batches = append(batches, batch)
	}

	return batches, nil
}

func (r *batchRepository) SaveUser(
	userID string,
	batchTitle string,
	batchID string,
	schemas []string,
	totalNodes int,
) error {
	doc := model.Batch{
		UserID:     userID,
		Title:      batchTitle,
		BatchID:    batchID,
		Schemas:    schemas,
		TotalNodes: totalNodes,
	}
	_, err := mongo.Client.InsertOne(constant.MongoIndex.Batch, doc)
	if err != nil {
		return err
	}

	return nil
}

func (r *batchRepository) UpdateBatchTitle(
	batchTitle string,
	batchID string,
) ([]string, error) {
	filter := bson.M{"batch_id": batchID}
	update := bson.M{"$set": bson.M{"title": batchTitle}}
	result, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Batch,
		filter,
		update,
	)
	if err != nil {
		return nil, err
	}

	var batch model.Batch
	if err := result.Decode(&batch); err != nil {
		return nil, err
	}
	return batch.Schemas, nil
}

func (r *batchRepository) SaveProfile(profile map[string]interface{}) error {
	_, err := mongo.Client.InsertOne(constant.MongoIndex.Profile, profile)
	if err != nil {
		return err
	}

	return nil
}

func (r *batchRepository) SaveNodeID(
	profileID string,
	profile map[string]interface{},
) error {
	filter := bson.M{"cuid": profileID}
	update := bson.M{"$set": profile}
	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Profile,
		filter,
		update,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *batchRepository) CheckUser(
	userID string,
	batchID string,
) (bool, error) {
	filter := bson.M{"user_id": userID, "batch_id": batchID}
	count, err := mongo.Client.Count(constant.MongoIndex.Batch, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *batchRepository) GetProfileByCuid(
	cuid string,
) (map[string]interface{}, error) {
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

func (r *batchRepository) GetProfilesByBatchID(
	batchID string,
) ([]map[string]interface{}, error) {
	filter := bson.M{"batch_id": batchID}
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

func (r *batchRepository) GetProfileOidsAndHashesByBatchID(
	batchID string,
) (map[string][2]string, error) {
	filter := bson.M{"batch_id": batchID}
	opts := options.Find().
		SetProjection(bson.D{{Key: "_id", Value: 0}, {Key: "oid", Value: 1}, {Key: "cuid", Value: 1}, {Key: "source_data_hash", Value: 1}})
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
		profileOidsAndHashes[profile["oid"]] = [2]string{
			profile["cuid"],
			profile["source_data_hash"],
		}
	}

	return profileOidsAndHashes, nil
}

func (r *batchRepository) UpdateProfile(
	profileID string,
	profile map[string]interface{},
) error {
	filter := bson.M{"cuid": profileID}

	// This part is used to remove the fields that are not in the new profile
	// get field name from mongoDB
	result := mongo.Client.FindOne(constant.MongoIndex.Profile, filter)
	var oldProfile map[string]interface{}
	if err := result.Decode(&oldProfile); err != nil {
		return err
	}
	// we don't need _id, __v, cuid
	delete(oldProfile, "_id")
	delete(oldProfile, "__v")
	delete(oldProfile, "cuid")

	unsetKeys := make(map[string]interface{})
	for key := range oldProfile {
		if _, ok := profile[key]; !ok {
			unsetKeys[key] = nil
		}
	}

	update := bson.M{"$set": profile, "$unset": unsetKeys}
	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Profile,
		filter,
		update,
	)
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

func (r *batchRepository) DeleteProfilesByBatchID(batchID string) error {
	filter := bson.M{"batch_id": batchID}
	err := mongo.Client.DeleteMany(constant.MongoIndex.Profile, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *batchRepository) DeleteBatchID(batchID string) error {
	filter := bson.M{"batch_id": batchID}
	err := mongo.Client.DeleteOne(constant.MongoIndex.Batch, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *batchRepository) UpdateBatchError(
	batchID string,
	errorMessage string,
) error {
	filter := bson.M{"batch_id": batchID}
	update := bson.M{"$set": bson.M{"error": errorMessage}}
	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Batch,
		filter,
		update,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *batchRepository) UpdateBatchProgress(
	batchID string,
	progress int,
) error {
	filter := bson.M{"batch_id": batchID}
	update := bson.M{"$set": bson.M{"processed_nodes": progress}}
	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Batch,
		filter,
		update,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *batchRepository) UpdateBatchStatus(
	batchID string,
	status string,
) error {
	filter := bson.M{"batch_id": batchID}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Batch,
		filter,
		update,
	)
	if err != nil {
		return err
	}
	return nil
}
