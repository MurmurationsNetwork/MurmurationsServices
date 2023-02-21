package usecase

import (
	"errors"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/repository/db"
	"github.com/lucsky/cuid"
	"math"
	"sort"
	"strconv"
	"strings"
)

type BatchUsecase interface {
	Validate([]string, [][]string) (int, error)
	Import([]string, [][]string, string, string, string) (string, int, error)
	Edit([]string, [][]string, string, string, string, string) (int, error)
	Delete(string, string) error
}

type batchUsecase struct {
	batchRepo db.BatchRepository
}

func NewBatchService(batchRepo db.BatchRepository) BatchUsecase {
	return &batchUsecase{
		batchRepo: batchRepo,
	}
}

func (s *batchUsecase) Validate(schemas []string, records [][]string) (int, error) {
	if len(records) > 1001 {
		return -1, errors.New("CSV rows can't be larger than 1,000")
	}

	rawProfiles := csvToMap(records)

	for line, rawProfile := range rawProfiles {
		profile, err := mapToProfile(rawProfile, schemas)
		if err != nil {
			return line, err
		}

		// validate data, once it has error response error
		validateUrl := config.Conf.Index.URL + "/v2/validate"
		isValid, failureReasons, err := importutil.Validate(validateUrl, profile)
		if err != nil {
			return line, err
		}
		if !isValid {
			return line, errors.New(failureReasons)
		}
	}

	return -1, nil
}

func (s *batchUsecase) Import(schemas []string, records [][]string, userId string, metaName string, metaUrl string) (string, int, error) {
	if len(records) > 1001 {
		return "", -1, errors.New("CSV rows can't be larger than 1,000")
	}

	// generate batch_id import using cuid and save it to mongo
	batchId := cuid.New()
	err := s.batchRepo.SaveUser(userId, batchId)
	if err != nil {
		return batchId, -1, err
	}

	rawProfiles := csvToMap(records)

	for line, rawProfile := range rawProfiles {
		profile, err := mapToProfile(rawProfile, schemas)
		if err != nil {
			return batchId, line, err
		}

		// hash profile
		profileHash, err := importutil.HashProfile(profile)
		if err != nil {
			return batchId, line, err
		}
		profile["source_data_hash"] = profileHash

		// add metadata
		if metaName != "" && metaUrl != "" {
			metadata := map[string]interface{}{
				"sources": []map[string]interface{}{
					{
						"name": metaName,
						"url":  metaUrl,
					},
				},
			}
			profile["metadata"] = metadata
		}

		profile["batch_id"] = batchId

		// generate cuid for profile
		profileCuid := cuid.New()
		profile["cuid"] = profileCuid

		// import profile to Mongo
		err = s.batchRepo.SaveProfile(profile)
		if err != nil {
			return batchId, line, err
		}

		// import profile to MurmurationsServices Index
		postNodeUrl := config.Conf.Index.URL + "/v2/nodes"
		profileUrl := config.Conf.DataProxy.URL + "/v1/profiles/" + profileCuid
		nodeId, err := importutil.PostIndex(postNodeUrl, profileUrl)
		if err != nil {
			return batchId, line, errors.New("Import to MurmurationsServices Index failed: " + err.Error())
		}

		// save node_id to mongo
		profile["node_id"] = nodeId
		profile["is_posted"] = true
		err = s.batchRepo.SaveNodeId(profileCuid, profile)
		if err != nil {
			return batchId, line, errors.New("Save node_id to Mongo failed: " + err.Error())
		}
	}

	return batchId, -1, nil
}

