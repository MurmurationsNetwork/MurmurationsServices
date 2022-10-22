package main

import (
	"encoding/json"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/internal/service"
	"github.com/lucsky/cuid"
	"io"
	"net/http"
	"os"
	"strings"
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
	apiEntry := "https://api.ofdb.io/v0/entries/"

	svc := service.NewProfileService(db.NewProfileRepository(mongo.Client.GetClient()))

	curTime := time.Now().Unix()
	refreshBefore := curTime - config.Conf.RefreshTTL
	profiles, err := svc.FindLessThan(refreshBefore)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error("no profile found", err)
		} else {
			logger.Error("failed to find data from profiles", err)
		}
		cleanUp()
	}

	// get mapping
	mapping, err := importutil.GetMapping(schemaName)
	if err != nil {
		logger.Error("failed to get mapping", err)
		cleanUp()
	}

	// check the profile status
	for _, profile := range profiles {
		url := apiEntry + profile.Oid
		res, err := http.Get(url)
		if err != nil {
			logger.Error("failed to get data from api, profile cuid:"+profile.Cuid, err)
			cleanUp()
		}
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error("failed to read data from api, profile cuid:"+profile.Cuid, err)
			cleanUp()
		}

		var profileData []interface{}
		err = json.Unmarshal(bodyBytes, &profileData)
		if err != nil {
			logger.Error("failed to unmarshal data from api, profile cuid:"+profile.Cuid, err)
			cleanUp()
		}

		// If the node still exist, don't delete it and update access_time
		if len(profileData) > 0 {
			profileJson := importutil.MapFieldsName(profileData[0].(map[string]interface{}), mapping)
			doc, err := json.Marshal(profileJson)
			if err != nil {
				logger.Error("failed to marshal data, profile cuid:"+profile.Cuid, err)
				cleanUp()
			}
			profileHash, err := importutil.Hash(string(doc))
			if err != nil {
				logger.Error("failed to hash data, profile cuid:"+profile.Cuid, err)
				cleanUp()
			}

			if profileHash != profile.SourceDataHash {
				logger.Info("source data hash mismatch: " + profile.Cuid + " - " + profile.Oid + " : " + profile.SourceDataHash + " - " + profileHash)

				// reconstruct data
				profileJson, err = importutil.MapProfile(profileData[0].(map[string]interface{}), mapping, schemaName)
				if err != nil {
					logger.Error("map profile failed, profile id is "+profile.Oid, err)
					cleanUp()
				}
				oid := profileJson["oid"].(string)

				if profileJson["primary_url"] == nil {
					logger.Info("primary_url is empty, profile id is " + oid)
					continue
				}

				// validate data
				validateUrl := config.Conf.Index.URL + "/v2/validate"
				isValid, failureReasons, err := importutil.Validate(validateUrl, profileJson)
				if err != nil {
					logger.Error("validate profile failed, profile id is "+profile.Oid+". error message: ", err)
					cleanUp()
				}
				if !isValid {
					logger.Info("validate profile failed, profile id is " + profile.Oid + ". failure reasons: " + failureReasons)
					cleanUp()
				}
				profileSvc := service.NewProfileService(db.NewProfileRepository(mongo.Client.GetClient()))
				// save to Mongo
				count, err := profileSvc.Count(profile.Oid)
				if err != nil {
					logger.Error("can't count profile, profile id is "+profile.Oid, err)
					cleanUp()
				}
				if count <= 0 {
					profileJson["cuid"] = cuid.New()
					err := profileSvc.Add(profileJson)
					if err != nil {
						logger.Error("can't add a profile, profile id is "+profile.Oid, err)
						cleanUp()
					}
				} else {
					result, err := profileSvc.Update(profile.Oid, profileJson)
					if err != nil {
						logger.Error("can't update a profile, profile id is "+profile.Oid, err)
						cleanUp()
					}
					profileJson["cuid"] = result["cuid"]
				}

				// post update to Index
				postNodeUrl := config.Conf.Index.URL + "/v2/nodes"
				profileUrl := config.Conf.DataProxy.URL + "/v1/profiles/" + profileJson["cuid"].(string)
				nodeId, err := importutil.PostIndex(postNodeUrl, profileUrl)
				if err != nil {
					logger.Error("failed to post profile to Index, profile url is "+profileUrl, err)
					cleanUp()
				}

				// save node_id to profile
				err = profileSvc.UpdateNodeId(oid, nodeId)
				if err != nil {
					logger.Error("update node id failed. profile id is "+oid, err)
					cleanUp()
				}
			} else {
				err = svc.UpdateAccessTime(profile.Oid)
				if err != nil {
					logger.Error("failed to update profile's access time, profile cuid:"+profile.Cuid, err)
					cleanUp()
				}
			}
		} else {
			err = svc.Delete(profile.Cuid)
			if err != nil {
				logger.Error("failed to delete data from profiles, profile cuid:"+profile.Cuid, err)
				cleanUp()
			}
			deleteNodeUrl := config.Conf.Index.URL + "/v2/nodes/" + profile.NodeId

			client := &http.Client{}
			req, err := http.NewRequest(http.MethodDelete, deleteNodeUrl, nil)
			if err != nil {
				logger.Error("failed to delete data from index service, profile node id:"+profile.NodeId, err)
				cleanUp()
			}
			res, err = client.Do(req)
			if err != nil {
				logger.Error("failed to delete data from index service, profile node id:"+profile.NodeId, err)
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
					logger.Info("failed to delete data from index service, profile node id:" + profile.NodeId + ", error message: " + errorsStr)
				} else {
					logger.Info("failed to delete data from index service, profile node id:" + profile.NodeId + ".")
				}
			}
		}
	}

	cleanUp()
}
