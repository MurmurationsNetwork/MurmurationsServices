package main

import (
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/adapter/mongodb"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/adapter/redisadapter"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/service"
	"time"
)

func init() {
	config.Init()
	mongodb.Init()
}

func main() {
	svc := service.NewSchemaService(db.NewSchemaRepository(), redisadapter.NewClient())

	url := config.Conf.CDN.URL + "/api/schemas"
	dnsInfo, err := svc.GetDNSInfo(url)
	if err != nil {
		logger.Error("Error when trying to get last_commit and schema_list from: "+url, err)
		return
	}

	if dnsInfo.Error != "" {
		logger.Error("Error when trying to get last_commit and schema_list from: "+dnsInfo.Error, err)
		return
	}

	hasNewCommit, err := svc.HasNewCommit(dnsInfo.LastCommit)
	if err != nil {
		logger.Error("Error when trying to get schemas:lastCommit", err)
		return
	}
	if !hasNewCommit {
		return
	}

	// issue-390: delay schema update for 5 minutes.
	now := time.Now()
	commitTime, err := time.Parse(time.RFC3339, dnsInfo.LastCommit)
	if err != nil {
		logger.Error("Error when converting lastCommit to time format", err)
		return
	}
	if now.Before(commitTime.Add(time.Minute * time.Duration(5))) {
		return
	}

	err = svc.UpdateSchemas(dnsInfo.SchemaList)
	if err != nil {
		logger.Error("Error when trying to update schemas", err)
		return
	}

	err = svc.SetLastCommit(dnsInfo.LastCommit)
	if err != nil {
		logger.Panic("Error when trying to set schemas:lastCommit", err)
		return
	}

	logger.Info(fmt.Sprintf("Update Library schemas"))
}