func (s *batchUsecase) Edit(schemas []string, records [][]string, userId string, batchId string, metaName string, metaUrl string) (int, error) {
	if len(records) > 1001 {
		return -1, errors.New("CSV rows can't be larger than 1,000")
	}

	// check if batch_id belongs to user
	isValid, err := s.batchRepo.CheckUser(userId, batchId)
	if err != nil {
		return -1, err
	}
	if !isValid {
		return -1, errors.New("batch_id doesn't belong to user")
	}

	// get profile oid, cuid hash by batch_id
	profileOidsAndHashes, err := s.batchRepo.GetProfileOidsAndHashesByBatchId(batchId)
	if err != nil {
		return -1, err
	}

	rawProfiles := csvToMap(records)

	for line, rawProfile := range rawProfiles {
		profile, err := mapToProfile(rawProfile, schemas)
		if err != nil {
			return line, err
		}

		// hash profile
		profileHash, err := importutil.HashProfile(profile)
		if err != nil {
			return line, err
		}
		profile["source_data_hash"] = profileHash
		profile["batch_id"] = batchId

		// add metadata
		if metaName != "" && metaUrl != "" {
			metadata := map[string]interface{}{
				"sources": []map[string]interface{}{
					{
						"name": metaName,
						"url":  metaUrl,
					},
				},
			}
			profile["metadata"] = metadata
		}

		// check if profile exists in mongo
		_, ok := profileOidsAndHashes[profile["oid"].(string)]
		var profileCuid string
		if ok {
			profileCuid = profileOidsAndHashes[profile["oid"].(string)][0]
			// if current profile's oid and profile_hash match the data in mongo, skip it
			if profileOidsAndHashes[profile["oid"].(string)][1] == profileHash {
				continue
			}
			// update profile to Mongo
			err = s.batchRepo.UpdateProfile(profileCuid, profile)
			if err != nil {
				return line, err
			}
			// delete oid from profileOidsAndHashes, so that the rest of data in it needs to be deleted later
			delete(profileOidsAndHashes, profile["oid"].(string))
		} else {
			// if profile doesn't have cuid, generate one
			profileCuid = cuid.New()
			profile["cuid"] = profileCuid

			// import profile to Mongo
			err = s.batchRepo.SaveProfile(profile)
			if err != nil {
				return line, err
			}
		}

		// import profile to MurmurationsServices Index
		postNodeUrl := config.Conf.Index.URL + "/v2/nodes"
		profileUrl := config.Conf.DataProxy.URL + "/v1/profiles/" + profileCuid
		nodeId, err := importutil.PostIndex(postNodeUrl, profileUrl)
		if err != nil {
			return line, errors.New("Import to MurmurationsServices Index failed: " + err.Error())
		}

		// save node_id to mongo
		profile["node_id"] = nodeId
		profile["is_posted"] = true
		err = s.batchRepo.SaveNodeId(profileCuid, profile)
		if err != nil {
			return line, errors.New("Save node_id to Mongo failed: " + err.Error())
		}
	}

	// rest of data which are not in the csv file needs to be deleted
	if len(profileOidsAndHashes) > 0 {
		for _, cuidAndHash := range profileOidsAndHashes {
			// get profile by cuid
			profile, err := s.batchRepo.GetProfileByCuid(cuidAndHash[0])

			// delete profiles from mongo
			err = s.batchRepo.DeleteProfileByCuid(cuidAndHash[0])
			if err != nil {
				return -1, err
			}

			// delete profiles from MurmurationsServices Index
			if profile["is_posted"].(bool) {
				nodeId := profile["node_id"].(string)
				deleteNodeUrl := config.Conf.Index.URL + "/v2/nodes/" + nodeId
				err := importutil.DeleteIndex(deleteNodeUrl, nodeId)
				if err != nil {
					return -1, errors.New("Delete from MurmurationsServices Index failed: " + err.Error())
				}
			}
		}
	}

	return -1, nil
}

func (s *batchUsecase) Delete(userId string, batchId string) error {
	// check if batch_id belongs to user
	isValid, err := s.batchRepo.CheckUser(userId, batchId)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.New("batch_id doesn't belong to user")
	}

	// get profiles by batch_id
	profiles, err := s.batchRepo.GetProfilesByBatchId(batchId)
	if err != nil {
		return err
	}

	// delete profiles from mongo
	err = s.batchRepo.DeleteProfilesByBatchId(batchId)
	if err != nil {
		return err
	}

	// delete profiles from MurmurationsServices Index
	for _, profile := range profiles {
		if profile["is_posted"].(bool) {
			nodeId := profile["node_id"].(string)
			deleteNodeUrl := config.Conf.Index.URL + "/v2/nodes/" + nodeId
			err := importutil.DeleteIndex(deleteNodeUrl, nodeId)
			if err != nil {
				return errors.New("Delete profile from MurmurationsServices Index failed: " + err.Error())
			}
		}
	}

	// delete batch_id from mongo
	err = s.batchRepo.DeleteBatchId(batchId)
	if err != nil {
		return err
	}

	return nil
}

