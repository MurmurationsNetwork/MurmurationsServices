package service

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/lucsky/cuid"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/profile/profilevalidator"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/repository/mongo"
)

type BatchService interface {
	GetBatchesByUserID(string) ([]model.Batch, error)
	Validate([]string, [][]string) (int, []jsonapi.Error, error)
	Import(
		string,
		[]string,
		[][]string,
		string,
		string,
		string,
	) (string, int, []jsonapi.Error, error)
	Edit(
		string,
		[][]string,
		string,
		string,
		string,
		string,
	) (int, []jsonapi.Error, error)
	Delete(string, string) error
}

type batchService struct {
	batchRepo mongo.BatchRepository
}

func NewBatchService(batchRepo mongo.BatchRepository) BatchService {
	return &batchService{
		batchRepo: batchRepo,
	}
}

func (s *batchService) GetBatchesByUserID(
	userID string,
) ([]model.Batch, error) {
	batches, err := s.batchRepo.GetBatchesByUserID(userID)
	if err != nil {
		return nil, err
	}

	return batches, nil
}

// Validate processes and validates CSV records against the provided schemas.
// It returns the line number where an error occurred (if any), a slice of
// jsonapi.Errors containing validation errors, and an error if a non-validation
// issue occurs.
//
// Parameters:
//   - schemaNames: A slice of schema names to validate against.
//   - csvRecords: CSV records represented as a slice of string slices.
//
// Returns:
//   - int: The line number where an error occurred (-1 if not applicable).
//   - []jsonapi.Error: A slice of validation errors (if any).
//   - error: An error object if a non-validation error occurred.
func (s *batchService) Validate(
	schemaNames []string,
	csvRecords [][]string,
) (int, []jsonapi.Error, error) {
	// Maximum number of data rows allowed in the CSV (excluding header).
	const maxDataRows = 1000

	// Check if the CSV exceeds the maximum allowed data rows.
	if len(csvRecords) > maxDataRows+1 { // +1 for header row
		return -1, nil, fmt.Errorf(
			"the CSV file cannot contain more than %d data rows",
			maxDataRows,
		)
	}

	_, validationErrors, err := s.validateProfiles(csvRecords, schemaNames)
	if err != nil {
		return -1, nil, err
	}

	// Return collected validation errors if any.
	if len(validationErrors) > 0 {
		return -1, validationErrors, nil
	}

	// All profiles are valid.
	return -1, nil, nil
}

// Import processes CSV records by first validating all profiles and then importing them
// into the system if they pass validation. It returns the generated batch ID, the line
// number where an error occurred (-1 if not applicable), a slice of jsonapi.Errors containing
// validation errors, and an error if a non-validation issue occurs.
//
// Parameters:
//   - title: The title for the batch import.
//   - schemaNames: A slice of schema names to validate against.
//   - csvRecords: CSV records represented as a slice of string slices.
//   - userID: The ID of the user performing the import.
//   - metaName: Metadata name to be added to each profile (optional).
//   - metaURL: Metadata URL to be added to each profile (optional).
//
// Returns:
//   - string: The generated batch ID.
//   - int: The line number where an error occurred (-1 if not applicable).
//   - []jsonapi.Error: A slice of validation errors (if any).
//   - error: An error object if a non-validation error occurred.
func (s *batchService) Import(
	title string,
	schemaNames []string,
	records [][]string,
	userID string,
	metaName string,
	metaURL string,
) (string, int, []jsonapi.Error, error) {
	// Maximum number of data rows allowed in the CSV (excluding header).
	const maxDataRows = 1000

	// Check if the CSV exceeds the maximum allowed data rows.
	if len(records) > maxDataRows+1 { // +1 for header row
		return "", -1, nil, fmt.Errorf(
			"the CSV file cannot contain more than %d data rows",
			maxDataRows,
		)
	}

	// Generate a new batch ID.
	batchID := cuid.New()

	mappedProfiles, validationErrors, err := s.validateProfiles(
		records,
		schemaNames,
	)

	if err != nil {
		return batchID, -1, nil, err
	}

	// If there are any validation errors, return them without proceeding further.
	if len(validationErrors) > 0 {
		return batchID, -1, validationErrors, nil
	}

	// Save the batch information to the database.
	err = s.batchRepo.SaveUser(
		userID,
		title,
		batchID,
		schemaNames,
		len(mappedProfiles),
	)
	if err != nil {
		return batchID, -1, nil, err
	}

	// Second pass: Process valid profiles and save them to the database.
	go s.ProcessImportAsync(batchID, mappedProfiles, metaName, metaURL)

	// All profiles have been successfully imported.
	return batchID, -1, nil, nil
}

