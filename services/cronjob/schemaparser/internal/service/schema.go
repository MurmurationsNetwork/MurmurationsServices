package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/redis"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/schemaparser"
)

const (
	LastCommitKey = "schemas:lastCommit"
	maxGoroutines = 10
)

type SchemaService interface {
	GetBranchInfo(url string) (*model.BranchInfo, error)
	HasNewCommit(lastCommit string) (bool, error)
	UpdateSchemas(branchSha string) error
	SetLastCommit(lastCommit string) error
}

type schemaService struct {
	mongoRepo mongo.SchemaRepository
	redis     redis.Redis
}

func NewSchemaService(
	mongoRepo mongo.SchemaRepository,
	redis redis.Redis,
) SchemaService {
	return &schemaService{
		mongoRepo: mongoRepo,
		redis:     redis,
	}
}

// GetBranchInfo retrieves information about a specific GitHub branch given its URL.
func (s *schemaService) GetBranchInfo(url string) (*model.BranchInfo, error) {
	// GetByteWithBearerToken sends a GET request to the provided URL, using
	// the provided bearer token for authentication.
	bytes, err := httputil.GetByteWithBearerToken(
		url,
		config.Values.Github.TOKEN,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from %s: %w", url, err)
	}

	var data model.BranchInfo

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &data, nil
}

// HasNewCommit checks if there's a new commit in the schema repository.
func (s *schemaService) HasNewCommit(lastCommit string) (bool, error) {
	val, err := s.redis.Get(LastCommitKey)
	if err != nil {
		return false, fmt.Errorf(
			"failed to retrieve last commit from Redis: %w",
			err,
		)
	}

	return val != lastCommit, nil
}

func (s *schemaService) UpdateSchemas(branchSha string) error {
	// Get schema folder list and field folder list
	schemaListURL, fieldListURL, err := getSchemaAndFieldFolderURLs(branchSha)
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

	// Create a semaphore channel to limit goroutines.
	semaphore := make(chan struct{}, maxGoroutines)

	g, ctx := errgroup.WithContext(context.Background())

	for _, schemaName := range schemaList {
		schemaNameMap := schemaName.(map[string]interface{})
		url := schemaNameMap["url"].(string)

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				parser := schemaparser.NewSchemaParser(fieldListMap)
				result, err := parser.GetSchema(url)
				if err != nil {
					return err
				}
				return s.updateSchema(result.Schema, result.FullJSON)
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
	schema *model.SchemaJSON,
	fullJSON bson.D,
) error {
	doc := &model.Schema{
		Title:       schema.Title,
		Description: schema.Description,
		Name:        schema.Metadata.Schema.Name,
		URL:         schema.Metadata.Schema.URL,
		FullSchema:  fullJSON,
	}
	err := s.mongoRepo.Update(doc)
	if err != nil {
		return err
	}
	return nil
}

// getBranchFolders retrieves URLs for 'schemas' and 'fields' directories
// given the branch's SHA.
func getSchemaAndFieldFolderURLs(branchSha string) (string, string, error) {
	// Construct the URL to the GitHub tree API endpoint for the branch.
	rootURL := config.Values.Github.TreeURL + "/" + branchSha

	// Fetch the tree (list of files and directories) from the GitHub API.
	rootList, err := getGithubTree(rootURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch tree from %s: %w",
			rootURL, err)
	}

	var (
		schemaListURL string
		fieldListURL  string
	)

	// Iterate over each item in the root directory.
	for _, item := range rootList {
		itemMap := item.(map[string]interface{})
		// If the item is the 'schemas' directory, save its URL.
		if itemMap["path"].(string) == "schemas" {
			schemaListURL = itemMap["url"].(string)
		}
		// If the item is the 'fields' directory, save its URL.
		if itemMap["path"].(string) == "fields" {
			fieldListURL = itemMap["url"].(string)
		}
	}

	if schemaListURL == "" {
		return "", "", fmt.Errorf("'schemas' directory not found")
	}
	if fieldListURL == "" {
		return "", "", fmt.Errorf("'fields' directory not found")
	}

	return schemaListURL, fieldListURL, nil
}

// getFieldsURLMap fetches the GitHub tree (list of files and directories)
// from the given URL, and then creates a map that maps field path to field URL.
func getFieldsURLMap(fieldListURL string) (map[string]string, error) {
	fieldList, err := getGithubTree(fieldListURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tree from %s: %w",
			fieldListURL, err)
	}

	// Initialize a map to hold the field path to URL mapping.
	fieldListMap := make(map[string]string)

	// Iterate over each field in the list.
	for _, field := range fieldList {
		fieldMap := field.(map[string]interface{})
		fieldPath := fieldMap["path"].(string)
		fieldURL := fieldMap["url"].(string)

		// If the field has a non-empty path and URL, add it to the map.
		if len(fieldPath) > 0 && len(fieldURL) > 0 {
			fieldListMap[fieldPath] = fieldURL
		} else {
			return nil, fmt.Errorf("field path '%s' or url '%s' is empty",
				fieldPath, fieldURL)
		}
	}

	// Return the map.
	return fieldListMap, nil
}

// getGithubTree fetches and returns the GitHub repository's tree (list of
// files and directories) given a URL to the tree API endpoint.
//
// https://docs.github.com/en/rest/git/trees?apiVersion=2022-11-28#get-a-tree.
func getGithubTree(url string) ([]interface{}, error) {
	resp, err := httputil.GetWithBearerToken(url, config.Values.Github.TOKEN)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}

	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	tree := data["tree"].([]interface{})

	return tree, nil
}
