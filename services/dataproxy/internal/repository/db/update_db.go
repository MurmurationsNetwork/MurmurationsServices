package db

import (
	"errors"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateRepository interface {
	GetUpdate(schemaName string) (*entity.Updates, resterr.RestErr)
}

type updateRepository struct{}

func NewUpdateRepository() UpdateRepository {
	return &updateRepository{}
}

func (r *updateRepository) GetUpdate(schemaName string) (*entity.Updates, resterr.RestErr) {
	filter := bson.M{"schema": schemaName}

	result := mongo.Client.FindOne(constant.MongoIndex.Update, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, resterr.NewNotFoundError(fmt.Sprintf("Could not find schema: %s", schemaName))
		}
		logger.Error("Error when trying to find an update", result.Err())
		return nil, resterr.NewInternalServerError("Error when trying to find an update.", errors.New("database error"))
	}

	var update *entity.Updates
	err := result.Decode(&update)
	if err != nil {
		logger.Error("Error when trying to parse database response", err)
		return nil, resterr.NewInternalServerError("Error when trying to find an update.", err)
	}

	return update, nil
}
