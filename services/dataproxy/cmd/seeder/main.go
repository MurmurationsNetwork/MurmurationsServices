package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
	"github.com/lucsky/cuid"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var headerAlphabets = map[int]string{
	0:  "A",
	1:  "B",
	2:  "C",
	3:  "D",
	4:  "E",
	5:  "F",
	6:  "G",
	7:  "H",
	8:  "I",
	9:  "J",
	10: "K",
	11: "L",
	12: "M",
	13: "N",
	14: "O",
	15: "P",
	16: "Q",
	17: "R",
	18: "S",
	19: "T",
	20: "U",
	21: "V",
	22: "W",
	23: "X",
	24: "Y",
	25: "Z",
}

// global variables
var fileName = "schema.xlsx"
var sheetName = "sheet1"

func Init() {
	config.Init()
	mongoInit()
}

func mongoInit() {
	uri := mongo.GetURI(config.Conf.Mongo.USERNAME, config.Conf.Mongo.PASSWORD, config.Conf.Mongo.HOST)

	err := mongo.NewClient(uri, config.Conf.Mongo.DBName)
	if err != nil {
		fmt.Println("error when trying to connect to MongoDB", err)
		os.Exit(1)
	}
	err = mongo.Client.Ping()
	if err != nil {
		fmt.Println("error when trying to ping the MongoDB", err)
		os.Exit(1)
	}
}

func readArgs(args []string) (string, string, int, int, error) {
	/*
		There are four arguments.
		1. EXCEL_URL 2. SCHEMA_NAME 3. FROM (row) 4. TO (row)
	*/
	if len(args) != 5 {
		return "", "", 0, 0, fmt.Errorf("missing arguments, please check the arguments")
	}
	from, err := strconv.Atoi(args[3])
	to, err := strconv.Atoi(os.Args[4])
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("from or to argument must be the number")
	}
	return args[1], args[2], from, to, nil
}

func downloadExcel(url string) error {
	fmt.Println("Downloading excel from remote server...")
	output, err := os.Create(fileName)
	defer output.Close()
	if err != nil {
		return fmt.Errorf("error while create file %s , error message: %s", fileName, err)
	}

	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("error while downloading from %s , error message: %s", url, err)
	}

	n, err := io.Copy(output, res.Body)
	if err != nil {
		return fmt.Errorf("error while receiving file %s data, error message: %s", fileName, err)
	}
	fmt.Println("Retrieve Excel successful:", n, "bytes downloaded.")
	return nil
}

func getMapping(schemaName string) (map[string]interface{}, error) {
	filter := bson.M{"schema": schemaName}
	result := mongo.Client.FindOne(constant.MongoIndex.Mapping, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("could not find mapping for schema %s", schemaName)
		}
		return nil, fmt.Errorf("error when trying to find the mapping, error message: %s", result.Err())
	}
	schemaRaw := make(map[string]interface{})
	err := result.Decode(schemaRaw)
	if err != nil {
		return nil, fmt.Errorf("error when trying to parse database response, error message: %s", result.Err())
	}

	// remove id and __v
	schema := make(map[string]interface{})
	for i, v := range schemaRaw {
		if i == "__v" || i == "_id" {
			continue
		}
		schema[i] = v
	}
	return schema, nil
}

func headerMapping(schema map[string]interface{}, f *excelize.File) (map[string]string, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error while getting rows from excel, error message: %s", err)
	}

	headerMap := make(map[string]string)
	headerMap["oid"] = "A"

	for i := 0; i < 1; i++ {

	}
	for rowIndex, row := range rows {
		// only need the header
		if rowIndex != 0 {
			break
		}
		for colIndex, colCell := range row {
			// column exceeds the limit
			if colIndex > 25 {
				return nil, fmt.Errorf("excel header can't have more than 25 columns, contact administrator to expend the size")
			}
			for fieldName, header := range schema {
				if colCell == header {
					headerMap[fieldName] = headerAlphabets[colIndex]
				}
			}
		}
	}
	return headerMap, nil
}

