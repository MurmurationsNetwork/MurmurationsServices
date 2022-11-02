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
)

// global variables
var fileName = "schema-uk-data.xlsx"
var sheetName = "sheet1"

func Init() {
	config.Init()
	mongoInit()
}

func mongoInit() {
	uri := mongo.GetURI(config.Conf.Mongo.USERNAME, config.Conf.Mongo.PASSWORD, config.Conf.Mongo.HOST)

	err := mongo.NewClient(uri, config.Conf.Mongo.DBName)
	if err != nil {
		fmt.Println("Error when trying to connect to MongoDB.", err)
		os.Exit(1)
	}
	err = mongo.Client.Ping()
	if err != nil {
		fmt.Println("Error when trying to ping MongoDB.", err)
		os.Exit(1)
	}
}

func readArgs(args []string) (string, int, int, error) {
	/*
		There are four arguments.
		1. EXCEL_URL 2. FROM (row) 3. TO (row)
	*/
	if len(args) != 4 {
		return "", 0, 0, fmt.Errorf("Missing arguments: please check the arguments.")
	}
	from, err := strconv.Atoi(args[2])
	to, err := strconv.Atoi(os.Args[3])
	if err != nil {
		return "", 0, 0, fmt.Errorf("From or to argument must be an integer.")
	}
	return args[1], from, to, nil
}

func downloadExcel(url string) error {
	fmt.Println("Downloading Excel file from remote server...")
	output, err := os.Create(fileName)
	defer output.Close()
	if err != nil {
		return fmt.Errorf("Error while creating %s file: %s", fileName, err)
	}

	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("Error while downloading from %s: %s", url, err)
	}

	n, err := io.Copy(output, res.Body)
	if err != nil {
		return fmt.Errorf("Error while receiving data from %s: %s", fileName, err)
	}
	fmt.Println("Excel file retrieved successfully: ", n, "bytes downloaded.")
	return nil
}

func importData(row int, file *excelize.File) (bool, error) {
	// get excel oid
	axis := "C" + strconv.Itoa(row)
	name, err := file.GetCellValue(sheetName, axis)
	if err != nil {
		return false, fmt.Errorf("Error reading Excel file. Axis: %s, error message: %s", axis, err)
	}

	// If database has same oid item, keep the old data and show warning message
	filter := bson.M{"name": name}
	result, err := mongo.Client.Count(constant.MongoIndex.Profile, filter)
	if err != nil {
		return false, fmt.Errorf("Error when trying to find a profile: %s", err)
	}
	if result > 0 {
		return true, fmt.Errorf("Warning: profile already exists. The old data has not been overwritten. Row: %v, Name: %s\n", row, name)
	}

	// Todo: This seeder will be abandoned in near future, so I create data in simple way.
	locality, err := file.GetCellValue(sheetName, "D"+strconv.Itoa(row))
	region, err := file.GetCellValue(sheetName, "E"+strconv.Itoa(row))
	countryName, err := file.GetCellValue(sheetName, "G"+strconv.Itoa(row))
	description, err := file.GetCellValue(sheetName, "I"+strconv.Itoa(row))
	primaryUrl, err := file.GetCellValue(sheetName, "J"+strconv.Itoa(row))
	lat, err := file.GetCellValue(sheetName, "K"+strconv.Itoa(row))
	lon, err := file.GetCellValue(sheetName, "L"+strconv.Itoa(row))

	if err != nil {
		return false, fmt.Errorf("Error when getting excel cell value: %s", err)
	}

	profile := make(map[string]interface{})
	profile["linked_schemas"] = []string{"organizations_schema-v1.0.0"}
	profile["name"] = name
	profile["description"] = description
	profile["primary_url"] = primaryUrl
	profile["locality"] = locality
	profile["region"] = region
	profile["country_name"] = countryName
	geolocation := make(map[string]interface{})
	geolocation["lat"], err = strconv.ParseFloat(lat, 64)
	geolocation["lon"], err = strconv.ParseFloat(lon, 64)
	if err != nil {
		return false, fmt.Errorf("Error when parsing float number: %s", err)
	}
	profile["geolocation"] = geolocation
	profile["tags"] = []string{"Co-op"}

	metadata := map[string]interface{}{
		"sources": []map[string]interface{}{
			{
				"name":             "Co-operatives UK",
				"url":              "https://www.uk.coop",
				"profile_data_url": "https://www.uk.coop/resources/open-data",
				"access_time":      1630450800,
			},
		},
	}
	profile["metadata"] = metadata

	// Validate data
	validateUrl := config.Conf.Index.URL + "/v2/validate"
	isValid, failureReasons, err := importutil.Validate(validateUrl, profile)
	if err != nil {
		return false, fmt.Errorf("Error when trying to validate a profile: %s", err)
	}
	if !isValid {
		return true, fmt.Errorf("Warning: skipped importing this row because profile validation failed. Row: %v, Name: %s, Failure Reasons: %s", row, name, failureReasons)
	}

	// Save to MongoDB, return url to post index
	// Generate cid for item
	profile["cuid"] = cuid.New()
	_, err = mongo.Client.InsertOne(constant.MongoIndex.Profile, profile)
	if err != nil {
		return false, fmt.Errorf("Error when trying to save a profile: %s", err)
	}

	// Post to index service
	postNodeUrl := config.Conf.Index.URL + "/v2/nodes"
	profileUrl := config.Conf.DataProxy.URL + "/v1/profiles/" + profile["cuid"].(string)
	nodeId, err := importutil.PostIndex(postNodeUrl, profileUrl)
	if err != nil {
		return false, fmt.Errorf("failed to post %s to Index: %s", profileUrl, err)
	}

	// update NodeId
	update := bson.M{"$set": bson.M{"node_id": nodeId, "is_posted": true}}
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
		return fmt.Errorf("Error when deleting the file %s - %s", fileName, err)
	}
	return nil
}

func main() {
	fmt.Println("Hi, Welcome to Murmurations Seeder. 🎉🚀")

	// Init config and mongoDB connection
	Init()

	// Get the arguments
	url, from, to, err := readArgs(os.Args)
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

	// Open Excel file
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Printf("Error while reading Excel file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Import data
	successNums := 0
	skippedNums := 0
	for i := from; i <= to; i++ {
		isSkipped, err := importData(i, f)
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

	fmt.Printf("Successfully imported profiles. Total profiles: %v, Success: %v , Skipped: %v, Failed: %v\n", totalNums, successNums, skippedNums, failedNums)

	// Disconnect MongoDB and delete excel file
	err = cleanUp()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
