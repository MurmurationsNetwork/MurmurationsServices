package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/internal/repository/db"
)

type ProfilesService interface {
	FindLessThan(timestamp int64) ([]entity.Profile, error)
	UpdateAccessTime(profileId string) error
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

func (svc *profilesService) FindLessThan(timestamp int64) ([]entity.Profile, error) {
	return svc.profileRepo.FindLessThan(timestamp)
}

func (svc *profilesService) UpdateAccessTime(profileId string) error {
	return svc.profileRepo.UpdateAccessTime(profileId)
}

func (svc *profilesService) Delete(profileId string) error {
	return svc.profileRepo.Delete(profileId)
}