func (s *batchService) ProcessImportAsync(
	batchID string,
	validProfiles []map[string]interface{},
	metaName string,
	metaURL string,
) {
	err := s.batchRepo.UpdateBatchStatus(batchID, "processing")
	if err != nil {
		logger.Error("Failed to update batch status to processing", err)
		return
	}

	processedNodes := 0
	totalNodes := len(validProfiles)

	for i, mappedProfile := range validProfiles {
		// Generate a new CUID for the profile and set "cuid".
		profileCUID := cuid.New()
		mappedProfile["cuid"] = profileCUID

		// Compute a hash of the profile and store it in "source_data_hash".
		profileHash, err := jsonutil.Hash(mappedProfile)
		if err != nil {
			batchErr := s.batchRepo.UpdateBatchError(
				batchID,
				"failed",
				fmt.Sprintf("Failed to compute hash of profile: %v", err),
			)
			if batchErr != nil {
				logger.Error("Failed to update batch error", batchErr)
			}
			break
		}
		mappedProfile["source_data_hash"] = profileHash

		// Add metadata if provided.
		if metaName != "" || metaURL != "" {
			sourceInfo := make(map[string]interface{})
			if metaName != "" {
				sourceInfo["name"] = metaName
			}
			if metaURL != "" {
				sourceInfo["url"] = metaURL
			}
			metadata := map[string]interface{}{
				"sources": []map[string]interface{}{
					sourceInfo,
				},
			}
			mappedProfile["metadata"] = metadata
		}

		// Set "batch_id" to associate the profile with the batch.
		mappedProfile["batch_id"] = batchID

		// Save the profile to the database.
		err = s.batchRepo.SaveProfile(mappedProfile)
		if err != nil {
			batchErr := s.batchRepo.UpdateBatchError(
				batchID,
				"failed",
				fmt.Sprintf("Failed to save profile: %v", err),
			)
			if batchErr != nil {
				logger.Error("Failed to update batch error", batchErr)
			}
			break
		}

		// Post the profile to the index service.
		postNodeURL := config.Values.Index.URL + "/v2/nodes"
		profileURL := config.Values.DataProxy.URL + "/v1/profiles/" + profileCUID
		nodeID, err := importutil.PostIndex(postNodeURL, profileURL)
		if err != nil {
			batchErr := s.batchRepo.UpdateBatchError(
				batchID,
				"failed",
				fmt.Sprintf("Failed to post profile to index service: %v", err),
			)
			if batchErr != nil {
				logger.Error("Failed to update batch error", batchErr)
			}
			break
		}

		// Update the profile with the node ID and mark it as posted.
		mappedProfile["node_id"] = nodeID
		mappedProfile["is_posted"] = true

		// Save the node ID to the database.
		err = s.batchRepo.SaveNodeID(profileCUID, mappedProfile)
		if err != nil {
			batchErr := s.batchRepo.UpdateBatchError(
				batchID,
				"failed",
				fmt.Sprintf("Failed to save node ID to database: %v", err),
			)
			if batchErr != nil {
				logger.Error("Failed to update batch error", batchErr)
			}
			break
		}

		processedNodes++

		if processedNodes%10 == 0 || i == totalNodes-1 {
			batchErr := s.batchRepo.UpdateBatchProgress(batchID, processedNodes)
			if batchErr != nil {
				logger.Error("Failed to update batch progress", batchErr)
			}
		}
	}

	err = s.batchRepo.UpdateBatchStatus(batchID, "completed")
	if err != nil {
		logger.Error("Failed to update batch status to completed", err)
		return
	}
}

