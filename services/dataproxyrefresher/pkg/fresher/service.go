package fresher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lucsky/cuid"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyrefresher/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyrefresher/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyrefresher/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyrefresher/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxyrefresher/internal/service"
)

func init() {
	global.Init()
}

const (
	// SchemaName represents the name of the schema being used.
	SchemaName = "karte_von_morgen-v1.0.0"
	// APIEntry is the base endpoint for the API from which profiles are
	// retrieved.
	APIEntry = "https://api.ofdb.io/v0/entries/"

	// APIValidatePath is the API path used for validating profiles.
	APIValidatePath = "/v2/validate"
	// APINodesPath is the API path used for operations related to nodes.
	APINodesPath = "/v2/nodes"
	// APIProfilesPath is the API path used for operations related to profiles.
	APIProfilesPath = "/v1/profiles"
)

type DataproxyRefresher struct {
	svc service.ProfilesService
}

func NewRefresher() *DataproxyRefresher {
	return &DataproxyRefresher{
		svc: service.NewProfileService(
			mongo.NewProfileRepository(mongodb.Client.GetClient()),
		),
	}
}

func (r *DataproxyRefresher) Run() error {
	defer r.cleanUp()

	profiles, err := r.getProfiles(SchemaName)
	if err != nil {
		return err
	}

	mapping, err := r.getMapping(SchemaName)
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		if err := r.processProfile(mapping, profile); err != nil {
			return err
		}
	}

	return nil
}

func (r *DataproxyRefresher) getProfiles(
	schemaName string,
) ([]model.Profile, error) {
	curTime := time.Now().Unix()
	refreshBefore := curTime - config.Conf.RefreshTTL

	profiles, err := r.svc.FindLessThan(schemaName, refreshBefore)
	if err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no profile found: %w", err)
		}
		return nil, fmt.Errorf("failed to find data from profiles: %w", err)
	}

	return profiles, nil
}

func (r *DataproxyRefresher) getMapping(
	schemaName string,
) (map[string]string, error) {
	mapping, err := importutil.GetMapping(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get mapping: %w", err)
	}
	return mapping, nil
}

func (r *DataproxyRefresher) processProfile(
	mapping map[string]string,
	profile model.Profile,
) error {
	url := APIEntry + profile.Oid
	profileData, err := r.getProfileData(url, profile.Cuid)
	if err != nil {
		return fmt.Errorf("failed to get profile data: %w", err)
	}

	if len(profileData) > 0 {
		err = r.processExistingProfile(mapping, profile, profileData)
		if err != nil {
			return fmt.Errorf("failed to process existing profile: %w", err)
		}
	} else {
		err = r.deleteProfile(profile)
		if err != nil {
			return fmt.Errorf("failed to delete profile: %w", err)
		}
	}

	return nil
}

func (r *DataproxyRefresher) getProfileData(
	url string,
	cuid string,
) ([]interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get data from API for Profile CUID %s: %w",
			cuid,
			err,
		)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"received non-OK status code: %d when fetching data for Profile CUID %s",
			res.StatusCode,
			cuid,
		)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read data from API for Profile CUID %s: %w",
			cuid,
			err,
		)
	}

	var profileData []interface{}
	err = json.Unmarshal(bodyBytes, &profileData)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal data from API for Profile CUID %s: %w",
			cuid,
			err,
		)
	}

	return profileData, nil
}

func (r *DataproxyRefresher) processExistingProfile(
	mapping map[string]string,
	profile model.Profile,
	profileData []interface{},
) error {
	// Extract data from profileData and perform mapping operations.
	profileJSON, err := mapProfileData(profileData, mapping)
	if err != nil {
		return fmt.Errorf("failed to map profile data: %w", err)
	}

	// Serialize and hash the data.
	hashedData, err := hashProfileData(profileJSON)
	if err != nil {
		return fmt.Errorf("failed to hash profile data: %w", err)
	}

	// Compare hash with profile's existing hash.
	if hashedData != profile.SourceDataHash {
		// Reconstruct profileJSON.
		profileJSON, err = importutil.MapProfile(
			profileData[0].(map[string]interface{}),
			mapping,
			SchemaName,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to reconstruct profile data for Profile ID %s: %w",
				profile.Oid,
				err,
			)
		}
		if err := r.updateProfileIfValid(profile, profileJSON); err != nil {
			return fmt.Errorf("failed to update profile: %w", err)
		}
	} else {
		if err := r.updateProfileAccessTime(profile); err != nil {
			return fmt.Errorf("failed to update profile access time: %w", err)
		}
	}

	return nil
}

func mapProfileData(
	profileData []interface{},
	mapping map[string]string,
) (map[string]interface{}, error) {
	if len(profileData) == 0 {
		return nil, fmt.Errorf("no data to map")
	}

	rawData, ok := profileData[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("profile data is not in the expected format")
	}

	mappedData := importutil.MapFieldsName(rawData, mapping)

	return mappedData, nil
}

