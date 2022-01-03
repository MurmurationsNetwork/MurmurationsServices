package service

import (
	"encoding/json"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/redis"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/internal/domain"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/internal/repository/db"
)

type SchemaService interface {
	GetDNSInfo(url string) (*domain.DnsInfo, error)
	HasNewCommit(lastCommit string) (bool, error)
	UpdateSchemas(schemaList []string) error
	SetLastCommit(lastCommit string) error
}

type schemaService struct {
	repo  db.SchemaRepository
	redis redis.Redis
}

func NewSchemaService(repo db.SchemaRepository, redis redis.Redis) SchemaService {
	return &schemaService{
		repo:  repo,
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

func (s *schemaService) UpdateSchemas(schemaList []string) error {
	for _, schemaName := range schemaList {
		schemaURL := s.getSchemaURL(schemaName)
		json, err := s.getSchema(schemaURL)
		if err != nil {
			return err
		}

		err = s.updateSchema(json)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *schemaService) SetLastCommit(newLastCommitTime string) error {
	oldLastCommitTime, err := s.redis.Get("schemas:lastCommit")
	if err != nil {
		return err
	}

	ok, err := shouldSetLastCommitTime(oldLastCommitTime, newLastCommitTime)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	err = s.redis.Set("schemas:lastCommit", newLastCommitTime, 0)
	if err != nil {
		return err
	}

	return nil
}

func shouldSetLastCommitTime(oldTime, newTime string) (bool, error) {
	if oldTime == "" {
		return true, nil
	}
	if newTime == "" {
		return false, nil
	}

	t1, err := time.Parse(time.RFC3339, oldTime)
	if err != nil {
		return false, err
	}
	t2, err := time.Parse(time.RFC3339, newTime)
	if err != nil {
		return false, err
	}

	// To make sure DNS updates the content of schemas.
	// The system won't update the last commit time until it passes a certain period.
	if int(t2.Sub(t1).Minutes()) < 10 {
		return false, nil
	}

	return true, nil
}

func (s *schemaService) updateSchema(json *domain.SchemaJSON) error {
	doc := &domain.Schema{
		Title:       json.Title,
		Description: json.Description,
		Name:        json.Metadata.Schema.Name,
		Version:     json.Metadata.Schema.Version,
		URL:         json.Metadata.Schema.URL,
	}
	err := s.repo.Update(doc)
	if err != nil {
		return err
	}
	return nil
}

func (s *schemaService) getSchemaURL(schemaName string) string {
	return config.Conf.CDN.URL + "/schemas/" + schemaName + ".json"
}

func (s *schemaService) getSchema(url string) (*domain.SchemaJSON, error) {
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
