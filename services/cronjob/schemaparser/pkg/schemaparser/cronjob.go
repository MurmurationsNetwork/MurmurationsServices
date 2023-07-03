package schemaparser

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/redis"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/service"
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
		logger.Panic("error when trying to ping Redis", err)
		return nil
	}

	return &SchemaCron{
		svc: service.NewSchemaService(
			mongo.NewSchemaRepository(),
			redisClient,
		),
	}
}

// Run executes the cron job. It fetches the latest branch info, checks for new
// commits, updates schemas if a new commit is found, and sets the last commit.
func (sc *SchemaCron) Run() {
	if err := sc.connectToMongoDB(); err != nil {
		logger.Panic("error when trying to connect to MongoDB", err)
		return
	}

	url := config.Values.Github.BranchURL

	// Use the SchemaService to get branch info.
	branchInfo, err := sc.svc.GetBranchInfo(url)
	if err != nil {
		logger.Error(
			"Error when trying to get last_commit and schema_list from: "+url,
			err,
		)
		return
	}

	// Check if the latest commit date indicates a new commit.
	hasNewCommit, err := sc.svc.HasNewCommit(
		branchInfo.Commit.InnerCommit.Author.Date,
	)
	if err != nil {
		logger.Error("Error when trying to get schemas:lastCommit", err)
		return
	}

	// If there's no new commit, there's nothing to do.
	if !hasNewCommit {
		logger.Info(
			"No new commit found. Latest commit on GitHub: " + branchInfo.Commit.InnerCommit.Author.Date,
		)
		return
	}

	// If there is a new commit, update the schemas.
	err = sc.svc.UpdateSchemas(branchInfo.Commit.Sha)
	if err != nil {
		logger.Error("Error when trying to update schemas", err)
		return
	}

	// After successfully updating the schemas, update the last commit date.
	err = sc.svc.SetLastCommit(branchInfo.Commit.InnerCommit.Author.Date)
	if err != nil {
		logger.Panic("Error when trying to set schemas:lastCommit", err)
	}
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
