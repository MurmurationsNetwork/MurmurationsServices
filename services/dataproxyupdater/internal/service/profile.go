package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyupdater/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyupdater/internal/repository/mongo"
)

type ProfilesService interface {
	Count(profileID string) (int64, error)
	Add(profileJSON map[string]interface{}) error
	Update(
		schemaName string,
		profileJSON map[string]interface{},
	) (map[string]interface{}, error)
	UpdateNodeID(profileID string, nodeID string) error
	GetNotPosted() ([]model.Profile, error)
	UpdateIsPosted(nodeID string) error
	Delete(profileID string) error
}

type profilesService struct {
	mongoRepo mongo.ProfileRepository
}

func NewProfileService(profileRepo mongo.ProfileRepository) ProfilesService {
	return &profilesService{
		mongoRepo: profileRepo,
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

func (svc *profilesService) GetNotPosted() ([]model.Profile, error) {
	return svc.mongoRepo.GetNotPosted()
}

func (svc *profilesService) UpdateIsPosted(nodeID string) error {
	return svc.mongoRepo.UpdateIsPosted(nodeID)
}

func (svc *profilesService) Delete(profileID string) error {
	return svc.mongoRepo.Delete(profileID)
}
