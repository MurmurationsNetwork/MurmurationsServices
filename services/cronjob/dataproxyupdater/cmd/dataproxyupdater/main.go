package main

import (
	"encoding/json"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/service"
	"github.com/lucsky/cuid"
	"os"
	"strconv"
	"time"
)

func init() {
	global.Init()
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

	if len(mapping) == 0 {
		logger.Info("can't find the mapping: " + schemaName)
		cleanUp()
	}

	if update == nil {
		// last_updated: according to recent_changes API, it can't retrieve the data before 100 days ago, so set default as 100 days ago.
		lastUpdated := time.Now().AddDate(0, 0, -100).Unix()
		err := svc.Save(schemaName, lastUpdated, apiEntry)
		if err != nil {
			logger.Error("save update status to server failed, error message:", err)
			cleanUp()
		}
		// get newer update again
		update = svc.Get(schemaName)
	}

	if update.HasError {
		logger.Info("last error didn't solve, can't continue the cronjob")
		cleanUp()
	}

	// recent-changes API
	// only process 100 data in once and process only 7 days
	entry := update.ApiEntry + "/entries/recently-changed"
	limit := 100
	offset := 0
	until := time.Unix(update.LastUpdated, 0).AddDate(0, 0, 7).Unix()

	url := getUrl(entry, update.LastUpdated, until, limit, offset)
	profiles, err := getProfiles(url)
	if err != nil {
		logger.Info("get profile failed, error message: " + err.Error())
		cleanUp()
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
				logger.Info("can't count profile, profile id is " + oid)
				cleanUp()
			}
			if count <= 0 {
				profileJson["cuid"] = cuid.New()
				err = profileSvc.Add(profileJson)
				if err != nil {
					logger.Info("can't add profile, profile id is " + oid)
					cleanUp()
				}
			} else {
				profileSvc.Update(oid, profileJson)
			}
			total++

			// todo: post update to Index
		}
		// if the data total is less than limit, no need to request data again
		if total < limit {
			break
		}
		offset += limit
		url = getUrl(entry, update.LastUpdated, until, limit, offset)
		profiles, err = getProfiles(url)
		if err != nil {
			logger.Info("get profile failed, error message: " + err.Error())
			cleanUp()
		}
	}

	// save back to update
	svc.Update(schemaName, until)

	// todo: error handling
}

func getUrl(entry string, since int64, until int64, limit int, offset int) string {
	sinceStr := strconv.FormatInt(since, 10)
	limitStr := strconv.Itoa(limit)
	offsetStr := strconv.Itoa(offset)
	apiUrl := entry + "/?since=" + sinceStr + "&limit=" + limitStr + "&offset=" + offsetStr

	curTime := time.Now().Unix()
	if until < curTime {
		untilStr := strconv.FormatInt(until, 10)
		apiUrl += "&until=" + untilStr
	}
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
