package schemaparser

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/adapter/mongodb"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/adapter/redisadapter"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/service"
)

func init() {
	config.Init()
	mongodb.Init()
}

// SchemaCron represents a cron job for managing schema updates.
type SchemaCron struct {
	// svc is an instance of the SchemaService to handle schema related operations.
	svc service.SchemaService
}

// NewCronJob creates a new instance of CronJob and initializes the
// SchemaService.
func NewCronJob() *SchemaCron {
	return &SchemaCron{
		svc: service.NewSchemaService(
			db.NewSchemaRepository(),
			redisadapter.NewClient(),
		),
	}
}

// Run executes the cron job. It fetches the latest branch info, checks for new
// commits, updates schemas if a new commit is found, and sets the last commit.
func (sc *SchemaCron) Run() {
	url := config.Conf.Github.BranchURL

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
