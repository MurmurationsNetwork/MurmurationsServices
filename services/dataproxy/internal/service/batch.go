package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/lucsky/cuid"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/profile/profilevalidator"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/validatenode"
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

func (s *batchService) Validate(
	schemas []string,
	records [][]string,
) (int, []jsonapi.Error, error) {
	if len(records) > 1001 {
		return -1, nil, errors.New(
			"the CSV file cannot contain more than 1,000 rows",
		)
	}

	rawProfiles := csvToMap(records)

	// parse schemas to json string schema from library and validate
	validateJSONSchemas, _, err := parseSchemas(schemas)
	if err != nil {
		return -1, nil, err
	}

	for line, rawProfile := range rawProfiles {
		profile, err := mapToProfile(rawProfile, schemas)
		if err != nil {
			return line, nil, err
		}

		validator, err := profilevalidator.NewBuilder().
			WithMapProfile(profile).
			WithStrSchemas(validateJSONSchemas).
			WithCustomValidation().
			Build()
		if err != nil {
			return -1, nil, err
		}

		result := validator.Validate()
		if !result.Valid {
			return line, jsonapi.NewError(
				result.ErrorMessages,
				result.Details,
				result.Sources,
				result.ErrorStatus,
			), nil
		}
	}

	return -1, nil, nil
}

func (s *batchService) Import(
	title string,
	schemas []string,
	records [][]string,
	userID string,
	metaName string,
	metaURL string,
) (string, int, []jsonapi.Error, error) {
	if len(records) > 1001 {
		return "", -1, nil, errors.New(
			"the CSV file cannot contain more than 1,000 rows",
		)
	}

	// Generate `batch_id` using cuid and save it to MongoDB
	batchID := cuid.New()
	err := s.batchRepo.SaveUser(userID, title, batchID, schemas)
	if err != nil {
		return batchID, -1, nil, err
	}

	rawProfiles := csvToMap(records)

	// parse schemas to json string schema from library and validate
	validateJSONSchemas, validateSchemas, err := parseSchemas(schemas)
	if err != nil {
		return batchID, -1, nil, err
	}

	for line, rawProfile := range rawProfiles {
		profile, err := mapToProfile(rawProfile, schemas)
		if err != nil {
			return batchID, line, nil, err
		}

		// Validate data and if needed, respond with error
		result := validatenode.ValidateAgainstSchemasWithoutURL(
			validateJSONSchemas,
			validateSchemas,
			profile,
		)
		if !result.Valid {
			return batchID, line, jsonapi.NewError(
				result.ErrorMessages,
				result.Details,
				result.Sources,
				result.ErrorStatus,
			), nil
		}

		// TODO
		profileHash, err := jsonutil.Hash(profile)
		if err != nil {
			return batchID, line, nil, err
		}
		profile["source_data_hash"] = profileHash

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

		profile["batch_id"] = batchID

		// Generate cuid for profile
		profileCuid := cuid.New()
		profile["cuid"] = profileCuid

		// Import profile to MongoDB
		err = s.batchRepo.SaveProfile(profile)
		if err != nil {
			return batchID, line, nil, err
		}

		// Import profile to Index
		postNodeURL := config.Conf.Index.URL + "/v2/nodes"
		profileURL := config.Conf.DataProxy.URL + "/v1/profiles/" + profileCuid
		nodeID, err := importutil.PostIndex(postNodeURL, profileURL)
		if err != nil {
			return batchID, line, nil, errors.New(
				"Import to MurmurationsServices Index failed: " + err.Error(),
			)
		}

		// Save `node_id` to MongoDB
		profile["node_id"] = nodeID
		profile["is_posted"] = true
		err = s.batchRepo.SaveNodeID(profileCuid, profile)
		if err != nil {
			return batchID, line, nil, errors.New(
				"Save node_id to Mongo failed: " + err.Error(),
			)
		}
	}

	return batchID, -1, nil, nil
}

