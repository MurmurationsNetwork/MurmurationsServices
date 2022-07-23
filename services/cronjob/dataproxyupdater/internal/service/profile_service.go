package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/repository/db"
)

type ProfilesService interface {
	Count(profileId string) (int64, error)
	Add(profileJson map[string]interface{}) error
	Update(schemaName string, profileJson map[string]interface{}) error
	UpdateNodeId(profileId string, nodeId string) error
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

func (svc *profilesService) Update(profileId string, profileJson map[string]interface{}) error {
	return svc.profileRepo.Update(profileId, profileJson)
}

func (svc *profilesService) UpdateNodeId(profileId string, nodeId string) error {
	return svc.profileRepo.UpdateNodeId(profileId, nodeId)
}