// convert csv to one-to-one map[string]string
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

// convert one-to-one map[string]string to profile data structure
func mapToProfile(rawProfile map[string]string, schemas []string) (map[string]interface{}, error) {
	profile := make(map[string]interface{})
	// handle geolocation
	if rawProfile["lat"] != "" && rawProfile["lon"] != "" {
		geolocation := make(map[string]float64)
		lat, err := getGeolocation(rawProfile["lat"])
		if err != nil {
			return nil, errors.New("parse location failed, err: " + err.Error())
		}
		lon, err := getGeolocation(rawProfile["lon"])
		if err != nil {
			return nil, errors.New("parse location failed, err: " + err.Error())
		}
		geolocation["lat"] = lat
		geolocation["lon"] = lon
		profile["geolocation"] = geolocation
	}
	delete(rawProfile, "lat")
	delete(rawProfile, "lon")

	// parse field name with hyphen, including array and array-object
	// array example: "tags-1" will be parsed to arrayFields["tags"] = value
	// array-object example: "urls-1-name" will be parsed to arrayObjectFields["urls"] = ["1-name"]
	arrayFields := make(map[string][]string)
	arrayObjectFields := make(map[string][]string)
	for key, value := range rawProfile {
		if strings.Contains(key, "-") {
			// one hyphen means it's an array field
			// delete the field from rawProfile, because we already save value to arrayFields
			hyphenIndex := strings.Index(key, "-")
			if strings.Count(key, "-") == 1 {
				arrayField := key[:hyphenIndex]
				arrayFields[arrayField] = append(arrayFields[arrayField], value)
				delete(rawProfile, key)
			} else if strings.Count(key, "-") == 2 {
				// handle array-object - two hyphens
				arrayField := key[:hyphenIndex]
				arrayFieldValue := key[hyphenIndex+1:]
				arrayObjectFields[arrayField] = append(arrayObjectFields[arrayField], arrayFieldValue)
			}
			// if there are more than three hyphens, it's not a valid field name, we will directly save it later.
		}
	}

	// handle array field
	for key, value := range arrayFields {
		profile[key] = value
	}

	// handle array-object field
	// it will sort arrayObjectFields by array number, combine them together if the number is the same
	for index, arrayObjectField := range arrayObjectFields {
		sort.Strings(arrayObjectField)
		currentNum := ""
		var objects []map[string]string
		object := make(map[string]string)
		for _, fieldName := range arrayObjectField {
			arrayNumIndex := strings.Index(fieldName, "-")
			arrayNum := fieldName[:arrayNumIndex]
			arrayVal := fieldName[arrayNumIndex+1:]
			if currentNum != "" && currentNum != arrayNum {
				objects = append(objects, object)
				object = make(map[string]string)
			}
			currentNum = arrayNum
			object[arrayVal] = rawProfile[index+"-"+fieldName]

			// after process data, remove from profile
			delete(rawProfile, index+"-"+fieldName)
		}
		objects = append(objects, object)
		profile[index] = objects
	}

	// handle rest of data
	for key, value := range rawProfile {
		profile[key] = value
	}

	// put schema here
	profile["linked_schemas"] = schemas
	return profile, nil
}

func getGeolocation(geolocation string) (float64, error) {
	precision := math.Pow(10, float64(8))
	float, err := strconv.ParseFloat(geolocation, 64)
	if err != nil {
		return 0, err
	}
	truncatedValue := math.Round(float*precision) / precision
	return truncatedValue, nil
}
