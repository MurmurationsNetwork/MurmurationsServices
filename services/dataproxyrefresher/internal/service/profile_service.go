package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyrefresher/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyrefresher/internal/repository/mongo"
)

type ProfilesService interface {
	Count(profileID string) (int64, error)
	Add(profileJSON map[string]interface{}) error
	Update(
		profileID string,
		profileJSON map[string]interface{},
	) (map[string]interface{}, error)
	UpdateNodeID(profileID string, nodeID string) error
	FindLessThan(schemaName string, timestamp int64) ([]model.Profile, error)
	UpdateAccessTime(profileID string) error
	Delete(profileID string) error
}

type profilesService struct {
	mongoRepo mongo.ProfileRepository
}

func NewProfileService(mongoRepo mongo.ProfileRepository) ProfilesService {
	return &profilesService{
		mongoRepo: mongoRepo,
	}
}

func (svc *profilesService) Count(profileID string) (int64, error) {
	return svc.mongoRepo.Count(profileID)
}

func (svc *profilesService) Add(profileJSON map[string]interface{}) error {
	return svc.mongoRepo.Add(profileJSON)
}

func (svc *profilesService) Update(
	profileID string,
	profileJSON map[string]interface{},
) (map[string]interface{}, error) {
	return svc.mongoRepo.Update(profileID, profileJSON)
}

func (svc *profilesService) UpdateNodeID(
	profileID string,
	nodeID string,
) error {
	return svc.mongoRepo.UpdateNodeID(profileID, nodeID)
}

func (svc *profilesService) FindLessThan(
	schemaName string,
	timestamp int64,
) ([]model.Profile, error) {
	return svc.mongoRepo.FindLessThan(schemaName, timestamp)
}

func (svc *profilesService) UpdateAccessTime(profileID string) error {
	return svc.mongoRepo.UpdateAccessTime(profileID)
}

func (svc *profilesService) Delete(profileID string) error {
	return svc.mongoRepo.Delete(profileID)
}
