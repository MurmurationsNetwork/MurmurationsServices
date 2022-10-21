package importutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var KvmCategory = map[string]string{
	"2cd00bebec0c48ba9db761da48678134": "#non-profit",
	"77b3c33a92554bcf8e8c2c86cedd6f6f": "#commercial",
	"c2dc278a2d6a4b9b8a50cb606fc017ed": "#event",
}

type Node struct {
	NodeId         string   `json:"node_id,omitempty"`
	ProfileUrl     string   `json:"profile_url,omitempty"`
	Status         string   `json:"status,omitempty"`
	FailureReasons []string `json:"failure_reasons,omitempty"`
}

type NodeData struct {
	Data   Node
	Status int `json:"status,omitempty"`
}

func GetMapping(schemaName string) (map[string]string, error) {
	filter := bson.M{"schema": schemaName}
	result := mongo.Client.FindOne(constant.MongoIndex.Mapping, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("could not find mapping for schema %s", schemaName)
		}
		return nil, fmt.Errorf("error when trying to find the mapping, error message: %s", result.Err())
	}
	schemaRaw := make(map[string]interface{})
	err := result.Decode(schemaRaw)
	if err != nil {
		return nil, fmt.Errorf("error when trying to parse database response, error message: %s", err)
	}

	// remove id and __v
	schema := make(map[string]string)
	for i, v := range schemaRaw {
		if i == "__v" || i == "_id" {
			continue
		}
		schema[i] = v.(string)
	}
	return schema, nil
}

func Hash(doc string) (string, error) {
	// ref: https://stackoverflow.com/questions/55256365/how-to-obtain-same-hash-from-json
	var v interface{}
	err := json.Unmarshal([]byte(doc), &v)
	if err != nil {
		return "", err
	}
	hashDoc, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(hashDoc)
	return hex.EncodeToString(sum[0:]), nil
}

func MapFieldsName(profile map[string]interface{}, mapping map[string]string) map[string]interface{} {
	profileJson := make(map[string]interface{})

	for k, v := range mapping {
		if profile[v] == nil || profile[v] == "" {
			continue
		}
		// todo: temporary fix latitude & longitude, will delete them in the future
		if k == "latitude" || k == "longitude" {
			// set precision to 8
			precision := math.Pow(10, float64(8))
			truncatedValue := math.Round(profile[v].(float64)*precision) / precision
			profileJson[k] = truncatedValue
			continue
		}
		// todo: temporary fix extra space in the string, might delete them in the future
		if k != "tags" && k != "kvm_category" {
			profileJson[k] = strings.TrimSpace(profile[v].(string))
			continue
		}
		profileJson[k] = profile[v]
	}

	return profileJson
}

func MapProfile(profile map[string]interface{}, mapping map[string]string, schema string) (map[string]interface{}, error) {
	// change field name
	profileJson := MapFieldsName(profile, mapping)

	// hash the updated data
	doc, err := json.Marshal(profileJson)
	if err != nil {
		return nil, err
	}
	profileHash, err := Hash(string(doc))
	if err != nil {
		return nil, err
	}

	// hash
	profileJson["source_data_hash"] = profileHash

	// oid
	profileJson["oid"] = profile["id"]
	// schema
	s := []string{schema}
	profileJson["linked_schemas"] = s
	// metadata
	// todo: seeder's access_time might need to set as "1646697600"
	metadata := map[string]interface{}{
		"sources": []map[string]interface{}{
			{
				"name":             "Karte von Morgen / Map of Tomorrow",
				"url":              "https://kartevonmorgen.org",
				"profile_data_url": "https://api.ofdb.io/v0/entries/" + profileJson["oid"].(string),
				"access_time":      time.Now().Unix(),
			},
		},
	}
	profileJson["metadata"] = metadata

	// replace kvm_category with real name
	if profileJson["kvm_category"] != nil {
		categoriesInterface := profileJson["kvm_category"].([]interface{})
		categoriesString := make([]string, len(categoriesInterface))
		for i, v := range categoriesInterface {
			categoriesString[i] = KvmCategory[v.(string)]
		}
		profileJson["kvm_category"] = categoriesString
	}

	return profileJson, nil
}

func Validate(validateUrl string, profile map[string]interface{}) (bool, string, error) {
	profileJson, err := json.Marshal(profile)
	if err != nil {
		return false, "", err
	}

	// Validate from index service
	res, err := http.Post(validateUrl, "application/json", bytes.NewBuffer(profileJson))
	if err != nil {
		return false, "", err
	}

	var resBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&resBody)
	if res.StatusCode != 200 {
		if resBody["errors"] != nil {
			var errors []string
			for _, item := range resBody["errors"].([]interface{}) {
				errors = append(errors, item.(string))
			}
			errorsStr := strings.Join(errors, ",")
			return false, errorsStr, nil
		}
		return false, "failed without reasons!", nil
	}
	return true, "", nil
}

func PostIndex(postNodeUrl string, profileUrl string) (string, error) {
	postProfile := make(map[string]string)
	postProfile["profile_url"] = profileUrl
	postProfileJson, err := json.Marshal(postProfile)
	if err != nil {
		errStr := "error when trying to marshal a profile, url: " + postProfile["profile_url"]
		logger.Error(errStr, err)
	}
	res, err := http.Post(postNodeUrl, "application/json", bytes.NewBuffer(postProfileJson))
	if err != nil {
		return "", err
	}
	// todo: need to show error message
	if res.StatusCode != 200 {
		return "", fmt.Errorf("post failed, the status code is " + strconv.Itoa(res.StatusCode) + ", url: " + postProfile["profile_url"])
	}

	// get post node body response
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var nodeData NodeData
	err = json.Unmarshal(bodyBytes, &nodeData)
	if err != nil {
		return "", err
	}

	return nodeData.Data.NodeId, nil
}
