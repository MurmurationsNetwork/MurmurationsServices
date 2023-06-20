package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/repository/db"
)

// SchemaService defines mtehods for operations on Schemas.
type SchemaService interface {
	Get(schemaName string) (interface{}, error)
	Search() (*model.Schemas, error)
}

type schemaService struct {
	repo db.SchemaRepo
}

// NewSchemaService creates a new SchemaService with the given SchemaRepo.
func NewSchemaService(repo db.SchemaRepo) SchemaService {
	return &schemaService{
		repo: repo,
	}
}

// Get fetches a Schema with the given name.
func (s *schemaService) Get(schemaName string) (interface{}, error) {
	result, err := s.repo.Get(schemaName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Search retrieves all Schemas.
func (s *schemaService) Search() (*model.Schemas, error) {
	result, err := s.repo.Search()
	if err != nil {
		return nil, err
	}
	return result, nil
}
