package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/repository/mongo"
)

// SchemaService defines mtehods for operations on Schemas.
type SchemaService interface {
	Get(schemaName string) (interface{}, error)
	Search() (*model.Schemas, error)
}

type schemaService struct {
	mongoRepo mongo.SchemaRepo
}

// NewSchemaService creates a new SchemaService with the given SchemaRepo.
func NewSchemaService(mongoRepo mongo.SchemaRepo) SchemaService {
	return &schemaService{
		mongoRepo: mongoRepo,
	}
}

// Get fetches a Schema with the given name.
func (s *schemaService) Get(schemaName string) (interface{}, error) {
	result, err := s.mongoRepo.Get(schemaName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Search retrieves all Schemas.
func (s *schemaService) Search() (*model.Schemas, error) {
	result, err := s.mongoRepo.Search()
	if err != nil {
		return nil, err
	}
	return result, nil
}
