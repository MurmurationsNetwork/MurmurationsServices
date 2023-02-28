package usecase

import (
	"errors"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/repository/db"
	"github.com/lucsky/cuid"
	"sort"
	"strconv"
	"strings"
)

type BatchUsecase interface {
	GetBatchesByUserID(string) ([]string, error)
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

func (s *batchUsecase) GetBatchesByUserID(userId string) ([]string, error) {
	batches, err := s.batchRepo.GetBatchesByUserID(userId)
	if err != nil {
		return nil, err
	}

	return batches, nil
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

		// validate data, once it has error response error
		validateUrl := config.Conf.Index.URL + "/v2/validate"
		isValid, failureReasons, err := importutil.Validate(validateUrl, profile)
		if err != nil {
			return batchId, line, err
		}
		if !isValid {
			return batchId, line, errors.New(failureReasons)
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
		oid := profile["oid"].(string)
		_, ok := profileOidsAndHashes[oid]
		var profileCuid string
		if ok {
			profileCuid = profileOidsAndHashes[oid][0]
			// if current profile's oid and profile_hash match the data in mongo, skip it
			if profileOidsAndHashes[oid][1] == profileHash {
				delete(profileOidsAndHashes, oid)
				continue
			}
			// update profile to Mongo
			err = s.batchRepo.UpdateProfile(profileCuid, profile)
			if err != nil {
				return line, err
			}
			// delete oid from profileOidsAndHashes, so that the rest of data in it needs to be deleted later
			delete(profileOidsAndHashes, oid)
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
	// sort rawProfile by key
	keys := make([]string, 0, len(rawProfile))
	// validate oid field
	hasOid := false
	for k := range rawProfile {
		if k == "oid" {
			hasOid = true
		}
		keys = append(keys, k)
	}

	if !hasOid {
		return nil, errors.New("oid field is required")
	}

	// sort the keys
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

	// put schema here
	profile["linked_schemas"] = schemas
	return profile, nil
}

// destruct field name and save field value to profile data structure
func destructField(profile map[string]interface{}, field string, value string) (map[string]interface{}, error) {
	// destruct field name
	// e.g. "urls[0].name" -> ["urls", 0, "name"], "tags[0]" -> ["tags", 0]
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
	current := profile
	for i, p := range path {
		// if the current path is a number, skip it, because it's already handled in the previous loop
		if _, err := strconv.Atoi(p); err == nil {
			continue
		}
		// if the next path is a number, and it's the last element, it means it's an array
		if i == len(path)-2 {
			if _, err := strconv.Atoi(path[i+1]); err == nil {
				if _, ok := current[path[i]]; !ok {
					current[path[i]] = make([]interface{}, 0)
				}
				current[path[i]] = append(current[path[i]].([]interface{}), destructValue(value))
				break
			}
		}
		// if the next path is a number, it means it's an array-object
		if i+1 < len(path) {
			if arrayNum, err := strconv.Atoi(path[i+1]); err == nil {
				if _, ok := current[path[i]]; !ok {
					current[path[i]] = make([]map[string]interface{}, 0)
				}
				if _, ok := current[path[i]].([]map[string]interface{}); !ok {
					return nil, errors.New("Check if the fields are duplicated or have different types of fields with the same name. Invalid field name: " + field)
				}
				if len(current[path[i]].([]map[string]interface{})) <= arrayNum {
					current[path[i]] = append(current[path[i]].([]map[string]interface{}), make(map[string]interface{}))
				}
				if len(current[path[i]].([]map[string]interface{}))-1 != arrayNum {
					return nil, errors.New("Check the filed name's array number is sequential and start from 0. Invalid field name: " + field)
				}
				current = current[path[i]].([]map[string]interface{})[arrayNum]
				continue
			}
		}
		// if the last element, put value into it.
		if i == len(path)-1 {
			current[p] = destructValue(value)
			break
		}
		if _, ok := current[p]; !ok {
			current[p] = make(map[string]interface{})
		}
		if _, ok := current[p].(map[string]interface{}); !ok {
			return nil, errors.New("Check if the fields are duplicated or have different types of fields with the same name. Invalid field name: " + field)
		}
		current = current[p].(map[string]interface{})
	}

	return profile, nil
}

func destructValue(value string) interface{} {
	// if the string has dot, it means it's possible a float number
	// if not, it's possible an int number
	// in other cases, it's a string
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