func (s *batchService) Edit(
	title string,
	records [][]string,
	userID string,
	batchID string,
	metaName string,
	metaURL string,
) (int, []jsonapi.Error, error) {
	const maxDataRows = 1000

	if len(records) > maxDataRows+1 {
		return -1, nil, fmt.Errorf(
			"the CSV file cannot contain more than %d data rows",
			maxDataRows,
		)
	}

	// Check if `batch_id` belongs to user
	isValid, err := s.batchRepo.CheckUser(userID, batchID)
	if err != nil {
		return -1, nil, err
	}
	if !isValid {
		return -1, nil, fmt.Errorf(
			"the `batch_id` doesn't belong to the specified user",
		)
	}

	batchStatus, err := s.batchRepo.GetBatchStatus(batchID)
	if err != nil {
		return -1, nil, err
	}
	if batchStatus != "completed" {
		return -1, nil, fmt.Errorf(
			"batch is not completed. Please wait for the batch to be completed before editing it",
		)
	}

	// save current schemas to batch collection
	schemas, err := s.batchRepo.UpdateBatchTitle(title, batchID)
	if err != nil {
		return -1, nil, err
	}

	mappedProfiles, validationErrors, err := s.validateProfiles(
		records,
		schemas,
	)

	if err != nil {
		return -1, nil, err
	}

	if len(validationErrors) > 0 {
		return -1, validationErrors, nil
	}

	go s.ProcessEditAsync(
		batchID,
		mappedProfiles,
		metaName,
		metaURL,
	)

	return -1, nil, nil
}

func (s *batchService) validateProfiles(
	records [][]string,
	schemas []string,
) ([]map[string]interface{}, []jsonapi.Error, error) {
	rawProfiles := csvToMap(records)

	parsed, err := ParseSchemas(schemas)
	if err != nil {
		return nil, nil, err
	}

	jsonSchemas := parsed.JSONSchemas
	schemaNames := parsed.SchemaNames

	var validationErrors []jsonapi.Error
	mappedProfiles := make([]map[string]interface{}, 0, len(rawProfiles))

	for line, rawProfile := range rawProfiles {
		profile, err := mapToProfile(rawProfile, schemas)
		if err != nil {
			return nil, nil, fmt.Errorf("line %d: %v", line, err)
		}

		validator, err := profilevalidator.NewBuilder().
			WithMapProfile(profile).
			WithJSONSchemas(schemaNames, jsonSchemas).
			Build()
		if err != nil {
			return nil, nil, err
		}

		result := validator.Validate()

		for idx := range result.Sources {
			result.Sources[idx] = append(
				result.Sources[idx],
				"oid",
				rawProfile["oid"],
			)
		}

		if !result.Valid {
			validationErrors = append(validationErrors,
				jsonapi.NewError(
					result.ErrorMessages,
					result.Details,
					result.Sources,
					result.ErrorStatus,
				)...,
			)
			continue
		}

		mappedProfiles = append(mappedProfiles, profile)
	}

	if len(validationErrors) > 0 {
		return nil, validationErrors, nil
	}

	return mappedProfiles, nil, nil
}