func hashProfileData(profileJSON map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(profileJSON)
	if err != nil {
		return "", fmt.Errorf("failed to marshal profile data: %w", err)
	}

	hashedData, err := jsonutil.Hash(string(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to hash profile data: %w", err)
	}

	return hashedData, nil
}

func (r *DataproxyRefresher) updateProfileIfValid(
	profile model.Profile,
	profileJSON map[string]interface{},
) error {
	isValid, failureReasons, err := r.validateProfile(profileJSON)
	if err != nil {
		return fmt.Errorf(
			"error during profile validation for Profile ID %s: %w",
			profile.Oid,
			err,
		)
	}

	if !isValid {
		logger.Info(
			fmt.Sprintf(
				"Validation failed for Profile ID %s. Failure reasons: %s",
				profile.Oid,
				failureReasons,
			),
		)
		return nil
	}

	if err := r.UpdateProfile(profile, profileJSON); err != nil {
		return fmt.Errorf(
			"failed to update profile for Profile ID %s: %w",
			profile.Oid,
			err,
		)
	}

	return nil
}

func (r *DataproxyRefresher) validateProfile(
	profileJSON map[string]interface{},
) (bool, string, error) {
	validateURL := config.Conf.Index.URL + APIValidatePath

	isValid, failureReasons, err := importutil.Validate(
		validateURL,
		profileJSON,
	)
	if err != nil {
		return false, "",
			fmt.Errorf("error during profile validation: %w", err)
	}

	return isValid, failureReasons, nil
}

func (r *DataproxyRefresher) UpdateProfile(
	profile model.Profile,
	profileJSON map[string]interface{},
) error {
	// Check if the profile already exists.
	count, err := r.svc.Count(profile.Oid)
	if err != nil {
		return fmt.Errorf(
			"can't count profile. Profile ID: %s, error: %v",
			profile.Oid,
			err,
		)
	}

	if count <= 0 {
		// If profile doesn't exist, create a new one.
		profileJSON["cuid"] = cuid.New()
		if err := r.svc.Add(profileJSON); err != nil {
			return fmt.Errorf(
				"can't add a profile. Profile ID: %s, error: %v",
				profile.Oid,
				err,
			)
		}
	} else {
		// If profile already exists, update it.
		result, err := r.svc.Update(profile.Oid, profileJSON)
		if err != nil {
			return fmt.Errorf(
				"can't update a profile. Profile ID: %s, error: %v",
				profile.Oid, err,
			)
		}
		profileJSON["cuid"] = result["cuid"]
	}

	// Post the update to the Index
	postNodeURL := config.Conf.Index.URL + APINodesPath
	profileURL := config.Conf.DataProxy.URL + APIProfilesPath + "/" +
		profileJSON["cuid"].(string)
	nodeID, err := importutil.PostIndex(postNodeURL, profileURL)
	if err != nil {
		return fmt.Errorf(
			"failed to post profile to Index. Profile URL: %s, error: %v",
			profileURL,
			err,
		)
	}

	// Save node_id to profile
	if err := r.svc.UpdateNodeID(profile.Oid, nodeID); err != nil {
		return fmt.Errorf(
			"update node id failed. Profile ID: %s, error: %v",
			profile.Oid,
			err,
		)
	}

	return nil
}

func (r *DataproxyRefresher) updateProfileAccessTime(
	profile model.Profile,
) error {
	err := r.svc.UpdateAccessTime(profile.Oid)
	if err != nil {
		return fmt.Errorf(
			"failed to update access time for Profile CUID %s: %w",
			profile.Cuid,
			err,
		)
	}
	return nil
}

func (r *DataproxyRefresher) deleteProfile(profile model.Profile) error {
	err := r.svc.Delete(profile.Cuid)
	if err != nil {
		return fmt.Errorf(
			"failed to delete profile from datastore with CUID %s: %w",
			profile.Cuid,
			err,
		)
	}

	// Delete from the Index service.
	deleteNodeURL := config.Conf.Index.URL + APINodesPath + "/" + profile.NodeID

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, deleteNodeURL, nil)
	if err != nil {
		return fmt.Errorf(
			"failed to create DELETE request for Index service with NodeID %s: %w",
			profile.NodeID,
			err,
		)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(
			"failed to execute DELETE request for Index service with NodeID %s: %w",
			profile.NodeID,
			err,
		)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var resBody map[string]interface{}
		if decodeErr := json.NewDecoder(res.Body).Decode(&resBody); decodeErr != nil {
			return fmt.Errorf(
				"failed to decode error response from Index service: %w",
				decodeErr,
			)
		}

		if errors, exists := resBody["errors"].([]interface{}); exists {
			var errorMessages []string
			for _, item := range errors {
				errorMessages = append(errorMessages, fmt.Sprintf("%v", item))
			}
			return fmt.Errorf(
				"failed to delete from Index service with NodeID %s, errors: %s",
				profile.NodeID,
				strings.Join(errorMessages, ", "),
			)
		}

		return fmt.Errorf(
			"failed to delete from Index service with NodeID %s, unknown error",
			profile.NodeID,
		)
	}

	return nil
}

func (r *DataproxyRefresher) cleanUp() {
	mongodb.Client.Disconnect()
}