func (s *batchService) Edit(
	title string,
	records [][]string,
	userID string,
	batchID string,
	metaName string,
	metaURL string,
) (int, []jsonapi.Error, error) {
	if len(records) > 1001 {
		return -1, nil, errors.New(
			"the CSV file cannot contain more than 1,000 rows",
		)
	}

	// Check if `batch_id` belongs to user
	isValid, err := s.batchRepo.CheckUser(userID, batchID)
	if err != nil {
		return -1, nil, err
	}
	if !isValid {
		return -1, nil, errors.New(
			"the `batch_id` doesn't belong to the specified user",
		)
	}

	// save current schemas to batch collection
	schemas, err := s.batchRepo.UpdateBatchTitle(title, batchID)
	if err != nil {
		return -1, nil, err
	}

	// Get profile `oid`, cuid and hash by `batch_id`
	profileOidsAndHashes, err := s.batchRepo.GetProfileOidsAndHashesByBatchID(
		batchID,
	)
	if err != nil {
		return -1, nil, err
	}

	rawProfiles := csvToMap(records)

	// parse schemas to json string schema from library and validate
	validateJSONSchemas, validateSchemas, err := parseSchemas(schemas)
	if err != nil {
		return -1, nil, err
	}

	for line, rawProfile := range rawProfiles {
		profile, err := mapToProfile(rawProfile, schemas)
		if err != nil {
			return line, nil, err
		}

		// Validate data and if needed, respond with error
		result := validatenode.ValidateAgainstSchemasWithoutURL(
			validateJSONSchemas,
			validateSchemas,
			profile,
		)
		if !result.Valid {
			return line, jsonapi.NewError(
				result.ErrorMessages,
				result.Details,
				result.Sources,
				result.ErrorStatus,
			), nil
		}

		// TODO
		profileHash, err := jsonutil.Hash(profile)
		if err != nil {
			return line, nil, err
		}
		profile["source_data_hash"] = profileHash
		profile["batch_id"] = batchID

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

		// Check if profile exists in MongoDB
		oid := profile["oid"].(string)
		_, ok := profileOidsAndHashes[oid]
		var profileCuid string
		if ok {
			profileCuid = profileOidsAndHashes[oid][0]
			// If current profile's `oid` and `profile_hash` match the data in MongoDB, skip it
			if profileOidsAndHashes[oid][1] == profileHash {
				delete(profileOidsAndHashes, oid)
				continue
			}
			// Otherwise update the profile in MongoDB
			err = s.batchRepo.UpdateProfile(profileCuid, profile)
			if err != nil {
				return line, nil, err
			}
			// Delete `oid` from profileOidsAndHashes, so that the rest of data in it needs to be deleted later
			delete(profileOidsAndHashes, oid)
		} else {
			// If profile doesn't have cuid, generate one
			profileCuid = cuid.New()
			profile["cuid"] = profileCuid

			// Import profile to MongoDB
			err = s.batchRepo.SaveProfile(profile)
			if err != nil {
				return line, nil, err
			}
		}

		// Import profile to Index
		postNodeURL := config.Conf.Index.URL + "/v2/nodes"
		profileURL := config.Conf.DataProxy.URL + "/v1/profiles/" + profileCuid
		nodeID, err := importutil.PostIndex(postNodeURL, profileURL)
		if err != nil {
			return line, nil, errors.New(
				"Import to Index failed: " + err.Error(),
			)
		}

		// Save `node_id` to MongoDB
		profile["node_id"] = nodeID
		profile["is_posted"] = true
		err = s.batchRepo.SaveNodeID(profileCuid, profile)
		if err != nil {
			return line, nil, errors.New(
				"Save node_id to MongoDB failed: " + err.Error(),
			)
		}
	}

	// The rest of the profiles which are not in the CSV file need to be deleted
	if len(profileOidsAndHashes) > 0 {
		for _, cuidAndHash := range profileOidsAndHashes {
			// Get profile by cuid
			profile, err := s.batchRepo.GetProfileByCuid(cuidAndHash[0])
			if err != nil {
				return -1, nil, err
			}

			// Delete profiles from mongo
			err = s.batchRepo.DeleteProfileByCuid(cuidAndHash[0])
			if err != nil {
				return -1, nil, err
			}

			// Delete profiles from Index
			if profile["is_posted"].(bool) {
				nodeID := profile["node_id"].(string)
				deleteNodeURL := config.Conf.Index.URL + "/v2/nodes/" + nodeID
				err := importutil.DeleteIndex(deleteNodeURL, nodeID)
				if err != nil {
					return -1, nil, errors.New(
						"failed to delete from Index : " + err.Error(),
					)
				}
			}
		}
	}

	return -1, nil, nil
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

	// Get profiles by batch_id
	profiles, err := s.batchRepo.GetProfilesByBatchID(batchID)
	if err != nil {
		return err
	}

	// Delete profiles from MongoDB
	err = s.batchRepo.DeleteProfilesByBatchID(batchID)
	if err != nil {
		return err
	}

	// Delete profiles from Index
	for _, profile := range profiles {
		if profile["is_posted"].(bool) {
			nodeID := profile["node_id"].(string)
			deleteNodeURL := config.Conf.Index.URL + "/v2/nodes/" + nodeID
			err := importutil.DeleteIndex(deleteNodeURL, nodeID)
			if err != nil {
				return errors.New(
					"Delete profile from MurmurationsServices Index failed: " + err.Error(),
				)
			}
		}
	}

	// Delete batch_id from MongoDB
	err = s.batchRepo.DeleteBatchID(batchID)
	if err != nil {
		return err
	}

	return nil
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

	fmt.Println(profile)

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
	// Check if field name ends with "(list)" and adjust processing accordingly
	isList := false
	if strings.HasSuffix(field, "(list)") {
		// Remove the "(list)" suffix
		field = strings.TrimSuffix(field, "(list)")
		isList = true
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
				current[p] = append(current[p].([]interface{}), destructValue(v))
			}
			break
		}

		// If the next path is a number, and it's the last element, it means it's an array
		if i == len(path)-2 {
			if _, err := strconv.Atoi(path[i+1]); err == nil {
				if _, ok := current[path[i]]; !ok {
					current[path[i]] = make([]interface{}, 0)
				}
				current[path[i]] = append(
					current[path[i]].([]interface{}),
					destructValue(value),
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

func parseSchemas(schemas []string) ([]string, []string, error) {
	// validate against the default schema.
	validateSchemas := []string{"default-v2.0.0"}
	validateSchemas = append(validateSchemas, schemas...)
	validateJSONSchemas := make([]string, len(validateSchemas))
	libraryURL := config.Conf.Library.InternalURL + "/v2/schemas"
	for i, schema := range validateSchemas {
		res, err := http.Get(libraryURL + "/" + schema)
		if err != nil {
			return nil, nil, err
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, nil, err
		}
		validateJSONSchemas[i] = string(body)
	}
	return validateJSONSchemas, validateSchemas, nil
}

func splitEscapedComma(s string) []string {
	var result []string
	var current strings.Builder
	for i := 0; i < len(s); i++ {
		// Split the string by comma, but if the comma is at the first one or the comma is escaped, ignore it
		if s[i] == ',' && (i == 0 || s[i-1] != '\\') {
			result = append(result, current.String())
			current.Reset()
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