func (s *batchService) ProcessEditAsync(
	batchID string,
	mappedProfiles []map[string]interface{},
	metaName string,
	metaURL string,
) {
	// update status = processing
	err := s.batchRepo.UpdateBatchStatus(batchID, "editing")
	if err != nil {
		logger.Error("Failed to update batch status to editing", err)
		return
	}

	profileOidsAndHashes, err := s.batchRepo.GetProfileOidsAndHashesByBatchID(
		batchID,
	)
	if err != nil {
		batchErr := s.batchRepo.UpdateBatchError(
			batchID,
			"edit_failed",
			err.Error(),
		)
		if batchErr != nil {
			logger.Error("Failed to update batch error", batchErr)
		}
		return
	}

	totalNodes := len(profileOidsAndHashes)
	err = s.batchRepo.UpdateBatchTotalNodesAndProgress(batchID, totalNodes, 0)
	if err != nil {
		batchErr := s.batchRepo.UpdateBatchError(
			batchID,
			"edit_failed",
			err.Error(),
		)
		if batchErr != nil {
			logger.Error("Failed to update batch error", batchErr)
		}
		return
	}

	processed := 0

	for _, profile := range mappedProfiles {
		// Add metadata
		if metaName != "" || metaURL != "" {
			source := make(map[string]interface{})
			if metaName != "" {
				source["name"] = metaName
			}
			if metaURL != "" {
				source["url"] = metaURL
			}

			metadata := map[string]interface{}{
				"sources": []map[string]interface{}{
					source,
				},
			}
			profile["metadata"] = metadata
		}

		profileHash, err := jsonutil.Hash(profile)
		if err != nil {
			batchErr := s.batchRepo.UpdateBatchError(
				batchID,
				"edit_failed",
				err.Error(),
			)
			if batchErr != nil {
				logger.Error("Failed to update batch error", batchErr)
			}
			return
		}
		profile["source_data_hash"] = profileHash
		profile["batch_id"] = batchID

		oid := profile["oid"].(string)
		var profileCUID string
		_, ok := profileOidsAndHashes[oid]
		if ok {
			profileCUID = profileOidsAndHashes[oid][0]
			// If current profile's `oid` and `profile_hash` match the data in MongoDB, skip it
			if profileOidsAndHashes[oid][1] == profileHash {
				delete(profileOidsAndHashes, oid)
				processed++
				if processed%10 == 0 || len(mappedProfiles) == processed {
					err = s.batchRepo.UpdateBatchProgress(batchID, processed)
					if err != nil {
						batchErr := s.batchRepo.UpdateBatchError(
							batchID,
							"edit_failed",
							err.Error(),
						)
						if batchErr != nil {
							logger.Error(
								"Failed to update batch error",
								batchErr,
							)
						}
						return
					}
				}
				continue
			}
			// Otherwise update the profile in MongoDB
			err = s.batchRepo.UpdateProfile(profileCUID, profile)
			if err != nil {
				batchErr := s.batchRepo.UpdateBatchError(
					batchID,
					"edit_failed",
					err.Error(),
				)
				if batchErr != nil {
					logger.Error("Failed to update batch error", batchErr)
				}
				return
			}
			// Delete `oid` from profileOidsAndHashes, so that the rest of data in it needs to be deleted later
			delete(profileOidsAndHashes, oid)
		} else {
			profileCUID = cuid.New()
			profile["cuid"] = profileCUID
			err = s.batchRepo.SaveProfile(profile)
			if err != nil {
				batchErr := s.batchRepo.UpdateBatchError(batchID, "edit_failed", err.Error())
				if batchErr != nil {
					logger.Error("Failed to update batch error", batchErr)
				}
				return
			}
		}

		// post to index
		postNodeURL := config.Values.Index.URL + "/v2/nodes"
		profileURL := config.Values.DataProxy.URL + "/v1/profiles/" + profileCUID
		nodeID, err := importutil.PostIndex(postNodeURL, profileURL)
		if err != nil {
			batchErr := s.batchRepo.UpdateBatchError(
				batchID,
				"edit_failed",
				err.Error(),
			)
			if batchErr != nil {
				logger.Error("Failed to update batch error", batchErr)
			}
			return
		}

		profile["node_id"] = nodeID
		profile["is_posted"] = true
		err = s.batchRepo.SaveNodeID(profileCUID, profile)
		if err != nil {
			batchErr := s.batchRepo.UpdateBatchError(
				batchID,
				"edit_failed",
				err.Error(),
			)
			if batchErr != nil {
				logger.Error("Failed to update batch error", batchErr)
			}
			return
		}

		processed++
		if processed%10 == 0 || len(mappedProfiles) == processed {
			err = s.batchRepo.UpdateBatchProgress(batchID, processed)
			if err != nil {
				batchErr := s.batchRepo.UpdateBatchError(
					batchID,
					"edit_failed",
					err.Error(),
				)
				if batchErr != nil {
					logger.Error("Failed to update batch error", batchErr)
				}
				return
			}
		}
	}

	// The rest of the profiles which are not in the CSV file need to be deleted
	if len(profileOidsAndHashes) > 0 {
		for _, cuidAndHash := range profileOidsAndHashes {
			// Get profile by cuid
			profile, err := s.batchRepo.GetProfileByCuid(cuidAndHash[0])
			if err != nil {
				batchErr := s.batchRepo.UpdateBatchError(
					batchID,
					"edit_failed",
					err.Error(),
				)
				if batchErr != nil {
					logger.Error("Failed to update batch error", batchErr)
				}
				return
			}

			// Delete profiles from mongo
			err = s.batchRepo.DeleteProfileByCuid(cuidAndHash[0])
			if err != nil {
				batchErr := s.batchRepo.UpdateBatchError(
					batchID,
					"edit_failed",
					err.Error(),
				)
				if batchErr != nil {
					logger.Error("Failed to update batch error", batchErr)
				}
				return
			}

			// Delete profiles from Index
			if profile["is_posted"].(bool) {
				nodeID := profile["node_id"].(string)
				deleteNodeURL := config.Values.Index.URL + "/v2/nodes/" + nodeID
				err := importutil.DeleteIndex(deleteNodeURL, nodeID)
				if err != nil {
					batchErr := s.batchRepo.UpdateBatchError(
						batchID,
						"edit_failed",
						err.Error(),
					)
					if batchErr != nil {
						logger.Error("Failed to update batch error", batchErr)
					}
					return
				}
			}

			processed++
			if processed%10 == 0 || totalNodes == processed {
				err = s.batchRepo.UpdateBatchProgress(batchID, processed)
				if err != nil {
					batchErr := s.batchRepo.UpdateBatchError(
						batchID,
						"edit_failed",
						err.Error(),
					)
					if batchErr != nil {
						logger.Error("Failed to update batch error", batchErr)
					}
					return
				}
			}
		}
	}

	err = s.batchRepo.UpdateBatchStatus(batchID, "completed")
	if err != nil {
		logger.Error("Failed to update batch status to completed", err)
		return
	}
}

