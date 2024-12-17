package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/lucsky/cuid"
	excelize "github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/importutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
)

const (
	fileName  = "schema.xlsx"
	sheetName = "sheet1"
)

func Init() {
	config.Init()
	mongoInit()
}

func mongoInit() {
	uri := mongo.GetURI(
		config.Values.Mongo.USERNAME,
		config.Values.Mongo.PASSWORD,
		config.Values.Mongo.HOST,
	)

	err := mongo.NewClient(uri, config.Values.Mongo.DBName)
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

func readArgs(args []string) (string, string, int, int, error) {
	/*
		There are four arguments.
		1. EXCEL_URL 2. SCHEMA_NAME 3. FROM (row) 4. TO (row)
	*/
	if len(args) != 5 {
		return "", "", 0, 0, fmt.Errorf(
			"missing arguments: please check the arguments",
		)
	}

	from, err := strconv.Atoi(args[3])
	if err != nil {
		return "", "", 0, 0, fmt.Errorf(
			"from argument must be an integer: %v",
			err,
		)
	}

	to, err := strconv.Atoi(os.Args[4])
	if err != nil {
		return "", "", 0, 0, fmt.Errorf(
			"to argument must be an integer: %v",
			err,
		)
	}

	return args[1], args[2], from, to, nil
}

func downloadExcel(url string) error {
	fmt.Println("Downloading Excel file from remote server...")
	output, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error while creating %s file: %s", fileName, err)
	}
	defer output.Close()

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error while downloading from %s: %s", url, err)
	}
	defer res.Body.Close()

	n, err := io.Copy(output, res.Body)
	if err != nil {
		return fmt.Errorf(
			"error while receiving data from %s: %s",
			fileName,
			err,
		)
	}
	fmt.Println("Excel file retrieved successfully: ", n, "bytes downloaded.")
	return nil
}

func importData(
	row int,
	schemaName string,
	mapping map[string]string,
	file *excelize.File,
) (bool, error) {
	// get excel oid
	axis := "A" + strconv.Itoa(row)
	oid, err := file.GetCellValue(sheetName, axis)
	if err != nil {
		return false, fmt.Errorf(
			"error reading Excel file. Axis: %s, error message: %s",
			axis,
			err,
		)
	}

	// If database has same oid item, keep the old data and show warning message
	filter := bson.M{"oid": oid}
	result, err := mongo.Client.Count(constant.MongoIndex.Profile, filter)
	if err != nil {
		return false, fmt.Errorf("error when trying to find a profile: %s", err)
	}
	if result > 0 {
		return true, fmt.Errorf(
			"warning: profile already exists. The old data has not been overwritten. Row: %v, OID: %s",
			row,
			oid,
		)
	}

	url := "https://api.ofdb.io/v0/entries/" + oid
	res, err := httputil.Get(url)
	if err != nil {
		return false, fmt.Errorf("can't get data from " + url)
	}
	defer res.Body.Close()

	var oldProfiles []map[string]interface{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&oldProfiles)
	if err != nil {
		return false, fmt.Errorf("can't parse data from " + url)
	}

	if len(oldProfiles) == 0 {
		return true, fmt.Errorf("profile doesn't exist. OID: " + oid)
	}

	profileJSON, err := importutil.MapProfile(
		oldProfiles[0],
		mapping,
		schemaName,
	)
	if err != nil {
		return false, fmt.Errorf("error when trying to map a profile: %s", err)
	}

	// Validate data
	validateURL := config.Values.Index.URL + "/v2/validate"
	isValid, failureReasons, err := importutil.Validate(
		validateURL,
		profileJSON,
	)
	if err != nil {
		return false, fmt.Errorf(
			"error when trying to validate a profile: %s",
			err,
		)
	}
	if !isValid {
		return true, fmt.Errorf(
			"warning: skipped importing this row because profile validation failed. Row: %v, OID: %s, Failure Reasons: %s",
			row,
			oid,
			failureReasons,
		)
	}

	// Save to MongoDB, return url to post index
	// Generate cid for item
	profileJSON["cuid"] = cuid.New()
	_, err = mongo.Client.InsertOne(constant.MongoIndex.Profile, profileJSON)
	if err != nil {
		return false, fmt.Errorf("error when trying to save a profile: %s", err)
	}

	// Post to index service
	postNodeURL := config.Values.Index.URL + "/v2/nodes"
	profileURL := config.Values.DataProxy.URL + "/v1/profiles/" + profileJSON["cuid"].(string)
	nodeID, err := importutil.PostIndex(postNodeURL, profileURL)
	if err != nil {
		return false, fmt.Errorf(
			"failed to post %s to Index: %s",
			profileURL,
			err,
		)
	}

	// update NodeID
	update := bson.M{"$set": bson.M{"node_id": nodeID, "is_posted": true}}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	_, err = mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Profile,
		filter,
		update,
		opt,
	)
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
		return fmt.Errorf("error when deleting the file %s - %s", fileName, err)
	}
	return nil
}

func main() {
	fmt.Println("Hi, Welcome to Murmurations Seeder. 🎉🚀")

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
		fmt.Printf("Error while reading Excel file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Import data
	successNums := 0
	skippedNums := 0
	for i := from; i <= to; i++ {
		isSkipped, err := importData(i, schemaName, mapping, f)
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

	fmt.Printf(
		"Successfully imported profiles. Total profiles: %v, Success: %v , Skipped: %v, Failed: %v\n",
		totalNums,
		successNums,
		skippedNums,
		failedNums,
	)

	// Disconnect MongoDB and delete excel file
	err = cleanUp()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
