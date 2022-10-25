package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/internal/service"
	"github.com/lucsky/cuid"
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
	apiEntry := "https://api.ofdb.io/v0/entries/"

	svc := service.NewProfileService(db.NewProfileRepository(mongo.Client.GetClient()))

	curTime := time.Now().Unix()
	refreshBefore := curTime - config.Conf.RefreshTTL
	profiles, err := svc.FindLessThan(refreshBefore)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error("No profile found.", err)
		} else {
			logger.Error("Failed to find data from profiles.", err)
		}
		cleanUp()
	}

	// get mapping
	mapping, err := importutil.GetMapping(schemaName)
	if err != nil {
		logger.Error("Failed to get mapping.", err)
		cleanUp()
	}

	// check the profile status
	for _, profile := range profiles {
		url := apiEntry + profile.Oid
		res, err := http.Get(url)
		if err != nil {
			logger.Error("Failed to get data from API. Profile CUID:"+profile.Cuid, err)
			cleanUp()
		}
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error("Failed to read data from API. Profile CUID:"+profile.Cuid, err)
			cleanUp()
		}

		var profileData []interface{}
		err = json.Unmarshal(bodyBytes, &profileData)
		if err != nil {
			logger.Error("Failed to unmarshal data from API. Profile CUID:"+profile.Cuid, err)
			cleanUp()
		}

		// If the node still exist, don't delete it and update access_time
		if len(profileData) > 0 {
			profileJson := importutil.MapFieldsName(profileData[0].(map[string]interface{}), mapping)
			doc, err := json.Marshal(profileJson)
			if err != nil {
				logger.Error("Failed to marshal data. Profile CUID: "+profile.Cuid, err)
				cleanUp()
			}
			profileHash, err := importutil.Hash(string(doc))
			if err != nil {
				logger.Error("Failed to hash data. Profile CUID: "+profile.Cuid, err)
				cleanUp()
			}

			if profileHash != profile.SourceDataHash {
				logger.Info("Source data hash mismatch: " + profile.Cuid + " - " + profile.Oid + " : " + profile.SourceDataHash + " - " + profileHash)

				// reconstruct data
				profileJson, err = importutil.MapProfile(profileData[0].(map[string]interface{}), mapping, schemaName)
				if err != nil {
					logger.Error("Map profile failed. Profile ID: "+profile.Oid, err)
					cleanUp()
				}
				oid := profileJson["oid"].(string)

				if profileJson["primary_url"] == nil {
					logger.Info("The primary_url is empty. Profile ID: " + oid)
					continue
				}

				// validate data
				validateUrl := config.Conf.Index.URL + "/v2/validate"
				isValid, failureReasons, err := importutil.Validate(validateUrl, profileJson)
				if err != nil {
					logger.Error("Validate profile failed. Profile ID: "+profile.Oid+". error message: ", err)
					cleanUp()
				}
				if !isValid {
					logger.Info("Validate profile failed. Profile ID: " + profile.Oid + ". failure reasons: " + failureReasons)
					cleanUp()
				}
				profileSvc := service.NewProfileService(db.NewProfileRepository(mongo.Client.GetClient()))
				// save to Mongo
				count, err := profileSvc.Count(profile.Oid)
				if err != nil {
					logger.Error("Can't count profile. Profile ID: "+profile.Oid, err)
					cleanUp()
				}
				if count <= 0 {
					profileJson["cuid"] = cuid.New()
					err := profileSvc.Add(profileJson)
					if err != nil {
						logger.Error("Can't add a profile. Profile ID: "+profile.Oid, err)
						cleanUp()
					}
				} else {
					result, err := profileSvc.Update(profile.Oid, profileJson)
					if err != nil {
						logger.Error("Can't update a profile. Profile ID: "+profile.Oid, err)
						cleanUp()
					}
					profileJson["cuid"] = result["cuid"]
				}

				// post update to Index
				postNodeUrl := config.Conf.Index.URL + "/v2/nodes"
				profileUrl := config.Conf.DataProxy.URL + "/v1/profiles/" + profileJson["cuid"].(string)
				nodeId, err := importutil.PostIndex(postNodeUrl, profileUrl)
				if err != nil {
					logger.Error("Failed to post profile to Index. Profile URL: "+profileUrl, err)
					cleanUp()
				}

				// save node_id to profile
				err = profileSvc.UpdateNodeId(oid, nodeId)
				if err != nil {
					logger.Error("Update node id failed. Profile ID: "+oid, err)
					cleanUp()
				}
			} else {
				err = svc.UpdateAccessTime(profile.Oid)
				if err != nil {
					logger.Error("Failed to update profile's access time. Profile CUID: "+profile.Cuid, err)
					cleanUp()
				}
			}
		} else {
			err = svc.Delete(profile.Cuid)
			if err != nil {
				logger.Error("Failed to delete data from profiles. Profile CUID: "+profile.Cuid, err)
				cleanUp()
			}
			deleteNodeUrl := config.Conf.Index.URL + "/v2/nodes/" + profile.NodeId

			client := &http.Client{}
			req, err := http.NewRequest(http.MethodDelete, deleteNodeUrl, nil)
			if err != nil {
				logger.Error("Failed to delete data from Index service Profile node ID: "+profile.NodeId, err)
				cleanUp()
			}
			res, err = client.Do(req)
			if err != nil {
				logger.Error("Failed to delete data from Index service. Profile node ID: "+profile.NodeId, err)
				cleanUp()
			}
			defer res.Body.Close()

			if res.StatusCode != 200 {
				var resBody map[string]interface{}
				json.NewDecoder(res.Body).Decode(&resBody)
				if resBody["errors"] != nil {
					var errors []string
					for _, item := range resBody["errors"].([]interface{}) {
						errors = append(errors, fmt.Sprintf("%#v", item))
					}
					errorsStr := strings.Join(errors, ",")
					logger.Info("Failed to delete data from Index service. Profile node ID: " + profile.NodeId + " - Error message: " + errorsStr)
				} else {
					logger.Info("Failed to delete data from Index service. Profile node ID: " + profile.NodeId)
				}
			}
		}
	}

	cleanUp()
}