func (s *batchService) Delete(userID string, batchID string) error {
	// Check if `batch_id` belongs to user
	isValid, err := s.batchRepo.CheckUser(userID, batchID)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.New("batch_id doesn't belong to user")
	}

	batchStatus, err := s.batchRepo.GetBatchStatus(batchID)
	if err != nil {
		return err
	}
	if batchStatus != "completed" {
		return errors.New(
			"batch is not completed. Please wait for the batch to be completed before deleting it",
		)
	}

	go s.deleteBatchAsync(batchID)

	return nil
}

func (s *batchService) deleteBatchAsync(batchID string) {
	err := s.batchRepo.UpdateBatchStatus(batchID, "deleting")
	if err != nil {
		logger.Error("Failed to update batch status to deleting", err)
		return
	}

	// Get profiles by batch_id
	profiles, err := s.batchRepo.GetProfilesByBatchID(batchID)
	if err != nil {
		batchErr := s.batchRepo.UpdateBatchError(
			batchID,
			"delete_failed",
			err.Error(),
		)
		if batchErr != nil {
			logger.Error("Failed to update batch error", batchErr)
		}
		return
	}

	totalNodes := len(profiles)
	err = s.batchRepo.UpdateBatchTotalNodesAndProgress(batchID, totalNodes, 0)
	if err != nil {
		batchErr := s.batchRepo.UpdateBatchError(
			batchID,
			"delete_failed",
			err.Error(),
		)
		if batchErr != nil {
			logger.Error("Failed to update batch error", batchErr)
		}
		return
	}

	processed := 0

	// Delete profiles from MongoDB
	err = s.batchRepo.DeleteProfilesByBatchID(batchID)
	if err != nil {
		batchErr := s.batchRepo.UpdateBatchError(
			batchID,
			"delete_failed",
			err.Error(),
		)
		if batchErr != nil {
			logger.Error("Failed to update batch error", batchErr)
		}
		return
	}

	// Delete profiles from Index
	for _, profile := range profiles {
		if profile["is_posted"].(bool) {
			nodeID := profile["node_id"].(string)
			deleteNodeURL := config.Values.Index.URL + "/v2/nodes/" + nodeID
			err := importutil.DeleteIndex(deleteNodeURL, nodeID)
			if err != nil {
				batchErr := s.batchRepo.UpdateBatchError(
					batchID,
					"delete_failed",
					err.Error(),
				)
				if batchErr != nil {
					logger.Error("Failed to update batch error", batchErr)
				}
				return
			}
		}

		processed++
		if processed%10 == 0 || totalNodes == processed {
			err = s.batchRepo.UpdateBatchProgress(batchID, processed)
			if err != nil {
				batchErr := s.batchRepo.UpdateBatchError(
					batchID,
					"delete_failed",
					err.Error(),
				)
				if batchErr != nil {
					logger.Error("Failed to update batch error", batchErr)
				}
				return
			}
		}
	}

	// Delete batch_id from MongoDB
	err = s.batchRepo.DeleteBatchID(batchID)
	if err != nil {
		batchErr := s.batchRepo.UpdateBatchError(
			batchID,
			"delete_failed",
			err.Error(),
		)
		if batchErr != nil {
			logger.Error("Failed to update batch error", batchErr)
		}
		return
	}
}

