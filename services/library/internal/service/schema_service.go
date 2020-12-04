package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/domain/schema"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/repository/schemarepo"
)

var (
	SchemaService schemaServiceInterface = &schemaService{}
)

type schemaServiceInterface interface {
	Search() (schema.Schemas, resterr.RestErr)
}

type schemaService struct{}

func (s *schemaService) Search() (schema.Schemas, resterr.RestErr) {
	result, err := schemarepo.Schema.Search()
	if err != nil {
		return nil, err
	}
	return result, nil
}
