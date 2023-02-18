package usecase

import (
	"errors"
	"fmt"
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
	Import([]string, [][]string, string) (string, int, error)
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

func (s *batchUsecase) Import(schemas []string, records [][]string, userCuid string) (string, int, error) {
	if len(records) > 1001 {
		return "", -1, errors.New("CSV rows can't be larger than 1,000")
	}

	// generate batch_id import using cuid and save it to mongo
	batchCuid := cuid.New()
	err := s.batchRepo.SaveUser(batchCuid, userCuid)
	if err != nil {
		return batchCuid, -1, err
	}

	rawProfiles := csvToMap(records)

	for line, rawProfile := range rawProfiles {
		profile, err := mapToProfile(rawProfile, schemas)
		if err != nil {
			return batchCuid, line, err
		}

		profile["batch_id"] = batchCuid

		// generate cuid for profile
		profileCuid := cuid.New()
		profile["cuid"] = profileCuid

		// import profile to Mongo
		err = s.batchRepo.SaveProfile(profile)
		if err != nil {
			return batchCuid, line, err
		}

		// import profile to MurmurationsServices Index
		postNodeUrl := config.Conf.Index.URL + "/v2/nodes"
		profileUrl := config.Conf.DataProxy.URL + "/v1/profiles/" + profileCuid
		fmt.Println(postNodeUrl)
		fmt.Println(profileUrl)
		nodeId, err := importutil.PostIndex(postNodeUrl, profileUrl)
		if err != nil {
			return batchCuid, line, errors.New("Import to MurmurationsServices Index failed: " + err.Error())
		}

		// save node_id to mongo
		profile["node_id"] = nodeId
		profile["is_posted"] = true
		err = s.batchRepo.SaveNodeId(profileCuid, profile)
		if err != nil {
			return batchCuid, line, errors.New("Save node_id to Mongo failed: " + err.Error())
		}
	}

	return batchCuid, -1, nil
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
