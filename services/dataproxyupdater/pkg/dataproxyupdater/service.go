package dataproxyupdater

import (
	"encoding/json"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyupdater/internal/model"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lucsky/cuid"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyupdater/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyupdater/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyupdater/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyupdater/internal/service"
)

func init() {
	global.Init()
}

func cleanUpWithError(schemaName string, svc service.UpdatesService, errStr string, errStatus ...int) {
	status := -1
	if len(errStr) > 0 {
		status = errStatus[0]
	}
	err := svc.SaveError(schemaName, true, errStr, status)
	if err != nil {
		logger.Fatal("save error message failed", err)
	}
	cleanUp()
}

func cleanUp() {
	mongodb.Client.Disconnect()
	os.Exit(0)
}

func removeError(schemaName string, svc service.UpdatesService) {
	err := svc.SaveError(schemaName, false, "", model.ErrorStatusOK)
	if err != nil {
		logger.Fatal("remove error message failed", err)
	}
}

func Run() {
	schemaName := "karte_von_morgen-v1.0.0"
	apiEntry := "https://api.ofdb.io/v0"

	svc := service.NewUpdateService(
		mongo.NewUpdateRepository(mongodb.Client.GetClient()),
	)
	profileSvc := service.NewProfileService(
		mongo.NewProfileRepository(mongodb.Client.GetClient()),
	)

	update := svc.Get(schemaName)

	if update == nil {
		// last_updated: according to recent_changes API, it can't retrieve the data before 100 days ago,
		// so set default as 100 days ago.
		lastUpdated := time.Now().AddDate(0, 0, -100).Unix()
		err := svc.Save(schemaName, lastUpdated, apiEntry)
		if err != nil {
			errStr := "save update status to server failed: " + err.Error()
			logger.Error("save update status to server failed", err)
			cleanUpWithError(schemaName, svc, errStr)
		}
		// get newer update again
		update = svc.Get(schemaName)
	}

	if update.HasError && update.ErrorStatus == model.ErrorStatusAPIUnavailable {
		url := getURL(update.APIEntry+"/entries/recently-changed", update.LastUpdated, time.Now().Unix(), 1, 0)
		_, err := getProfiles(url)
		if err == nil {
			removeError(schemaName, svc)
			// get newer update again
			update = svc.Get(schemaName)
		}
	}

	// if the last error didn't solve, don't run
	if update.HasError {
		logger.Info("last error didn't solve, can't continue the cronjob")
		logger.Info("last error: " + update.ErrorMessage)
		cleanUp()
	}

	mapping, err := importutil.GetMapping(schemaName)
	if err != nil {
		errStr := "get mapping failed: " + err.Error()
		logger.Error("get mapping failed", err)
		cleanUpWithError(schemaName, svc, errStr)
	}

	if len(mapping) == 0 {
		errStr := "can't find the mapping: " + schemaName
		logger.Info(errStr)
		cleanUpWithError(schemaName, svc, errStr)
	}

	// recent-changes API
	// only process 100 data in once
	entry := update.APIEntry + "/entries/recently-changed"
	limit := 100
	offset := 0
	until := time.Now().Unix()

	url := getURL(entry, update.LastUpdated, until, limit, offset)
	profiles, err := getProfiles(url)
	if err != nil {
		errStr := "get profile failed: " + err.Error()
		logger.Error("get profile failed", err)
		cleanUpWithError(schemaName, svc, errStr, model.ErrorStatusAPIUnavailable)
	}
	for len(profiles) > 0 {
		for _, oldProfile := range profiles {
			profileJSON, err := importutil.MapProfile(
				oldProfile,
				mapping,
				schemaName,
			)
			if err != nil {
				errStr := "map profile failed, profile id is " + oldProfile["id"].(string) + ". error message: " + err.Error()
				logger.Error("map profile failed", err)
				cleanUpWithError(schemaName, svc, errStr)
			}
			oid := profileJSON["oid"].(string)

			if profileJSON["primary_url"] == nil {
				logger.Info("primary_url is empty, profile id is " + oid)
				continue
			}

			// validate data
			validateURL := config.Conf.Index.URL + "/v2/validate"
			isValid, failureReasons, err := importutil.Validate(
				validateURL,
				profileJSON,
			)
			if err != nil {
				errStr := "validate profile failed, profile id is " + oid + ". error message: " + err.Error()
				logger.Error("validate profile failed", err)
				cleanUpWithError(schemaName, svc, errStr)
			}
			if !isValid {
				errStr := "validate profile failed, profile id is " + oid + ". failure reasons: " + failureReasons
				logger.Info(errStr)
				continue
			}

			// save to Mongo
			count, err := profileSvc.Count(oid)
			if err != nil {
				errStr := "can't count profile, profile id is " + oid
				logger.Info(errStr)
				cleanUpWithError(schemaName, svc, errStr)
			}
			if count <= 0 {
				profileJSON["cuid"] = cuid.New()
				err := profileSvc.Add(profileJSON)
				if err != nil {
					errStr := "can't add a profile, profile id is " + oid
					logger.Info(errStr)
					cleanUpWithError(schemaName, svc, errStr)
				}
			} else {
				result, err := profileSvc.Update(oid, profileJSON)
				if err != nil {
					errStr := "can't update a profile, profile id is " + oid
					logger.Info(errStr)
					cleanUpWithError(schemaName, svc, errStr)
				}
				profileJSON["cuid"] = result["cuid"]
			}

			// post update to Index
			postNodeURL := config.Conf.Index.URL + "/v2/nodes"
			profileURL := config.Conf.DataProxy.URL + "/v1/profiles/" + profileJSON["cuid"].(string)
			nodeID, err := importutil.PostIndex(postNodeURL, profileURL)
			if err != nil {
				errStr := "failed to post profile to Index, profile url is " + profileURL + ". error message: " + err.Error()
				logger.Error(errStr, err)
				cleanUpWithError(schemaName, svc, errStr)
			}

			// save node_id to profile
			err = profileSvc.UpdateNodeID(oid, nodeID)
			if err != nil {
				errStr := "update node id failed. profile id is " + oid
				logger.Error(errStr, err)
				cleanUpWithError(schemaName, svc, errStr)
			}
		}
		offset += limit
		url = getURL(entry, update.LastUpdated, until, limit, offset)
		profiles, err = getProfiles(url)
		if err != nil {
			errStr := "get profile failed: " + err.Error()
			logger.Error("get profile failed", err)
			cleanUpWithError(schemaName, svc, errStr, model.ErrorStatusAPIUnavailable)
		}
	}

	// save back to update
	err = svc.Update(schemaName, until)
	if err != nil {
		errStr := "failed to update the updates: " + err.Error()
		logger.Error("failed to update the updates", err)
		cleanUpWithError(schemaName, svc, errStr)
	}

	// get profile with not posted
	notPostedProfiles, err := profileSvc.GetNotPosted()
	if err != nil {
		errStr := "failed to get not posted nodes: " + err.Error()
		logger.Error("failed to get not posted nodes", err)
		cleanUpWithError(schemaName, svc, errStr)
	}

	for _, notPostedProfile := range notPostedProfiles {
		getNodeURL := config.Conf.Index.URL + "/v2/nodes/" + notPostedProfile.NodeID
		res, err := http.Get(getNodeURL)
		if err != nil {
			errStr := "failed to get not posted nodes, node id is " + notPostedProfile.NodeID + err.Error()
			logger.Error("failed to get not posted nodes", err)
			cleanUpWithError(schemaName, svc, errStr)
		}
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			errStr := "read post body failed. node id is " + notPostedProfile.NodeID + err.Error()
			logger.Error(errStr, err)
			cleanUpWithError(schemaName, svc, errStr)
		}

		var nodeData importutil.NodeData
		err = json.Unmarshal(bodyBytes, &nodeData)
		if err != nil {
			errStr := "unmarshal body failed. node id is " + notPostedProfile.NodeID + err.Error()
			logger.Error(errStr, err)
			cleanUpWithError(schemaName, svc, errStr)
		}

		if res.StatusCode == 404 {
			err = profileSvc.Delete(notPostedProfile.Cuid)
			if err != nil {
				errStr := "delete profile failed. node cuid is " + notPostedProfile.Cuid + err.Error()
				logger.Error(errStr, err)
				cleanUpWithError(schemaName, svc, errStr)
			}
		}

		if nodeData.Data.Status == "posted" {
			err = profileSvc.UpdateIsPosted(notPostedProfile.NodeID)
			if err != nil {
				errStr := "update isPosted failed. node id is " + notPostedProfile.NodeID + err.Error()
				logger.Error(errStr, err)
				cleanUpWithError(schemaName, svc, errStr)
			}
		} else {
			if nodeData.Errors != nil {
				var errors []string
				for _, item := range nodeData.Errors {
					errors = append(errors, fmt.Sprintf("%#v", item))
				}
				errorsStr := strings.Join(errors, ",")
				logger.Info("node id " + notPostedProfile.NodeID + " is not posted. Profile url is " + nodeData.Data.ProfileURL + ". Error messages: " + errorsStr)
			} else {
				logger.Info("node id " + notPostedProfile.NodeID + " is not posted. Profile url is " + nodeData.Data.ProfileURL + ".")
			}
		}
	}

	cleanUp()
}

func getURL(
	entry string,
	since int64,
	until int64,
	limit int,
	offset int,
) string {
	sinceStr := strconv.FormatInt(since, 10)
	limitStr := strconv.Itoa(limit)
	offsetStr := strconv.Itoa(offset)
	untilStr := strconv.FormatInt(until, 10)
	apiURL := entry + "/?since=" + sinceStr + "&limit=" + limitStr + "&offset=" + offsetStr + "&until=" + untilStr

	return apiURL
}

func getProfiles(url string) ([]map[string]interface{}, error) {
	res, err := httputil.Get(url)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get data from " + url + "with the error message: " + err.Error(),
		)
	}

	defer res.Body.Close()

	var bodyJSON []map[string]interface{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("can't parse data from: " + url)
	}

	return bodyJSON, nil
}
