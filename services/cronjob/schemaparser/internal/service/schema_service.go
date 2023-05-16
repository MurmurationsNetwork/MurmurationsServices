package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/redis"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/domain"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/repository/db"
)

type SchemaService interface {
	GetBranchInfo(url string) (*domain.BranchInfo, error)
	HasNewCommit(lastCommit string) (bool, error)
	UpdateSchemas(branchSha string) error
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

func (s *schemaService) GetBranchInfo(url string) (*domain.BranchInfo, error) {
	bytes, err := httputil.GetByteWithBearerToken(url, config.Conf.Github.TOKEN)
	if err != nil {
		return nil, err
	}

	var data domain.BranchInfo
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

func (s *schemaService) UpdateSchemas(branchSha string) error {
	// get the branch root tree
	rootURL := config.Conf.Github.TreeURL + "/" + branchSha
	resp, err := httputil.GetWithBearerToken(rootURL, config.Conf.Github.TOKEN)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	tree := data["tree"].([]interface{})
	var (
		schemaListUrl string
		fieldListUrl  string
	)
	for _, item := range tree {
		itemMap := item.(map[string]interface{})
		if itemMap["path"].(string) == "schemas" {
			schemaListUrl = itemMap["url"].(string)
		} else if itemMap["path"].(string) == "fields" {
			fieldListUrl = itemMap["url"].(string)
		}
	}

	if schemaListUrl == "" {
		return fmt.Errorf("schemas folder not found")
	}

	// Get the fields tree first
	resp, err = httputil.GetWithBearerToken(fieldListUrl, config.Conf.Github.TOKEN)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get fields tree failed, status code: %d", resp.StatusCode)
	}

	var fieldData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&fieldData)
	if err != nil {
		return err
	}

	fieldList := fieldData["tree"].([]interface{})

	fieldListMap := make(map[string]string)
	for _, field := range fieldList {
		fieldMap := field.(map[string]interface{})
		fieldPath := fieldMap["path"].(string)
		fieldUrl := fieldMap["url"].(string)
		if len(fieldPath) > 0 && len(fieldUrl) > 0 {
			fieldListMap[fieldPath] = fieldUrl
		} else {
			fmt.Println("field path or url is empty" + fieldPath + fieldUrl)
			return fmt.Errorf("field path or url is empty" + fieldPath + fieldUrl)
		}
	}

	// get the schemas tree
	resp, err = httputil.GetWithBearerToken(schemaListUrl, config.Conf.Github.TOKEN)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	schemaList := data["tree"].([]interface{})

	for _, schemaName := range schemaList {
		schemaNameMap := schemaName.(map[string]interface{})
		schema, fullJson, err := s.getSchema(schemaNameMap["url"].(string), fieldListMap)
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
		URL:         schema.Metadata.Schema.URL,
		FullSchema:  fullJson,
	}
	err := s.repo.Update(doc)
	if err != nil {
		return err
	}
	return nil
}

func (s *schemaService) getSchema(url string, fieldListMap map[string]string) (*domain.SchemaJSON, map[string]interface{}, error) {
	// Get schema json from GitHub API
	resp, err := httputil.GetWithBearerToken(url, config.Conf.Github.TOKEN)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var decodedContent []byte
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, nil, err
		}
		decodedContent, err = base64.StdEncoding.DecodeString(data["content"].(string))
		if err != nil {
			return nil, nil, err
		}
		if len(decodedContent) == 0 {
			return nil, nil, fmt.Errorf("get schema failed, url: %s, content is empty", url)
		}
	} else {
		return nil, nil, fmt.Errorf("get schema failed, url: %s, status code: %d", url, resp.StatusCode)
	}

	// Get the schema json and full json
	var data domain.SchemaJSON
	err = json.Unmarshal(decodedContent, &data)
	if err != nil {
		return nil, nil, err
	}

	var fullData map[string]interface{}
	err = json.Unmarshal(decodedContent, &fullData)
	if err != nil {
		return nil, nil, err
	}

	// Parse json $ref
	parsedFullData := s.parseProperties(fullData, fieldListMap)

	return &data, parsedFullData, nil
}

func (s *schemaService) parseProperties(fullData map[string]interface{}, fieldListMap map[string]string) map[string]interface{} {
	if fullData["properties"] == nil {
		return fullData
	}
	propertiesMap := fullData["properties"].(map[string]interface{})
	for k, v := range propertiesMap {
		ref := v.(map[string]interface{})
		if ref["type"] == "array" {
			itemsMap := ref["items"].(map[string]interface{})
			parsedSubSchema := s.parseProperties(itemsMap, fieldListMap)
			arrayPropertiesMap := propertiesMap[k].(map[string]interface{})
			arrayPropertiesMap["items"] = parsedSubSchema
			propertiesMap[k] = arrayPropertiesMap
		}
		if ref["$ref"] != nil {
			subSchema, _ := s.schemaParser(ref["$ref"].(string), fieldListMap)
			parsedSubSchema := s.parseProperties(subSchema, fieldListMap)
			propertiesMap[k] = parsedSubSchema
		}
	}
	fullData["properties"] = propertiesMap
	return fullData
}

func (s *schemaService) schemaParser(url string, fieldListMap map[string]string) (map[string]interface{}, error) {
	_, fieldName := path.Split(url)

	fieldListUrl := fieldListMap[fieldName]

	if fieldListUrl == "" {
		return nil, fmt.Errorf("get schema failed, url: %s, fieldListUrl is empty", url)
	}

	resp, err := httputil.GetWithBearerToken(fieldListUrl, config.Conf.Github.TOKEN)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var decodedContent []byte
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, err
		}
		decodedContent, err = base64.StdEncoding.DecodeString(data["content"].(string))
		if err != nil {
			return nil, err
		}
		if len(decodedContent) == 0 {
			return nil, fmt.Errorf("get field failed, url: %s, content is empty", url)
		}
	} else {
		return nil, fmt.Errorf("get field failed, url: %s, status code: %d", url, resp.StatusCode)
	}

	var subSchema map[string]interface{}
	err = json.Unmarshal(decodedContent, &subSchema)
	if err != nil {
		return nil, err
	}

	return subSchema, nil
}
