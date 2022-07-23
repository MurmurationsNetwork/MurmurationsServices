package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxycleaner/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxycleaner/internal/repository/db"
)

type ProfilesService interface {
	FindLessThan(timestamp int64) ([]entity.Profile, error)
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

func (svc *profilesService) Delete(profileId string) error {
	return svc.profileRepo.Delete(profileId)
}
