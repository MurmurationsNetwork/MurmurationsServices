package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/internal/repository/db"
)

type ProfilesService interface {
	Count(profileID string) (int64, error)
	Add(profileJSON map[string]interface{}) error
	Update(
		profileID string,
		profileJSON map[string]interface{},
	) (map[string]interface{}, error)
	UpdateNodeID(profileID string, nodeID string) error
	FindLessThan(schemaName string, timestamp int64) ([]entity.Profile, error)
	UpdateAccessTime(profileID string) error
	Delete(profileID string) error
}

type profilesService struct {
	profileRepo db.ProfileRepository
}

func NewProfileService(profileRepo db.ProfileRepository) ProfilesService {
	return &profilesService{
		profileRepo: profileRepo,
	}
}

func (svc *profilesService) Count(profileID string) (int64, error) {
	return svc.profileRepo.Count(profileID)
}

func (svc *profilesService) Add(profileJSON map[string]interface{}) error {
	return svc.profileRepo.Add(profileJSON)
}

func (svc *profilesService) Update(
	profileID string,
	profileJSON map[string]interface{},
) (map[string]interface{}, error) {
	return svc.profileRepo.Update(profileID, profileJSON)
}

func (svc *profilesService) UpdateNodeID(
	profileID string,
	nodeID string,
) error {
	return svc.profileRepo.UpdateNodeID(profileID, nodeID)
}

func (svc *profilesService) FindLessThan(
	schemaName string,
	timestamp int64,
) ([]entity.Profile, error) {
	return svc.profileRepo.FindLessThan(schemaName, timestamp)
}

func (svc *profilesService) UpdateAccessTime(profileID string) error {
	return svc.profileRepo.UpdateAccessTime(profileID)
}

func (svc *profilesService) Delete(profileID string) error {
	return svc.profileRepo.Delete(profileID)
}
