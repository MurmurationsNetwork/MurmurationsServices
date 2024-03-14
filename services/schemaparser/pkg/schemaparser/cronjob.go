package schemaparser

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/redis"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/internal/service"
)

// SchemaCron represents a cron job for managing schema updates.
type SchemaCron struct {
	// svc is an instance of the SchemaService to handle schema related operations.
	svc service.SchemaService
}

// NewCronJob creates a new instance of CronJob and initializes the
// SchemaService.
func NewCronJob() *SchemaCron {
	redisClient := redis.NewClient(config.Values.Redis.URL)
	err := redisClient.Ping()
	if err != nil {
		logger.Error("error when trying to ping Redis", err)
		os.Exit(1)
	}

	return &SchemaCron{
		svc: service.NewSchemaService(
			mongo.NewSchemaRepository(),
			redisClient,
		),
	}
}

// Run start loading the schema.
func (sc *SchemaCron) Run() error {
	if err := sc.connectToMongoDB(); err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if config.Values.IsLocal {
		schemaData, err := readSchemaFilesFromDir("library/schemas")
		if err != nil {
			return fmt.Errorf("failed to read schema files from dir: %w", err)
		}
		fieldData, err := readSchemaFilesFromDir("library/fields")
		if err != nil {
			return fmt.Errorf("failed to read field files from dir: %w", err)
		}

		// if schema data is not empty, update the schemas
		if len(schemaData) > 0 {
			err = sc.svc.UpdateLocalSchemas(schemaData, fieldData)
			if err != nil {
				return fmt.Errorf("failed to update local schemas: %w", err)
			}
		}
	}

	branchInfo, err := sc.svc.GetBranchInfo(config.Values.Github.BranchURL)
	if err != nil {
		return fmt.Errorf(
			"failed to get last_commit and schema_list from %s: %w",
			config.Values.Github.BranchURL,
			err,
		)
	}

	hasNewCommit, err := sc.svc.HasNewCommit(
		branchInfo.Commit.InnerCommit.Author.Date,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to get schemas:lastCommit: %w",
			err,
		)
	}

	if !hasNewCommit {
		logger.Info(
			"No new commit found. Latest commit on GitHub: " +
				branchInfo.Commit.InnerCommit.Author.Date,
		)
		return nil
	}

	err = sc.svc.UpdateSchemas(branchInfo.Commit.Sha)
	if err != nil {
		return fmt.Errorf("failed to update schemas: %w", err)
	}

	// After successfully updating the schemas, update the last commit date.
	err = sc.svc.SetLastCommit(branchInfo.Commit.InnerCommit.Author.Date)
	if err != nil {
		return fmt.Errorf(
			"failed to set schemas:lastCommit: %w",
			err,
		)
	}

	return nil
}

// connectToMongoDB establishes a connection to MongoDB.
func (sc *SchemaCron) connectToMongoDB() error {
	uri := mongodb.GetURI(
		config.Values.Mongo.USERNAME,
		config.Values.Mongo.PASSWORD,
		config.Values.Mongo.HOST,
	)
	err := mongodb.NewClient(uri, config.Values.Mongo.DBName)
	if err != nil {
		return err
	}
	err = mongodb.Client.Ping()
	if err != nil {
		return err
	}
	return nil
}

// get Schemas from local folder.
func readSchemaFilesFromDir(dirPath string) (map[string][]byte, error) {
	filesData := make(
		map[string][]byte,
	) // Initialize a map to store filename and its content

	err := filepath.WalkDir(
		dirPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && filepath.Ext(path) == ".json" {
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				relPath, err := filepath.Rel(dirPath, path)
				if err != nil {
					return err
				}
				filesData[relPath] = data
			}
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return filesData, nil
}
