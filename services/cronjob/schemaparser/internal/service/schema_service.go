package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/iancoleman/orderedmap"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/redis"
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

func NewSchemaService(
	repo db.SchemaRepository,
	redis redis.Redis,
) SchemaService {
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
	schemaListURL, fieldListURL, err := getBranchFolders(branchSha)
	if err != nil {
		return err
	}

	// Read field folder list and create a map of field name and field URL
	fieldListMap, err := getFieldsURLMap(fieldListURL)
	if err != nil {
		return err
	}

	// Read schema folder list
	schemaList, err := getGithubTree(schemaListURL)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(context.Background())

	for _, schemaName := range schemaList {
		schemaNameMap := schemaName.(map[string]interface{})
		url := schemaNameMap["url"].(string)
		// Create a new goroutine to get and update each schema.
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				schema, fullJSON, err := s.getSchema(url, fieldListMap)
				if err != nil {
					return err
				}
				return s.updateSchema(schema, fullJSON)
			}
		})
	}

	return g.Wait()
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

func (s *schemaService) updateSchema(
	schema *domain.SchemaJSON,
	fullJSON bson.D,
) error {
	doc := &domain.Schema{
		Title:       schema.Title,
		Description: schema.Description,
		Name:        schema.Metadata.Schema.Name,
		URL:         schema.Metadata.Schema.URL,
		FullSchema:  fullJSON,
	}
	err := s.repo.Update(doc)
	if err != nil {
		return err
	}
	return nil
}

func (s *schemaService) getSchema(
	url string,
	fieldListMap map[string]string,
) (*domain.SchemaJSON, bson.D, error) {
	// Get schema json from GitHub API
	schemaJSON, err := getGithubFile(url)
	if err != nil {
		return nil, nil, err
	}

	// Get the schema json and full json
	var data domain.SchemaJSON
	err = json.Unmarshal(schemaJSON, &data)
	if err != nil {
		return nil, nil, err
	}

	fullData := orderedmap.New()
	err = json.Unmarshal(schemaJSON, &fullData)
	if err != nil {
		return nil, nil, err
	}

	// Parse json $ref
	parsedFullData := s.parseProperties(*fullData, fieldListMap)

	return &data, parsedFullData, nil
}

// Direct generate the bson.D type for saving to MongoDB.
func (s *schemaService) parseProperties(
	fullData orderedmap.OrderedMap,
	fieldListMap map[string]string,
) bson.D {
	properties, exist := fullData.Get("properties")
	if !exist {
		propertiesData := s.generateBsonDObject(fullData, fieldListMap)
		return propertiesData
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
			arrayPropertiesMap := s.generateBsonDArray(ref, fieldListMap)
			propertiesMap.Set(k, arrayPropertiesMap)
		} else if ok && refType == "object" {
			objectPropertiesMap := s.generateBsonDObject(ref, fieldListMap)
			propertiesMap.Set(k, objectPropertiesMap)
		} else {
			refMap := generateBsonD(ref)
			propertiesMap.Set(k, refMap)
		}
	}

	propertiesData := generateBsonD(propertiesMap)
	fullData.Set("properties", propertiesData)

	returnData := s.generateBsonDObject(fullData, fieldListMap)
	return returnData
}

func (s *schemaService) schemaParser(
	url string,
	fieldListMap map[string]string,
) (*orderedmap.OrderedMap, error) {
	_, fieldName := path.Split(url)

	fieldListURL := fieldListMap[fieldName]

	if fieldListURL == "" {
		return nil, fmt.Errorf(
			"get schema failed, url: %s, fieldListURL is empty",
			url,
		)
	}

	// Get field json from GitHub API
	fieldJSON, err := getGithubFile(fieldListURL)
	if err != nil {
		return nil, err
	}

	subSchema := orderedmap.New()
	err = json.Unmarshal(fieldJSON, &subSchema)
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
		schemaListURL string
		fieldListURL  string
	)
	for _, item := range rootList {
		itemMap := item.(map[string]interface{})
		if itemMap["path"].(string) == "schemas" {
			schemaListURL = itemMap["url"].(string)
		} else if itemMap["path"].(string) == "fields" {
			fieldListURL = itemMap["url"].(string)
		}
	}

	if schemaListURL == "" {
		return "", "", fmt.Errorf("schemas folder not found")
	}
	if fieldListURL == "" {
		return "", "", fmt.Errorf("fields folder not found")
	}

	return schemaListURL, fieldListURL, nil
}

func getFieldsURLMap(fieldListURL string) (map[string]string, error) {
	fieldList, err := getGithubTree(fieldListURL)
	if err != nil {
		return nil, err
	}

	fieldListMap := make(map[string]string)
	for _, field := range fieldList {
		fieldMap := field.(map[string]interface{})
		fieldPath := fieldMap["path"].(string)
		fieldURL := fieldMap["url"].(string)
		if len(fieldPath) > 0 && len(fieldURL) > 0 {
			fieldListMap[fieldPath] = fieldURL
		} else {
			return nil, fmt.Errorf("field path or url is empty" + fieldPath + fieldURL)
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
		decodedContent, err = base64.StdEncoding.DecodeString(
			data["content"].(string),
		)
		if err != nil {
			return nil, err
		}
		if len(decodedContent) == 0 {
			return nil, fmt.Errorf(
				"get file failed, url: %s, content is empty",
				url,
			)
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

func (s *schemaService) generateBsonDArray(
	orderedMap orderedmap.OrderedMap,
	fieldListMap map[string]string,
) bson.D {
	items, _ := orderedMap.Get("items")
	itemsMap := items.(orderedmap.OrderedMap)
	parsedSubSchema := s.parseProperties(itemsMap, fieldListMap)

	// In the array, items need to be parsed recursively
	bsonData := bson.D{}
	for _, key := range orderedMap.Keys() {
		if key == "items" {
			bsonData = append(
				bsonData,
				bson.E{Key: key, Value: parsedSubSchema},
			)
		} else {
			value, _ := orderedMap.Get(key)
			bsonData = append(bsonData, bson.E{Key: key, Value: value})
		}
	}

	return bsonData
}

func (s *schemaService) generateBsonDObject(
	orderedMap orderedmap.OrderedMap,
	fieldListMap map[string]string,
) bson.D {
	bsonData := bson.D{}
	for _, key := range orderedMap.Keys() {
		value, _ := orderedMap.Get(key)

		// If the value can be parsed into OrderedMap, it means it is a map, we need to parse it recursively
		vMap, ok := value.(orderedmap.OrderedMap)
		if !ok {
			bsonData = append(bsonData, bson.E{Key: key, Value: value})
		} else {
			bsonData = append(bsonData, bson.E{
				Key:   key,
				Value: s.parseProperties(vMap, fieldListMap),
			})
		}
	}
	return bsonData
}