// Convert csv to one-to-one map[string]string.
func csvToMap(records [][]string) []map[string]string {
	csvHeader := records[0]
	var rawProfiles []map[string]string

	for i := 1; i < len(records); i++ {
		rawProfile := make(map[string]string)
		for index, value := range records[i] {
			if value != "" {
				rawProfile[csvHeader[index]] = value
			}
		}
		rawProfiles = append(rawProfiles, rawProfile)
	}
	return rawProfiles
}

// Convert one-to-one map[string]string to profile data structure.
func mapToProfile(
	rawProfile map[string]string,
	schemas []string,
) (map[string]interface{}, error) {
	// Sort rawProfile by key
	keys := make([]string, 0, len(rawProfile))
	// Validate `oid` field
	hasOid := false
	for k := range rawProfile {
		if k == "oid" {
			hasOid = true
		}
		keys = append(keys, k)
	}

	if !hasOid {
		return nil, errors.New("the `oid` field is required")
	}

	// Sort the keys
	sort.Strings(keys)
	profile := make(map[string]interface{})
	var err error
	for _, k := range keys {
		value := rawProfile[k]
		if k == "oid" {
			profile["oid"] = value
			continue
		}
		if value != "" {
			profile, err = destructField(profile, k, value)
			if err != nil {
				return nil, err
			}
		}
	}

	// Put schema here
	profile["linked_schemas"] = schemas
	return profile, nil
}

