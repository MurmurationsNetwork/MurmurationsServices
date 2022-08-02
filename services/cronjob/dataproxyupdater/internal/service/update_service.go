package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/repository/db"
)

type UpdatesService interface {
	Get(schemaName string) *entity.Update
	Save(schemaName string, lastUpdated int64, apiEntry string) error
	Update(schemaName string, lastUpdated int64) error
	SaveError(schemaName string, errorMessage string) error
}

type updatesService struct {
	updateRepo db.UpdateRepository
}

func NewUpdateService(updateRepo db.UpdateRepository) UpdatesService {
	return &updatesService{
		updateRepo: updateRepo,
	}
}

func (svc *updatesService) Get(schemaName string) *entity.Update {
	return svc.updateRepo.Get(schemaName)
}

func (svc *updatesService) Save(schemaName string, lastUpdated int64, apiEntry string) error {
	return svc.updateRepo.Save(schemaName, lastUpdated, apiEntry)
}

func (svc *updatesService) Update(schemaName string, lastUpdated int64) error {
	return svc.updateRepo.Update(schemaName, lastUpdated)
}

func (svc *updatesService) SaveError(schemaName string, errorMessage string) error {
	return svc.updateRepo.SaveError(schemaName, errorMessage)
}
