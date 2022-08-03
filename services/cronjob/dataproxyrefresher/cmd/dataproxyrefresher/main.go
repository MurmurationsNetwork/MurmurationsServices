package main

import (
	"encoding/json"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/internal/service"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func init() {
	global.Init()
}

func cleanUp() {
	mongo.Client.Disconnect()
	os.Exit(0)
}

type Node struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func main() {
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

		var profileData map[string]interface{}
		err = json.Unmarshal(bodyBytes, &profileData)
		if err != nil {
			logger.Error("failed to unmarshal data from api, profile cuid:"+profile.Cuid, err)
			cleanUp()
		}

		// If the node still exist, don't delete it and update access_time
		if len(profileData) > 0 {
			doc, err := json.Marshal(profileData)
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
				cleanUp()
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

			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logger.Error("failed to read response when deleting data from index service, profile node id:"+profile.NodeId, err)
				cleanUp()
			}

			var node Node
			err = json.Unmarshal(resBody, &node)
			if err != nil {
				logger.Error("failed to unmarshal response when deleting data from index service, profile node id:"+profile.NodeId, err)
				cleanUp()
			}

			if node.Status != 200 {
				logger.Info("failed to delete data from index service, profile node id:" + profile.NodeId + ", error message: " + node.Message)
			}
		}
	}

	cleanUp()
}