func validate(schema string, profile map[string]interface{}) (bool, string, error) {
	for k, v := range profile {
		if v == "" {
			delete(profile, k)
			continue
		}
		// Array type
		if k == "tags" {
			profile[k] = strings.Split(v.(string), ",")
			continue
		}
		// Number type
		if k == "latitude" || k == "longitude" {
			num, err := strconv.ParseFloat(v.(string), 64)
			if err != nil {
				return false, "", fmt.Errorf("error when parsing number type data, error message: %s", err)
			}
			profile[k] = num
			continue
		}
		// Default is String type
		profile[k] = v
	}

	// Add linked_schemas
	var s []string
	s = append(s, schema)
	profile["linked_schemas"] = s
	profileJson, err := json.Marshal(profile)
	if err != nil {
		return false, "", err
	}

	// Validate from index service
	validateUrl := config.Conf.Index.URL + "/v2/validate"
	res, err := http.Post(validateUrl, "application/json", bytes.NewBuffer(profileJson))
	if err != nil {
		return false, "", err
	}
	if res.StatusCode != 200 {
		return false, "", fmt.Errorf("validate failed, the status code is %s. json data: %s", strconv.Itoa(res.StatusCode), string(profileJson))
	}

	var resBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&resBody)
	statusCode := int64(resBody["status"].(float64))
	if statusCode != 200 {
		if resBody["failure_reasons"] != nil {
			var failureReasons []string
			for _, item := range resBody["failure_reasons"].([]interface{}) {
				failureReasons = append(failureReasons, item.(string))
			}
			failureReasonsStr := strings.Join(failureReasons, ",")
			return false, failureReasonsStr, nil
		}
		return false, "failed without reasons!", nil
	}
	return true, "", nil
}

func importData(row int, schemaName string, headerMap map[string]string, file *excelize.File) (bool, error) {
	profileJson := make(map[string]interface{})
	for index, value := range headerMap {
		axis := value + strconv.Itoa(row)
		cell, err := file.GetCellValue(sheetName, axis)
		if err != nil {
			return false, fmt.Errorf("read Excel error, axis: %s, error message: %s", axis, err)
		}
		profileJson[index] = cell
	}

	// Validate data
	isValid, failureReasons, err := validate(schemaName, profileJson)
	if err != nil {
		return false, fmt.Errorf("error when trying to validate a profile, error message: %s", err)
	}
	if !isValid {
		return true, fmt.Errorf("warning: skip importing this row, validate profile failed, row: %v, id: %s, failure reasons: %s", row, profileJson["oid"], failureReasons)
	}

	// If database has same oid item, skip it and show warning message
	filter := bson.M{"oid": profileJson["oid"]}
	result, err := mongo.Client.Count(constant.MongoIndex.Profile, filter)
	if err != nil {
		return false, fmt.Errorf("error when trying to find a profile, error message: %s", err)
	}
	if result > 0 {
		return true, fmt.Errorf("warning: skip importing this row, profile exist, row: %v, id: %s", row, profileJson["oid"])
	}

	// Generate cid for item
	profileJson["cuid"] = cuid.New()

	// Save to MongoDB, return url to post index
	_, err = mongo.Client.InsertOne(constant.MongoIndex.Profile, profileJson)
	if err != nil {
		return false, fmt.Errorf("error when trying to save a profile, error message: %s", err)
	}

	// Post to index service
	postNodeUrl := config.Conf.Index.URL + "/v2/nodes"
	postProfile := make(map[string]string)
	postProfile["profile_url"] = config.Conf.DataProxy.URL + "/v1/profiles/" + profileJson["cuid"].(string)
	postProfileJson, err := json.Marshal(postProfile)
	if err != nil {
		return false, fmt.Errorf("error when trying to marshal a profile, url: %s, error message: %s", postProfile["profile_url"], err)
	}
	res, err := http.Post(postNodeUrl, "application/json", bytes.NewBuffer(postProfileJson))
	if err != nil {
		return false, fmt.Errorf("error when trying to post a profile, error message: %s", err)
	}
	if res.StatusCode != 200 {
		return false, fmt.Errorf("post failed, the status code is %s. url: %s", strconv.Itoa(res.StatusCode), postProfile["profile_url"])
	}
	return false, nil
}

func cleanUp() error {
	// turn off connection with MongoDB
	mongo.Client.Disconnect()

	// delete the local file
	err := os.Remove(fileName)
	if err != nil {
		return fmt.Errorf("error when deleting the file %s, error message: %s", fileName, err)
	}
	return nil
}

func main() {
	fmt.Println("Hi, Welcome to Murmurations Seeder. 🎉")

	// Init config and mongoDB connection
	Init()

	// Get the arguments
	url, schemaName, from, to, err := readArgs(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Download Excel
	err = downloadExcel(url)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Find the mapping schema from MongoDB
	schema, err := getMapping(schemaName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Open Excel file
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Printf("error while reading excel, error message: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Mapping excel header with schema mapping
	headerMap, err := headerMapping(schema, f)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Import data
	successNums := 0
	skippedNums := 0
	for i := from; i <= to; i++ {
		isSkipped, err := importData(i, schemaName, headerMap, f)
		if err != nil {
			fmt.Println(err.Error())
			if isSkipped {
				skippedNums++
				continue
			}
			os.Exit(1)
		}
		successNums++
	}
	totalNums := to - from + 1
	failedNums := totalNums - successNums - skippedNums

	fmt.Printf("successfully imported profiles, total profiles: %v, success: %v ,skipped: %v, failed: %v\n", totalNums, successNums, skippedNums, failedNums)

	// Disconnect MongoDB and delete excel file
	err = cleanUp()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
