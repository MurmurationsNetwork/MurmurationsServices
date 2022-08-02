package service

import (
	"encoding/json"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/redis"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/domain"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/repository/db"
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
		schema, fullJson, err := s.getSchema(schemaURL)
		if err != nil {
			return err
		}

		err = s.updateSchema(schema, fullJson)
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

func (s *schemaService) updateSchema(schema *domain.SchemaJSON, fullJson map[string]interface{}) error {
	doc := &domain.Schema{
		Title:       schema.Title,
		Description: schema.Description,
		Name:        schema.Metadata.Schema.Name,
		Version:     schema.Metadata.Schema.Version,
		URL:         schema.Metadata.Schema.URL,
		FullSchema:  fullJson,
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

func (s *schemaService) getSchema(url string) (*domain.SchemaJSON, map[string]interface{}, error) {
	bytes, err := httputil.GetByte(url)
	if err != nil {
		return nil, nil, err
	}

	var data domain.SchemaJSON
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, nil, err
	}

	var fullData map[string]interface{}
	err = json.Unmarshal(bytes, &fullData)
	if err != nil {
		return nil, nil, err
	}

	// parse json $ref
	parsedFullData := s.parseProperties(fullData)

	return &data, parsedFullData, nil
}

func (s *schemaService) parseProperties(fullData map[string]interface{}) map[string]interface{} {
	if fullData["properties"] == nil {
		return fullData
	}
	propertiesMap := fullData["properties"].(map[string]interface{})
	for k, v := range propertiesMap {
		ref := v.(map[string]interface{})
		if ref["type"] == "array" {
			itemsMap := ref["items"].(map[string]interface{})
			parsedSubSchema := s.parseProperties(itemsMap)
			arrayPropertiesMap := propertiesMap[k].(map[string]interface{})
			arrayPropertiesMap["items"] = parsedSubSchema
			propertiesMap[k] = arrayPropertiesMap
		}
		if ref["$ref"] != nil {
			subSchema, _ := s.schemaParser(ref["$ref"].(string))
			parsedSubSchema := s.parseProperties(subSchema)
			propertiesMap[k] = parsedSubSchema
		}
	}
	fullData["properties"] = propertiesMap
	return fullData
}

func (s *schemaService) schemaParser(url string) (map[string]interface{}, error) {
	// remove ".."
	url = url[2:]
	bytes, err := httputil.GetByte(config.Conf.CDN.URL + url)
	if err != nil {
		return nil, err
	}
	var subSchema map[string]interface{}
	err = json.Unmarshal(bytes, &subSchema)
	if err != nil {
		return nil, err
	}

	return subSchema, nil
}
