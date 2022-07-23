package db

import (
	"context"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxycleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxycleaner/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileRepository interface {
	FindLessThan(timestamp int64) ([]entity.Profile, error)
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

func (r *profileRepository) FindLessThan(timestamp int64) ([]entity.Profile, error) {
	filter := bson.M{
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

func (r *profileRepository) Delete(profileId string) error {
	filter := bson.M{"cuid": profileId}

	_, err := r.client.Database(config.Conf.Mongo.DBName).Collection(constant.MongoIndex.Profile).DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
