package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lucsky/cuid"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/internal/service"
)

const (
	schemaName = "karte_von_morgen-v1.0.0"
	apiEntry   = "https://api.ofdb.io/v0"
)

func init() {
	global.Init()
}

func main() {
	logger.Info("Starting dataproxy updater...")

	updater := NewUpdater()

	startTime := time.Now()

	if err := updater.Run(context.Background()); err != nil {
		logger.Panic("Error running dataproxy updater", err)
		return
	}

	duration := time.Since(startTime)
	logger.Info("Dataproxy updater completed successfully")
	logger.Info("Dataproxy updater run duration: " + duration.String())
}

type DataproxyUpdater struct {
	updateSvc  service.UpdatesService
	profileSvc service.ProfilesService
	update     *model.Update
}

func NewUpdater() *DataproxyUpdater {
	return &DataproxyUpdater{
		updateSvc: service.NewUpdateService(
			mongo.NewUpdateRepository(mongodb.Client.GetClient()),
		),
		profileSvc: service.NewProfileService(
			mongo.NewProfileRepository(mongodb.Client.GetClient()),
		),
	}
}

func (u *DataproxyUpdater) Run(ctx context.Context) error {
	var err error

	defer func() {
		if err != nil {
			err = u.updateSvc.SaveError(schemaName, err.Error())
		}
	}()
	defer cleanup()

	if err = u.initializeUpdateStatus(ctx); err != nil {
		return fmt.Errorf("initializing update status failed: %w", err)
	}

	if err = u.processRecentChanges(ctx); err != nil {
		return fmt.Errorf("processing recent changes failed: %w", err)
	}

	return nil
}

func (u *DataproxyUpdater) initializeUpdateStatus(_ context.Context) error {
	update := u.updateSvc.Get(schemaName)

	// If the update status doesn't exist, initialize it.
	if update == nil {
		// According to the recent_changes API, it can't retrieve the data before
		// 100 days ago, so set default as 100 days ago.
		lastUpdated := time.Now().AddDate(0, 0, -100).Unix()

		err := u.updateSvc.Save(schemaName, lastUpdated, apiEntry)
		if err != nil {
			return fmt.Errorf("save update status to server failed: %v", err)
		}
		// Fetch the newer update again.
		update = u.updateSvc.Get(schemaName)
	}

	if update.HasError {
		return fmt.Errorf("last error didn't solve: %s", update.ErrorMessage)
	}
	return nil
}

func (u *DataproxyUpdater) processRecentChanges(ctx context.Context) error {
	mapping, err := importutil.GetMapping(schemaName)
	if err != nil {
		return fmt.Errorf("get mapping failed: %v", err)
	}

	if len(mapping) == 0 {
		return fmt.Errorf("can't find the mapping: %s", schemaName)
	}

	// recent-changes API
	limit := 100
	offset := 0
	until := time.Now().Unix()

	url := buildQueryURL(
		u.update.APIEntry+"/entries/recently-changed",
		u.update.LastUpdated,
		until,
		limit,
		offset,
	)
	profiles, err := getProfiles(url)
	if err != nil {
		return fmt.Errorf("get profile failed: %v", err)
	}

	for len(profiles) > 0 {
		for _, oldProfile := range profiles {
			profileJSON, err := importutil.MapProfile(
				oldProfile,
				mapping,
				schemaName,
			)
			if err != nil {
				return fmt.Errorf(
					"map profile failed, profile id is %v. error message: %v",
					oldProfile["id"],
					err,
				)
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
				return fmt.Errorf(
					"validate profile failed, profile id is %s. error message: %v",
					oid,
					err,
				)
			}
			if !isValid {
				logger.Info(
					fmt.Sprintf(
						"validate profile failed, profile id is %s. failure reasons: %s",
						oid,
						failureReasons,
					),
				)
				continue
			}

			// save to Mongo
			if err := u.saveOrUpdateProfile(ctx, oid, profileJSON); err != nil {
				return err
			}
		}

		offset += limit
		url = buildQueryURL(
			u.update.APIEntry+"/entries/recently-changed",
			u.update.LastUpdated,
			until,
			limit,
			offset,
		)
		profiles, err = getProfiles(url)
		if err != nil {
			return fmt.Errorf("get profile failed: %v", err)
		}
	}

	// Update the last update time.
	if err = u.updateSvc.Update(schemaName, until); err != nil {
		return fmt.Errorf("failed to update the updates: %v", err)
	}

	return nil
}

func (u *DataproxyUpdater) saveOrUpdateProfile(
	_ context.Context,
	oid string,
	profileJSON map[string]interface{},
) error {
	// Check if profile exists.
	count, err := u.profileSvc.Count(oid)
	if err != nil {
		return fmt.Errorf("can't count profile, profile id is %s: %v", oid, err)
	}

	// If profile doesn't exist, add it.
	if count <= 0 {
		profileJSON["cuid"] = cuid.New()
		if err := u.profileSvc.Add(profileJSON); err != nil {
			return fmt.Errorf(
				"can't add a profile, profile id is %s: %v",
				oid,
				err,
			)
		}
	} else {
		// If profile exists, update it.
		updatedProfile, err := u.profileSvc.Update(oid, profileJSON)
		if err != nil {
			return fmt.Errorf("can't update a profile, profile id is %s: %v", oid, err)
		}
		profileJSON["cuid"] = updatedProfile["cuid"]
	}

	// Post profile details to Index.
	profileURL := fmt.Sprintf(
		"%s/v1/profiles/%s",
		config.Conf.DataProxy.URL,
		profileJSON["cuid"].(string),
	)
	nodeID, err := importutil.PostIndex(
		fmt.Sprintf("%s/v2/nodes", config.Conf.Index.URL),
		profileURL,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to post profile to Index, profile url is %s: %v",
			profileURL,
			err,
		)
	}

	// Save node_id to profile.
	if err := u.profileSvc.UpdateNodeID(oid, nodeID); err != nil {
		return fmt.Errorf(
			"update node id failed. profile id is %s: %v",
			oid,
			err,
		)
	}

	return nil
}

func buildQueryURL(entry string, since, until int64, limit, offset int) string {
	apiURL := fmt.Sprintf(
		"%s/?since=%d&limit=%d&offset=%d&until=%d",
		entry, since, limit, offset, until,
	)
	return apiURL
}

// getProfiles fetches and decodes JSON data from the provided URL.
func getProfiles(url string) ([]map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Decode the JSON response.
	var profiles []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
		return nil, fmt.Errorf("failed to parse data from %s: %w", url, err)
	}

	return profiles, nil
}

func cleanup() {
	mongodb.Client.Disconnect()
	os.Exit(0)
}
