package usecase

import (
	"errors"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/repository/db"
	"math"
	"sort"
	"strconv"
	"strings"
)

type BatchUsecase interface {
	Validate([]string, [][]string) (int, error)
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

	// put csv rows into map[string]string
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

	// construct profile
	for line, rawProfile := range rawProfiles {
		profile := make(map[string]interface{})
		// handle geolocation
		if rawProfile["lat"] != "" && rawProfile["lon"] != "" {
			geolocation := make(map[string]float64)
			lat, err := getGeolocation(rawProfile["lat"])
			if err != nil {
				return line, errors.New("parse location failed, err: " + err.Error())
			}
			lon, err := getGeolocation(rawProfile["lon"])
			if err != nil {
				return line, errors.New("parse location failed, err: " + err.Error())
			}
			geolocation["lat"] = lat
			geolocation["lon"] = lon
			profile["geolocation"] = geolocation
		}
		delete(rawProfile, "lat")
		delete(rawProfile, "lon")

		// handle tags
		if rawProfile["tags"] != "" {
			profile["tags"] = strings.Split(rawProfile["tags"], ",")
			delete(rawProfile, "tags")
		}

		// handle array
		arrayFields := make(map[string][]string)
		for key, _ := range rawProfile {
			if strings.Contains(key, "-") {
				arrayFieldIndex := strings.Index(key, "-")
				arrayField := key[:arrayFieldIndex]
				arrayFieldValue := key[arrayFieldIndex+1:]
				arrayFields[arrayField] = append(arrayFields[arrayField], arrayFieldValue)
			}
		}

		for index, arrayField := range arrayFields {
			sort.Strings(arrayField)
			currentNum := ""
			var objects []map[string]string
			object := make(map[string]string)
			for _, fieldName := range arrayField {
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

func getGeolocation(geolocation string) (float64, error) {
	precision := math.Pow(10, float64(8))
	float, err := strconv.ParseFloat(geolocation, 64)
	if err != nil {
		return 0, err
	}
	truncatedValue := math.Round(float*precision) / precision
	return truncatedValue, nil
}
