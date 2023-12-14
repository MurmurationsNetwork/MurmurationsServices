package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyupdater/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyupdater/internal/repository/mongo"
)

type UpdatesService interface {
	Get(schemaName string) *model.Update
	Save(schemaName string, lastUpdated int64, apiEntry string) error
	Update(schemaName string, lastUpdated int64) error
	SaveError(schemaName string, errorMessage string) error
}

type updatesService struct {
	mongoRepo mongo.UpdateRepository
}

func NewUpdateService(mongoRepo mongo.UpdateRepository) UpdatesService {
	return &updatesService{
		mongoRepo: mongoRepo,
	}
}

func (svc *updatesService) Get(schemaName string) *model.Update {
	return svc.mongoRepo.Get(schemaName)
}

func (svc *updatesService) Save(
	schemaName string,
	lastUpdated int64,
	apiEntry string,
) error {
	return svc.mongoRepo.Save(schemaName, lastUpdated, apiEntry)
}

func (svc *updatesService) Update(schemaName string, lastUpdated int64) error {
	return svc.mongoRepo.Update(schemaName, lastUpdated)
}

func (svc *updatesService) SaveError(
	schemaName string,
	errorMessage string,
) error {
	return svc.mongoRepo.SaveError(schemaName, errorMessage)
}
