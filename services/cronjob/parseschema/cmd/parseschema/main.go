package main

import (
	"encoding/json"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/internal/adapter/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type schemaFormat struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Metadata    struct {
		Schema struct {
			Name    string `json:"name"`
			Version int    `json:"version"`
			URL     string `json:"url"`
		} `json:"schema"`
	} `json:"metadata"`
}

type schemaInfo struct {
	Title       string `bson:"title,omitempty"`
	Description string `bson:"description,omitempty"`
	Name        string `bson:"name,omitempty"`
	Version     int    `bson:"version,omitempty"`
	URL         string `bson:"url,omitempty"`
}

func init() {
	config.Init()
	mongodb.Init()
}

func main() {
	schemaListingURL := config.Conf.CDN.URL + "/api/schemas"

	bytes, err := httputil.GetByte(schemaListingURL)
	if err != nil {
		logger.Panic("error when trying to get a list of schemas from: "+schemaListingURL, err)
		return
	}

	var schemas []string
	err = json.Unmarshal(bytes, &schemas)
	if err != nil {
		logger.Panic("error when trying to parse content from: "+schemaListingURL, err)
		return
	}

	for _, schema := range schemas {
		schemaURL := getSchemaURL(schema)

		bytes, err := httputil.GetByte(schemaURL)
		if err != nil {
			logger.Panic("error when trying to get the schema content from: "+schemaURL, err)
			return
		}

		var jsonData schemaFormat
		err = json.Unmarshal(bytes, &jsonData)
		if err != nil {
			logger.Panic("error when trying to parse content from: "+schemaURL, err)
			return
		}

		schemaInfo := schemaInfo{
			Title:       jsonData.Title,
			Description: jsonData.Description,
			Name:        jsonData.Metadata.Schema.Name,
			Version:     jsonData.Metadata.Schema.Version,
			URL:         jsonData.Metadata.Schema.URL,
		}

		filter := bson.M{"name": schemaInfo.Name}
		update := bson.M{"$set": schemaInfo}
		opt := options.FindOneAndUpdate().SetUpsert(true)

		_, err = mongo.Client.FindOneAndUpdate(constant.MongoIndex.Schema, filter, update, opt)
		if err != nil {
			logger.Panic("Error when trying to create a schema record", err)
			return
		}
	}
}

func getSchemaURL(schemaName string) string {
	return config.Conf.CDN.URL + "/schemas/" + schemaName + ".json"
}
