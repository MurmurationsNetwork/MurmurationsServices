package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/service"
	"github.com/lucsky/cuid"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	NodeId         string   `json:"node_id"`
	ProfileUrl     string   `json:"profile_url"`
	Status         string   `json:"status"`
	FailureReasons []string `json:"failure_reasons"`
}
type NodeData struct {
	Data Node
}

func init() {
	global.Init()
}

func errCleanUp(schemaName string, svc service.UpdatesService, errStr string) {
	err := svc.SaveError(schemaName, errStr)
	if err != nil {
		logger.Fatal("save error message failed", err)
	}
	cleanUp()
}

func cleanUp() {
	mongo.Client.Disconnect()
	os.Exit(0)
}

func main() {
	schemaName := "karte_von_morgen-v1.0.0"
	apiEntry := "https://api.ofdb.io/v0"

	svc := service.NewUpdateService(db.NewUpdateRepository(mongo.Client.GetClient()))
	mappingSvc := service.NewMappingService(db.NewMappingRepository(mongo.Client.GetClient()))
	profileSvc := service.NewProfileService(db.NewProfileRepository(mongo.Client.GetClient()))

	update := svc.Get(schemaName)
	mapping := mappingSvc.Get(schemaName)

	if update == nil {
		// last_updated: according to recent_changes API, it can't retrieve the data before 100 days ago, so set default as 100 days ago.
		lastUpdated := time.Now().AddDate(0, 0, -100).Unix()
		err := svc.Save(schemaName, lastUpdated, apiEntry)
		if err != nil {
			errStr := "save update status to server failed" + err.Error()
			logger.Error("save update status to server failed", err)
			errCleanUp(schemaName, svc, errStr)
		}
		// get newer update again
		update = svc.Get(schemaName)
	}

	// if the last error didn't solve, don't run
	if update.HasError {
		logger.Info("last error didn't solve, can't continue the cronjob")
		logger.Info("last error: " + update.ErrorMessage)
		cleanUp()
	}

	if len(mapping) == 0 {
		errStr := "can't find the mapping: " + schemaName
		logger.Info(errStr)
		errCleanUp(schemaName, svc, errStr)
	}

	// recent-changes API
	// only process 100 data in once
	entry := update.ApiEntry + "/entries/recently-changed"
	limit := 100
	offset := 0
	until := time.Now().Unix()

	url := getUrl(entry, update.LastUpdated, until, limit, offset)
	profiles, err := getProfiles(url)
	if err != nil {
		errStr := "get profile failed" + err.Error()
		logger.Error("get profile failed", err)
		errCleanUp(schemaName, svc, errStr)
	}
	for len(profiles) > 0 {
		total := 0
		for _, value := range profiles {
			var profileJson map[string]interface{}
			profileJson = mapData(value, mapping, schemaName)
			oid := profileJson["oid"].(string)
			if profileJson["primary_url"] == nil {
				logger.Info("primary_url is empty, profile id is " + oid)
				continue
			}
			count, err := profileSvc.Count(oid)
			if err != nil {
				errStr := "can't count profile, profile id is " + oid
				logger.Info(errStr)
				errCleanUp(schemaName, svc, errStr)
			}
			if count <= 0 {
				profileJson["cuid"] = cuid.New()
				err = profileSvc.Add(profileJson)
				if err != nil {
					errStr := "can't add profile, profile id is " + oid
					logger.Info(errStr)
					errCleanUp(schemaName, svc, errStr)
				}
			} else {
				profileSvc.Update(oid, profileJson)
			}
			total++

			// post update to Index
			postNodeUrl := config.Conf.Index.URL + "/v2/nodes"
			postProfile := make(map[string]string)
			postProfile["profile_url"] = config.Conf.DataProxy.URL + "/v1/profiles/" + profileJson["cuid"].(string)
			postProfileJson, err := json.Marshal(postProfile)
			if err != nil {
				errStr := "error when trying to marshal a profile, url: " + postProfile["profile_url"]
				logger.Error(errStr, err)
				errCleanUp(schemaName, svc, errStr)
			}
			res, err := http.Post(postNodeUrl, "application/json", bytes.NewBuffer(postProfileJson))
			if err != nil {
				errStr := "error when trying to post a profile"
				logger.Error(errStr, err)
				errCleanUp(schemaName, svc, errStr)
			}
			if res.StatusCode != 200 {
				errStr := "post failed, the status code is " + strconv.Itoa(res.StatusCode) + ", url: " + postProfile["profile_url"]
				logger.Info(errStr)
				errCleanUp(schemaName, svc, errStr)
			}

			// get post node body response
			bodyBytes, err := io.ReadAll(res.Body)
			if err != nil {
				errStr := "read post body failed. url: " + postProfile["profile_url"]
				logger.Error(errStr, err)
				errCleanUp(schemaName, svc, errStr)
			}

			var nodeData NodeData
			err = json.Unmarshal(bodyBytes, &nodeData)
			if err != nil {
				errStr := "unmarshal body failed. url: " + postProfile["profile_url"]
				logger.Error(errStr, err)
				errCleanUp(schemaName, svc, errStr)
			}
			// save node_id to profile
			err = profileSvc.UpdateNodeId(oid, nodeData.Data.NodeId)
			if err != nil {
				errStr := "update node id failed. profile id is " + oid
				logger.Error(errStr, err)
				errCleanUp(schemaName, svc, errStr)
			}
		}
		// if the data total is less than limit, no need to request data again
		if total < limit {
			break
		}
		offset += limit
		url = getUrl(entry, update.LastUpdated, until, limit, offset)
		profiles, err = getProfiles(url)
		if err != nil {
			errStr := "get profile failed" + err.Error()
			logger.Error("get profile failed", err)
			errCleanUp(schemaName, svc, errStr)
		}
	}

	// save back to update
	err = svc.Update(schemaName, until)
	if err != nil {
		errStr := "failed to update the updates" + err.Error()
		logger.Error("failed to update the updates", err)
		errCleanUp(schemaName, svc, errStr)
	}

	// get profile with not posted
	notPostedProfiles, err := profileSvc.GetNotPosted()
	if err != nil {
		errStr := "failed to get not posted nodes" + err.Error()
		logger.Error("failed to get not posted nodes", err)
		errCleanUp(schemaName, svc, errStr)
	}

	for _, notPostedProfile := range notPostedProfiles {
		getNodeUrl := config.Conf.Index.URL + "/v2/nodes/" + notPostedProfile.NodeId
		res, err := http.Get(getNodeUrl)
		if err != nil {
			errStr := "failed to get not posted nodes, node id is " + notPostedProfile.NodeId + err.Error()
			logger.Error("failed to get not posted nodes", err)
			errCleanUp(schemaName, svc, errStr)
		}
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			errStr := "read post body failed. node id is " + notPostedProfile.NodeId + err.Error()
			logger.Error(errStr, err)
			errCleanUp(schemaName, svc, errStr)
		}

		var nodeData NodeData
		err = json.Unmarshal(bodyBytes, &nodeData)
		if err != nil {
			errStr := "unmarshal body failed. node id is " + notPostedProfile.NodeId + err.Error()
			logger.Error(errStr, err)
			errCleanUp(schemaName, svc, errStr)
		}

		if nodeData.Data.Status == "posted" {
			err = profileSvc.UpdateIsPosted(notPostedProfile.NodeId)
			if err != nil {
				errStr := "update isPosted failed. node id is " + notPostedProfile.NodeId + err.Error()
				logger.Error(errStr, err)
				errCleanUp(schemaName, svc, errStr)
			}
		} else {
			failureReasons := strings.Join(nodeData.Data.FailureReasons, ",")
			logger.Info(notPostedProfile.NodeId + " is not posted. Profile url is " + nodeData.Data.ProfileUrl + ". Error messages: " + failureReasons)
		}
	}

	cleanUp()
}

func getUrl(entry string, since int64, until int64, limit int, offset int) string {
	sinceStr := strconv.FormatInt(since, 10)
	limitStr := strconv.Itoa(limit)
	offsetStr := strconv.Itoa(offset)
	untilStr := strconv.FormatInt(until, 10)
	apiUrl := entry + "/?since=" + sinceStr + "&limit=" + limitStr + "&offset=" + offsetStr + "&until=" + untilStr

	return apiUrl
}

func getProfiles(url string) ([]map[string]interface{}, error) {
	res, err := httputil.Get(url)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("can't get data from" + url)
	}

	var bodyJson []map[string]interface{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&bodyJson)
	if err != nil {
		return nil, fmt.Errorf("can't parse data from" + url)
	}

	return bodyJson, nil
}

func mapData(profile map[string]interface{}, mapping map[string]string, schema string) map[string]interface{} {
	profileJson := make(map[string]interface{})

	// change field name
	for k, v := range mapping {
		if profile[v] == nil {
			continue
		}
		profileJson[k] = profile[v]
	}

	// oid
	profileJson["oid"] = profile["id"]
	// schema
	s := []string{schema}
	profileJson["linked_schemas"] = s
	// metadata
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

	return profileJson
}
