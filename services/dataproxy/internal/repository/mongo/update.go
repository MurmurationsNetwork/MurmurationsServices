package mongo

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/model"
)

type UpdateRepository interface {
	GetUpdate(schemaName string) (*model.Updates, resterr.RestErr)
}

type updateRepository struct{}

func NewUpdateRepository() UpdateRepository {
	return &updateRepository{}
}

func (r *updateRepository) GetUpdate(
	schemaName string,
) (*model.Updates, resterr.RestErr) {
	filter := bson.M{"schema": schemaName}

	result := mongo.Client.FindOne(constant.MongoIndex.Update, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, resterr.NewNotFoundError(
				fmt.Sprintf("Could not find schema: %s", schemaName),
			)
		}
		logger.Error("Error when trying to find an update", result.Err())
		return nil, resterr.NewInternalServerError(
			"Error when trying to find an update.",
			errors.New("database error"),
		)
	}

	var update *model.Updates
	err := result.Decode(&update)
	if err != nil {
		logger.Error("Error when trying to parse database response", err)
		return nil, resterr.NewInternalServerError(
			"Error when trying to find an update.",
			errors.New("database error"),
		)
	}

	return update, nil
}
