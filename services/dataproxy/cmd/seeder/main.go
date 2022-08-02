package main

import (
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
	"github.com/lucsky/cuid"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func headerMapping(f *excelize.File) (map[string]string, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error while getting rows from excel, error message: %s", err)
	}

	headerMap := make(map[string]string)
	headerMap["id"] = "A"

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
			headerMap[colCell] = headerAlphabets[colIndex]
		}
	}
	return headerMap, nil
}

func importData(row int, schemaName string, headerMap map[string]string, mapping map[string]string, file *excelize.File) (bool, error) {
	oldProfile := make(map[string]interface{})
	for index, value := range headerMap {
		axis := value + strconv.Itoa(row)
		cell, err := file.GetCellValue(sheetName, axis)
		if err != nil {
			return false, fmt.Errorf("read Excel error, axis: %s, error message: %s", axis, err)
		}
		oldProfile[index] = cell
	}

	// deal with special types
	for k, v := range oldProfile {
		// Array type
		if k == "tags" {
			if v.(string) != "" {
				oldProfile[k] = strings.Split(v.(string), ",")
			}
			continue
		}
		if k == "categories" {
			oldProfileStr := strings.Split(v.(string), ",")
			oldProfileInterface := make([]interface{}, len(oldProfileStr))
			for i, v := range oldProfileStr {
				oldProfileInterface[i] = v
			}
			if len(oldProfileInterface) > 0 {
				oldProfile[k] = oldProfileInterface
			}
			continue
		}
		// Number type
		if k == "lat" || k == "lng" {
			num, err := strconv.ParseFloat(v.(string), 64)
			if err != nil {
				return false, fmt.Errorf("error when parsing number type data, error message: %s", err)
			}
			oldProfile[k] = num
			continue
		}
		// Default is String type
		oldProfile[k] = v
	}
	profileJson := importutil.MapProfile(oldProfile, mapping, schemaName)
	oid := profileJson["oid"].(string)

	// Validate data
	validateUrl := config.Conf.Index.URL + "/v2/validate"
	isValid, failureReasons, err := importutil.Validate(validateUrl, profileJson)
	if err != nil {
		return false, fmt.Errorf("error when trying to validate a profile, error message: %s", err)
	}
	if !isValid {
		return true, fmt.Errorf("warning: skip importing this row, validate profile failed, row: %v, id: %s, failure reasons: %s", row, oid, failureReasons)
	}

	// If database has same oid item, overwrite old data and show warning message
	filter := bson.M{"oid": oid}
	result, err := mongo.Client.Count(constant.MongoIndex.Profile, filter)
	if err != nil {
		return false, fmt.Errorf("error when trying to find a profile, error message: %s", err)
	}
	if result > 0 {
		return true, fmt.Errorf("warning: profile exist, the old data is not overwrited, row: %v, id: %s\n", row, oid)
	}

	// Save to MongoDB, return url to post index
	// Generate cid for item
	profileJson["cuid"] = cuid.New()
	_, err = mongo.Client.InsertOne(constant.MongoIndex.Profile, profileJson)
	if err != nil {
		return false, fmt.Errorf("error when trying to save a profile, error message: %s", err)
	}

	// Post to index service
	postNodeUrl := config.Conf.Index.URL + "/v2/nodes"
	profileUrl := config.Conf.DataProxy.URL + "/v1/profiles/" + profileJson["cuid"].(string)
	nodeId, err := importutil.PostIndex(postNodeUrl, profileUrl)
	if err != nil {
		return false, fmt.Errorf("failed to post profile to Index, profile url is %s, error message: %s", profileUrl, err)
	}

	// update NodeId
	update := bson.M{"$set": bson.M{"node_id": nodeId, "is_posted": false}}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	_, err = mongo.Client.FindOneAndUpdate(constant.MongoIndex.Profile, filter, update, opt)
	if err != nil {
		return true, err
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
	mapping, err := importutil.GetMapping(schemaName)
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
	headerMap, err := headerMapping(f)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Import data
	successNums := 0
	skippedNums := 0
	for i := from; i <= to; i++ {
		isSkipped, err := importData(i, schemaName, headerMap, mapping, f)
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
