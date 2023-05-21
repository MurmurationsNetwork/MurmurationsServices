package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/redis"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/domain"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/repository/db"
	"github.com/iancoleman/orderedmap"
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
	// Get schema folder list and field folder list
	schemaListUrl, fieldListUrl, err := getBranchFolders(branchSha)

	if err != nil {
		return err
	}

	// Read field folder list and create a map of field name and field url
	fieldListMap, err := getFieldsUrlMap(fieldListUrl)

	if err != nil {
		return err
	}

	// Read schema folder list
	schemaList, err := getGithubTree(schemaListUrl)

	if err != nil {
		return err
	}

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

func (s *schemaService) updateSchema(schema *domain.SchemaJSON, fullJson bson.D) error {
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

func (s *schemaService) getSchema(url string, fieldListMap map[string]string) (*domain.SchemaJSON, bson.D, error) {
	// Get schema json from GitHub API
	schemaJson, err := getGithubFile(url)

	// Get the schema json and full json
	var data domain.SchemaJSON
	err = json.Unmarshal(schemaJson, &data)
	if err != nil {
		return nil, nil, err
	}

	fullData := orderedmap.New()
	err = json.Unmarshal(schemaJson, &fullData)
	if err != nil {
		return nil, nil, err
	}

	// Parse json $ref
	parsedFullData := s.parseProperties(*fullData, fieldListMap)

	return &data, parsedFullData, nil
}

// Direct generate the bson.D type for saving to MongoDB
func (s *schemaService) parseProperties(fullData orderedmap.OrderedMap, fieldListMap map[string]string) bson.D {
	properties, exist := fullData.Get("properties")
	if !exist {
		mapData := bson.D{}
		for _, key := range fullData.Keys() {
			value, _ := fullData.Get(key)

			// If the value can be parsed into OrderedMap, it means it is a map, we need to parse it recursively
			vMap, ok := value.(orderedmap.OrderedMap)
			if !ok {
				mapData = append(mapData, bson.E{Key: key, Value: value})
			} else {
				mapData = append(mapData, bson.E{Key: key, Value: s.parseProperties(vMap, fieldListMap)})
			}
		}

		return mapData
	}

	propertiesMap := properties.(orderedmap.OrderedMap)
	for _, k := range propertiesMap.Keys() {
		v, _ := propertiesMap.Get(k)
		ref := v.(orderedmap.OrderedMap)

		// If $ref exists, parse the ref
		refPath, ok := ref.Get("$ref")
		if ok && refPath != nil {
			subSchema, _ := s.schemaParser(refPath.(string), fieldListMap)
			parsedSubSchema := s.parseProperties(*subSchema, fieldListMap)
			propertiesMap.Set(k, parsedSubSchema)
			continue
		}

		refType, ok := ref.Get("type")
		// If type is array, parse the items
		// Otherwise, directly parse the orderedmap to bson.D
		if ok && refType == "array" {
			refItems, _ := ref.Get("items")
			itemsMap := refItems.(orderedmap.OrderedMap)
			parsedSubSchema := s.parseProperties(itemsMap, fieldListMap)

			// In the array, items need to be parsed recursively
			arrayPropertiesMap := bson.D{}
			for _, key := range ref.Keys() {
				if key == "items" {
					arrayPropertiesMap = append(arrayPropertiesMap, bson.E{Key: key, Value: parsedSubSchema})
				} else {
					value, _ := ref.Get(key)
					arrayPropertiesMap = append(arrayPropertiesMap, bson.E{Key: key, Value: value})
				}
			}
			propertiesMap.Set(k, arrayPropertiesMap)
		} else {
			refMap := generateBsonD(ref)
			propertiesMap.Set(k, refMap)
		}
	}

	propertiesData := generateBsonD(propertiesMap)
	fullData.Set("properties", propertiesData)

	returnData := generateBsonD(fullData)
	return returnData
}

func (s *schemaService) schemaParser(url string, fieldListMap map[string]string) (*orderedmap.OrderedMap, error) {
	_, fieldName := path.Split(url)

	fieldListUrl := fieldListMap[fieldName]

	if fieldListUrl == "" {
		return nil, fmt.Errorf("get schema failed, url: %s, fieldListUrl is empty", url)
	}

	// Get field json from GitHub API
	fieldJson, err := getGithubFile(fieldListUrl)

	subSchema := orderedmap.New()
	err = json.Unmarshal(fieldJson, &subSchema)
	if err != nil {
		return nil, err
	}

	return subSchema, nil
}

func getBranchFolders(branchSha string) (string, string, error) {
	rootURL := config.Conf.Github.TreeURL + "/" + branchSha

	rootList, err := getGithubTree(rootURL)

	if err != nil {
		return "", "", err
	}

	var (
		schemaListUrl string
		fieldListUrl  string
	)
	for _, item := range rootList {
		itemMap := item.(map[string]interface{})
		if itemMap["path"].(string) == "schemas" {
			schemaListUrl = itemMap["url"].(string)
		} else if itemMap["path"].(string) == "fields" {
			fieldListUrl = itemMap["url"].(string)
		}
	}

	if schemaListUrl == "" {
		return "", "", fmt.Errorf("schemas folder not found")
	}
	if fieldListUrl == "" {
		return "", "", fmt.Errorf("fields folder not found")
	}

	return schemaListUrl, fieldListUrl, nil
}

func getFieldsUrlMap(fieldListUrl string) (map[string]string, error) {
	fieldList, err := getGithubTree(fieldListUrl)
	if err != nil {
		return nil, err
	}

	fieldListMap := make(map[string]string)
	for _, field := range fieldList {
		fieldMap := field.(map[string]interface{})
		fieldPath := fieldMap["path"].(string)
		fieldUrl := fieldMap["url"].(string)
		if len(fieldPath) > 0 && len(fieldUrl) > 0 {
			fieldListMap[fieldPath] = fieldUrl
		} else {
			return nil, fmt.Errorf("field path or url is empty" + fieldPath + fieldUrl)
		}
	}

	return fieldListMap, nil
}

func getGithubTree(url string) ([]interface{}, error) {
	resp, err := httputil.GetWithBearerToken(url, config.Conf.Github.TOKEN)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	tree := data["tree"].([]interface{})
	return tree, nil
}

func getGithubFile(url string) ([]byte, error) {
	resp, err := httputil.GetWithBearerToken(url, config.Conf.Github.TOKEN)
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
			return nil, fmt.Errorf("get file failed, url: %s, content is empty", url)
		}
	} else {
		return nil, fmt.Errorf("get file failed, url: %s, status code: %d", url, resp.StatusCode)
	}

	return decodedContent, nil
}

func generateBsonD(orderedMap orderedmap.OrderedMap) bson.D {
	bsonData := bson.D{}
	for _, key := range orderedMap.Keys() {
		value, _ := orderedMap.Get(key)
		bsonData = append(bsonData, bson.E{Key: key, Value: value})
	}
	return bsonData
}
