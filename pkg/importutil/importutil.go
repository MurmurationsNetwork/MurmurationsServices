package importutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
)

var KvmCategory = map[string]string{
	"2cd00bebec0c48ba9db761da48678134": "#non-profit",
	"77b3c33a92554bcf8e8c2c86cedd6f6f": "#commercial",
	"c2dc278a2d6a4b9b8a50cb606fc017ed": "#event",
}

type Node struct {
	NodeID     string `json:"node_id,omitempty"`
	ProfileURL string `json:"profile_url,omitempty"`
	Status     string `json:"status,omitempty"`
}

type NodeData struct {
	Data   Node
	Errors []interface{} `json:"errors,omitempty"`
}

func GetMapping(schemaName string) (map[string]string, error) {
	filter := bson.M{"schema": schemaName}
	result := mongo.Client.FindOne(constant.MongoIndex.Mapping, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(
				"Could not find mapping for the following schema: %s",
				schemaName,
			)
		}
		return nil, fmt.Errorf(
			"Error when trying to find schema mapping; error message: %s",
			result.Err(),
		)
	}
	schemaRaw := make(map[string]interface{})
	err := result.Decode(schemaRaw)
	if err != nil {
		return nil, fmt.Errorf(
			"Error when trying to parse database response; error message: %s",
			err,
		)
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

func MapFieldsName(
	profile map[string]interface{},
	mapping map[string]string,
) map[string]interface{} {
	profileJSON := make(map[string]interface{})

	for k, v := range mapping {
		if profile[v] == nil || profile[v] == "" {
			continue
		}
		// Truncate latitude & longitude after the 8th decimal place since extra precision is superfluous
		if k == "latitude" || k == "longitude" {
			precision := math.Pow(10, float64(8))
			truncatedValue := math.Round(
				profile[v].(float64)*precision,
			) / precision
			profileJSON[k] = truncatedValue
			continue
		}
		// Trim extra space (except in tags and kvm_category)
		if k != "tags" && k != "kvm_category" {
			profileJSON[k] = strings.TrimSpace(profile[v].(string))
			continue
		}
		profileJSON[k] = profile[v]
	}

	return profileJSON
}

func MapProfile(
	profile map[string]interface{},
	mapping map[string]string,
	schema string,
) (map[string]interface{}, error) {
	// Convert KVM field names to Org Schema field names
	profileJSON := MapFieldsName(profile, mapping)

	// Hash the updated data
	profileHash, err := HashProfile(profileJSON)
	if err != nil {
		return nil, err
	}

	// hash
	profileJSON["source_data_hash"] = profileHash
	// oid
	profileJSON["oid"] = profile["id"]
	// schema
	s := []string{schema}
	profileJSON["linked_schemas"] = s
	// metadata
	metadata := map[string]interface{}{
		"sources": []map[string]interface{}{
			{
				"name":             "Karte von Morgen / Map of Tomorrow",
				"url":              "https://kartevonmorgen.org",
				"profile_data_url": "https://api.ofdb.io/v0/entries/" + profileJSON["oid"].(string),
				"access_time":      time.Now().Unix(),
			},
		},
	}
	profileJSON["metadata"] = metadata

	// Replace kvm_category with real name
	if profileJSON["kvm_category"] != nil {
		categoriesInterface := profileJSON["kvm_category"].([]interface{})
		categoriesString := make([]string, len(categoriesInterface))
		for i, v := range categoriesInterface {
			categoriesString[i] = KvmCategory[v.(string)]
		}
		profileJSON["kvm_category"] = categoriesString
	}

	return profileJSON, nil
}

func HashProfile(profile map[string]interface{}) (string, error) {
	doc, err := json.Marshal(profile)
	if err != nil {
		return "", err
	}
	profileHash, err := Hash(string(doc))
	if err != nil {
		return "", err
	}
	return profileHash, nil
}

func Validate(
	validateURL string,
	profile map[string]interface{},
) (bool, string, error) {
	profileJSON, err := json.Marshal(profile)
	if err != nil {
		return false, "", err
	}

	// Validate from index service
	res, err := http.Post(
		validateURL,
		"application/json",
		bytes.NewBuffer(profileJSON),
	)
	if err != nil {
		return false, "", fmt.Errorf("failed to post request: %w", err)
	}

	var resBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return false, "", fmt.Errorf("failed to decode response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		if resBody["errors"] != nil {
			var errors []string
			for _, item := range resBody["errors"].([]interface{}) {
				errors = append(errors, fmt.Sprintf("%v", item))
			}
			errorsStr := strings.Join(errors, ",")
			return false, errorsStr, nil
		}
		return false, "failed without reasons!", nil
	}

	return true, "", nil
}

func PostIndex(postNodeURL string, profileURL string) (string, error) {
	postProfile := make(map[string]string)
	postProfile["profile_url"] = profileURL
	postProfileJSON, err := json.Marshal(postProfile)
	if err != nil {
		errStr := "Error when trying to marshal a profile at `profile_url`: " + postProfile["profile_url"]
		logger.Error(errStr, err)
	}
	res, err := http.Post(
		postNodeURL,
		"application/json",
		bytes.NewBuffer(postProfileJSON),
	)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		var resBody map[string]interface{}

		err = json.NewDecoder(res.Body).Decode(&resBody)
		if err != nil {
			return "", err
		}

		if resBody["errors"] != nil {
			var errors []string
			for _, item := range resBody["errors"].([]interface{}) {
				errors = append(errors, fmt.Sprintf("%#v", item))
			}
			errorsStr := strings.Join(errors, ",")
			return "", fmt.Errorf(
				"Post failed with status code: " + strconv.Itoa(
					res.StatusCode,
				) + " at `profile_url`: " + postProfile["profile_url"] + " with error: " + errorsStr,
			)
		}
		return "", fmt.Errorf(
			"Post failed with status code: " + strconv.Itoa(
				res.StatusCode,
			) + "at `profile_url`: " + postProfile["profile_url"],
		)
	}

	// Get POST /node body response
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var nodeData NodeData
	err = json.Unmarshal(bodyBytes, &nodeData)
	if err != nil {
		return "", err
	}

	return nodeData.Data.NodeID, nil
}

func DeleteIndex(deleteNodeURL string, nodeID string) error {
	req, err := http.NewRequest("DELETE", deleteNodeURL, nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("node_id", nodeID)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		var resBody map[string]interface{}

		err = json.NewDecoder(res.Body).Decode(&resBody)
		if err != nil {
			return err
		}

		if resBody["errors"] != nil {
			var errors []string
			for _, item := range resBody["errors"].([]interface{}) {
				errors = append(errors, fmt.Sprintf("%#v", item))
			}
			errorsStr := strings.Join(errors, ",")
			return fmt.Errorf(
				"Delete failed with status code: " + strconv.Itoa(
					res.StatusCode,
				) + " for `node_id`: " + nodeID + " with error: " + errorsStr,
			)
		}
		return fmt.Errorf(
			"Delete failed with status code: " + strconv.Itoa(
				res.StatusCode,
			) + " for `node_id`: " + nodeID,
		)
	}
	return nil
}
