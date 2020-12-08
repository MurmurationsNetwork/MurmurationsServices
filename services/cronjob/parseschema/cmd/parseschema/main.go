package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/internal/adapter/mongodb"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/internal/adapter/redisadapter"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/internal/service"
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
		logger.Error("error when trying to get last_commit and schema_list from: "+url, err)
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

	err = svc.UpdateSchemas(dnsInfo.SchemaList)
	if err != nil {
		logger.Error("error when trying to update schemas", err)
		return
	}

	err = svc.SetLastCommit(dnsInfo.LastCommit)
	if err != nil {
		logger.Panic("Error when trying to set schemas:lastCommit", err)
		return
	}
}
