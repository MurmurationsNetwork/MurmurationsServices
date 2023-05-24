package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/repository/db"
)

type ProfilesService interface {
	Count(profileId string) (int64, error)
	Add(profileJson map[string]interface{}) error
	Update(
		schemaName string,
		profileJson map[string]interface{},
	) (map[string]interface{}, error)
	UpdateNodeId(profileId string, nodeId string) error
	GetNotPosted() ([]entity.Profile, error)
	UpdateIsPosted(nodeId string) error
	Delete(profileId string) error
}

type profilesService struct {
	profileRepo db.ProfileRepository
}

func NewProfileService(profileRepo db.ProfileRepository) ProfilesService {
	return &profilesService{
		profileRepo: profileRepo,
	}
}

func (svc *profilesService) Count(profileId string) (int64, error) {
	return svc.profileRepo.Count(profileId)
}

func (svc *profilesService) Add(profileJson map[string]interface{}) error {
	return svc.profileRepo.Add(profileJson)
}

func (svc *profilesService) Update(
	profileId string,
	profileJson map[string]interface{},
) (map[string]interface{}, error) {
	return svc.profileRepo.Update(profileId, profileJson)
}

func (svc *profilesService) UpdateNodeId(
	profileId string,
	nodeId string,
) error {
	return svc.profileRepo.UpdateNodeId(profileId, nodeId)
}

func (svc *profilesService) GetNotPosted() ([]entity.Profile, error) {
	return svc.profileRepo.GetNotPosted()
}

func (svc *profilesService) UpdateIsPosted(nodeId string) error {
	return svc.profileRepo.UpdateIsPosted(nodeId)
}

func (svc *profilesService) Delete(profileId string) error {
	return svc.profileRepo.Delete(profileId)
}