// Destructure field name and save field value to profile data structure.
func destructField(
	profile map[string]interface{},
	field string,
	value string,
) (map[string]interface{}, error) {
	// Handle (list) fields - issue #727
	// Validating the field name
	if strings.Contains(field, "(") {
		// Use regex to validate the field name
		regexPattern := `^.+\(list-\d+\)$`
		re, err := regexp.Compile(regexPattern)
		if err != nil {
			return nil, err
		}
		if !re.MatchString(field) {
			return nil, errors.New(
				"field format error: please use (list-number) format for lists",
			)
		}
	}

	// Check if field name with multiple (list) - e.g., tags(list-0), tags(list-1)
	isList := false
	lastLeftParenIndex := strings.LastIndex(field, "(")
	if lastLeftParenIndex != -1 && strings.HasSuffix(field, ")") {
		isList = true
		field = field[:lastLeftParenIndex]
	}

	// Destructure field name
	// e.g., "urls[0].name" -> ["urls", 0, "name"], "tags[0]" -> ["tags", 0]
	fieldName := strings.Split(field, ".")
	var path []string
	for _, p := range fieldName {
		if i := strings.IndexByte(p, '['); i != -1 {
			index := strings.Trim(p[i:], "[]")
			path = append(path, p[:i], index)
		} else {
			path = append(path, p)
		}
	}

	// Iterate profile to put value into the correct path
	current := profile
	for i, p := range path {
		// If the current path is a number, skip it, because it's already handled in the previous loop
		if _, err := strconv.Atoi(p); err == nil {
			continue
		}

		// If the field is a list, and it's the last path, append the value to the array
		if isList && i == len(path)-1 {
			values := splitEscapedComma(value)
			_, exists := current[p]
			if !exists {
				current[p] = make([]interface{}, 0)
			}
			for _, v := range values {
				trimmedValue := strings.TrimSpace(v)
				current[p] = append(current[p].([]interface{}), trimmedValue)
			}
			break
		}

		// If the next path is a number, and it's the last element, it means it's an array
		if i == len(path)-2 {
			if _, err := strconv.Atoi(path[i+1]); err == nil {
				if _, ok := current[path[i]]; !ok {
					current[path[i]] = make([]interface{}, 0)
				}
				trimmedValue := strings.TrimSpace(value)
				current[path[i]] = append(
					current[path[i]].([]interface{}),
					trimmedValue,
				)
				break
			}
		}
		// If the next path is a number, it means it's an array-object
		if i+1 < len(path) {
			if arrayNum, err := strconv.Atoi(path[i+1]); err == nil {
				if _, ok := current[path[i]]; !ok {
					current[path[i]] = make([]map[string]interface{}, 0)
				}
				if _, ok := current[path[i]].([]map[string]interface{}); !ok {
					return nil, errors.New(
						"Check if the fields are duplicated or have different types of fields with the same name. Invalid field name: " + field,
					)
				}
				if len(
					current[path[i]].([]map[string]interface{}),
				) <= arrayNum {
					current[path[i]] = append(
						current[path[i]].([]map[string]interface{}),
						make(map[string]interface{}),
					)
				}
				if len(
					current[path[i]].([]map[string]interface{}),
				)-1 != arrayNum {
					return nil, errors.New(
						"Check the field name's array number is sequential and starts from 0. Invalid field name: " + field,
					)
				}
				current = current[path[i]].([]map[string]interface{})[arrayNum]
				continue
			}
		}
		// If the last element, put value into it.
		if i == len(path)-1 {
			current[p] = destructValue(value)
			break
		}
		if _, ok := current[p]; !ok {
			current[p] = make(map[string]interface{})
		}
		if _, ok := current[p].(map[string]interface{}); !ok {
			return nil, errors.New(
				"Check if the fields are duplicated or have different types of fields with the same name. Invalid field name: " + field,
			)
		}
		current = current[p].(map[string]interface{})
	}

	return profile, nil
}

func destructValue(value string) interface{} {
	// If the string has a dot, it means it's possibly a float number
	// If not, it's possible an int number
	// In all other cases, it's a string
	if strings.Contains(value, ".") {
		float, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return value
		}
		return float
	}
	integer, err := strconv.Atoi(value)
	if err != nil {
		return value
	}
	return integer
}

func splitEscapedComma(s string) []string {
	var result []string
	var current strings.Builder
	for i := 0; i < len(s); i++ {
		// Split the string by comma
		// If the first character is a comma, ignore it
		// If the current character is a comma and the previous character is not a backslash, split the string
		if s[i] == ',' && (i == 0 || s[i-1] != '\\') {
			if current.Len() > 0 || i != 0 {
				result = append(result, current.String())
				current.Reset()
			}
		} else if s[i] == ',' && i > 0 && s[i-1] == '\\' {
			// If the current character is a backslash and the next character is a comma, which means the comma is escaped
			current.WriteRune(',')
		} else if s[i] != '\\' {
			// Other characters are written to the current string
			current.WriteByte(s[i])
		}
	}
	result = append(result, current.String())
	return result
}
