package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/domain/schema"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/repository/db"
)

type SchemaService interface {
	Search() (schema.Schemas, resterr.RestErr)
}

type schemaService struct {
	repo db.SchemaRepo
}

func NewSchemaService(repo db.SchemaRepo) SchemaService {
	return &schemaService{
		repo: repo,
	}
}

func (s *schemaService) Search() (schema.Schemas, resterr.RestErr) {
	result, err := s.repo.Search()
	if err != nil {
		return nil, err
	}
	return result, nil
}
