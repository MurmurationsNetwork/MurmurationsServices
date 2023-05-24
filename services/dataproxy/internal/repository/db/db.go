package db

import (
	"errors"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MappingRepository interface {
	Save(mapping map[string]string) resterr.RestErr
}

type mappingRepository struct{}

func NewMappingRepository() MappingRepository {
	return &mappingRepository{}
}

func (r *mappingRepository) Save(mapping map[string]string) resterr.RestErr {
	filter := bson.M{"schema": mapping["schema"]}
	update := bson.M{"$set": mapping}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Mapping,
		filter,
		update,
		opt,
	)
	if err != nil {
		logger.Error("Error when trying to create a node", err)
		return resterr.NewInternalServerError(
			"Error when trying to add a node.",
			errors.New("database error"),
		)
	}

	return nil
}
