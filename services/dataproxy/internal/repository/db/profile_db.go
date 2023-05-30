package db

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/resterr"
)

type ProfileRepository interface {
	GetProfile(profileID string) (map[string]interface{}, resterr.RestErr)
}

type profileRepository struct{}

func NewProfileRepository() ProfileRepository {
	return &profileRepository{}
}

func (r *profileRepository) GetProfile(
	profileID string,
) (map[string]interface{}, resterr.RestErr) {
	filter := bson.M{"cuid": profileID}

	result := mongo.Client.FindOne(constant.MongoIndex.Profile, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, resterr.NewNotFoundError(
				fmt.Sprintf("Could not find profile_id: %s", profileID),
			)
		}
		logger.Error("Error when trying to find a node", result.Err())
		return nil, resterr.NewInternalServerError(
			"Error when trying to find a node.",
			errors.New("database error"),
		)
	}

	var profile map[string]interface{}
	err := result.Decode(&profile)
	if err != nil {
		logger.Error(
			"Error when trying to parse database response",
			result.Err(),
		)
		return nil, resterr.NewInternalServerError(
			"Error when trying to find a node.",
			errors.New("database error"),
		)
	}

	return profile, nil
}
