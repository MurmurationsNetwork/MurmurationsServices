package service

import (
	"encoding/json"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/redis"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/internal/domain"
)

type SchemaService interface {
	GetDNSInfo(url string) (*domain.DnsInfo, error)
	HasNewCommit(lastCommit string) (bool, error)
	SetLastCommit(lastCommit string) error
	GetSchemaURL(schemaName string) string
	GetSchema(url string) (*domain.SchemaJSON, error)
}

type schemaService struct {
	redis redis.Redis
}

func NewSchemaService(redis redis.Redis) SchemaService {
	return &schemaService{
		redis: redis,
	}
}

func (s *schemaService) GetDNSInfo(url string) (*domain.DnsInfo, error) {
	bytes, err := httputil.GetByte(url)
	if err != nil {
		return nil, err
	}

	var data domain.DnsInfo
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *schemaService) HasNewCommit(lastCommit string) (bool, error) {
	val, err := s.redis.Get("schemas:lastCommit")
	if err != nil {
		return false, err
	}
	return val != lastCommit, nil
}

func (s *schemaService) SetLastCommit(lastCommit string) error {
	err := s.redis.Set("schemas:lastCommit", lastCommit, 0)
	if err != nil {
		return err
	}
	return nil
}

func (s *schemaService) GetSchemaURL(schemaName string) string {
	return config.Conf.CDN.URL + "/schemas/" + schemaName + ".json"
}

func (s *schemaService) GetSchema(url string) (*domain.SchemaJSON, error) {
	bytes, err := httputil.GetByte(url)
	if err != nil {
		return nil, err
	}

	var data domain.SchemaJSON
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
