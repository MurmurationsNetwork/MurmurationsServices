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

type schema struct {
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

type schemaDoc struct {
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
	url := config.Conf.CDN.URL + "/api/schemas"

	schemaList, err := getSchemaList(url)
	if err != nil {
		logger.Panic("error when trying to get schemaList from: "+url, err)
		return
	}

	for _, schemaName := range schemaList {
		schemaURL := getSchemaURL(schemaName)

		data, err := getSchema(schemaURL)
		if err != nil {
			logger.Panic("error when trying to get schema from: "+schemaURL, err)
			return
		}

		doc := schemaDoc{
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
	}
}

func getSchemaList(url string) ([]string, error) {
	bytes, err := httputil.GetByte(url)
	if err != nil {
		return nil, err
	}

	type jsonFormat struct {
		LastCommit string   `json:"last_commit"`
		SchemaList []string `json:"schema_list"`
	}

	var data jsonFormat
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data.SchemaList, nil
}

func getSchemaURL(schemaName string) string {
	return config.Conf.CDN.URL + "/schemas/" + schemaName + ".json"
}

func getSchema(url string) (*schema, error) {
	bytes, err := httputil.GetByte(url)
	if err != nil {
		return nil, err
	}

	var data schema
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
