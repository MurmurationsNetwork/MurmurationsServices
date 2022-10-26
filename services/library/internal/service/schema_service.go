package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/domain/schema"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/repository/db"
)

type SchemaService interface {
	Get(schemaName string) (interface{}, []jsonapi.Error)
	Search() (schema.Schemas, []jsonapi.Error)
}

type schemaService struct {
	repo db.SchemaRepo
}

func NewSchemaService(repo db.SchemaRepo) SchemaService {
	return &schemaService{
		repo: repo,
	}
}

func (s *schemaService) Get(schemaName string) (interface{}, []jsonapi.Error) {
	result, err := s.repo.Get(schemaName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *schemaService) Search() (schema.Schemas, []jsonapi.Error) {
	result, err := s.repo.Search()
	if err != nil {
		return nil, err
	}
	return result, nil
}
