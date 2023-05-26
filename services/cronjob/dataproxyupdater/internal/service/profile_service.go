package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/repository/db"
)

type ProfilesService interface {
	Count(profileID string) (int64, error)
	Add(profileJSON map[string]interface{}) error
	Update(
		schemaName string,
		profileJSON map[string]interface{},
	) (map[string]interface{}, error)
	UpdateNodeID(profileID string, nodeID string) error
	GetNotPosted() ([]entity.Profile, error)
	UpdateIsPosted(nodeID string) error
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

func (svc *profilesService) GetNotPosted() ([]entity.Profile, error) {
	return svc.profileRepo.GetNotPosted()
}

func (svc *profilesService) UpdateIsPosted(nodeID string) error {
	return svc.profileRepo.UpdateIsPosted(nodeID)
}

func (svc *profilesService) Delete(profileID string) error {
	return svc.profileRepo.Delete(profileID)
}
