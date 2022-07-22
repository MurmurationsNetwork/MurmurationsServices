package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/repository/db"
)

type MappingsService interface {
	Get(schemaName string) map[string]string
}

type mappingsService struct {
	mappingRepo db.MappingRepository
}

func NewMappingService(mappingRepo db.MappingRepository) MappingsService {
	return &mappingsService{
		mappingRepo: mappingRepo,
	}
}

func (svc *mappingsService) Get(schemaName string) map[string]string {
	schemaRaw := svc.mappingRepo.Get(schemaName)

	schema := make(map[string]string)
	// remove id and __v
	for k, v := range schemaRaw {
		if k == "_id" || k == "__v" {
			continue
		}
		schema[k] = v.(string)
	}

	return schema
}
