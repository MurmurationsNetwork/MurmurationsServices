package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/internal/adapter/mongodb"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/internal/adapter/redisadapter"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/internal/domain"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	config.Init()
	mongodb.Init()
}

func main() {
	svc := service.NewSchemaService(redisadapter.NewClient())

	url := config.Conf.CDN.URL + "/api/schemas"
	dnsInfo, err := svc.GetDNSInfo(url)
	if err != nil {
		logger.Panic("error when trying to get last_commit and schema_list from: "+url, err)
		return
	}

	hasNewCommit, err := svc.HasNewCommit(dnsInfo.LastCommit)
	if err != nil {
		logger.Panic("Error when trying to get schemas:lastCommit", err)
		return
	}
	if !hasNewCommit {
		return
	}

	for _, schemaName := range dnsInfo.SchemaList {
		schemaURL := svc.GetSchemaURL(schemaName)
		data, err := svc.GetSchema(schemaURL)
		if err != nil {
			logger.Panic("error when trying to get schema from: "+schemaURL, err)
			return
		}

		doc := domain.Schema{
			Title:       data.Title,
			Description: data.Description,
			Name:        data.Metadata.Schema.Name,
			Version:     data.Metadata.Schema.Version,
			URL:         data.Metadata.Schema.URL,
		}

		filter := bson.M{"name": doc.Name}
		update := bson.M{"$set": doc}
		opt := options.FindOneAndUpdate().SetUpsert(true)

		_, err = mongo.Client.FindOneAndUpdate(constant.MongoIndex.Schema, filter, update, opt)
		if err != nil {
			logger.Panic("Error when trying to create a schema record", err)
			return
		}

		err = svc.SetLastCommit(dnsInfo.LastCommit)
		if err != nil {
			logger.Panic("Error when trying to set schemas:lastCommit", err)
			return
		}
	}
}
