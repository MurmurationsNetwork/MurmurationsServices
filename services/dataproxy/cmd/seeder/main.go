package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/global"
	"github.com/lucsky/cuid"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var alphabet = map[int]string{
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

func main() {
	fmt.Println("Hi, Welcome to Murmurations Seeder. 🎉")

	global.Init()
	mongoInit()

	// 1st argument: import excel url
	// 2nd argument: schema name
	// 3rd argument: from data id
	// 4th argument: to data id
	if len(os.Args) != 5 {
		fmt.Println("Missing argument. Please check the argument.")
		os.Exit(1)
	}

	url := os.Args[1]
	schemaName := os.Args[2]
	fromString := os.Args[3]
	toString := os.Args[4]
	from, err := strconv.Atoi(fromString)
	to, err := strconv.Atoi(toString)
	if err != nil {
		fmt.Println("from or to argument must be the number.")
		os.Exit(1)
	}

	// download excel from server
	fmt.Println("Downloading excel from remote server...")
	fileName := "schema.xlsx"
	output, err := os.Create(fileName)
	defer output.Close()

	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("Error while downloading from ", url, ", error message: ", err)
		os.Exit(1)
	}

	n, err := io.Copy(output, res.Body)
	fmt.Println("Retrieve Excel successful:", n, "bytes downloaded.")

	// check mapping db has the name with schema_name
	filter := bson.M{"schema": schemaName}
	result := mongo.Client.FindOne(constant.MongoIndex.Mapping, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			fmt.Println("Could not find mapping: ", schemaName)
			os.Exit(1)
		}
		fmt.Println("Error when trying to find mapping: ", result.Err())
		os.Exit(1)
	}

	schemaRaw := make(map[string]interface{})
	err = result.Decode(schemaRaw)
	if err != nil {
		fmt.Println("Error when trying to parse database response", result.Err())
		os.Exit(1)
	}

	// remove id and __v
	schema := make(map[string]interface{})
	for i, v := range schemaRaw {
		if i == "__v" || i == "_id" {
			continue
		}
		schema[i] = v
	}

	// Start processing data according to "from" and "to"
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println("Error while reading excel: ", err)
		os.Exit(1)
	}
	defer f.Close()

	rows, err := f.GetRows("sheet1")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	excelMap := make(map[string]string)
	excelMap["oid"] = "A"

	for index, row := range rows {
		if index != 0 {
			break
		}
		for colIndex, colCell := range row {
			if colIndex > 25 {
				fmt.Println("Excel header can't have more than 25 columns. Contact Administrator to expend the size.")
				os.Exit(1)
			}
			for i, v := range schema {
				if colCell == v {
					excelMap[i] = alphabet[colIndex]
				}
			}
		}
	}

	for i := from; i <= to; i++ {
		json := make(map[string]interface{})
		for index, value := range excelMap {
			axis := value + strconv.Itoa(i)
			cell, err := f.GetCellValue("sheet1", axis)
			if err != nil {
				fmt.Println("read Excel error, axis=", axis, " error message: ", err)
				os.Exit(1)
			}
			json[index] = cell
		}

		// validate data
		isValid, failureReasons, err := validate(schemaName, json)
		if err != nil {
			fmt.Println("Error when trying to validate a profile", err)
			os.Exit(1)
		}
		if !isValid {
			fmt.Println("Error: skip importing this row, validate profile failed, row=", i, " id=", json["oid"], "failure reasons =", failureReasons)
			continue
		}

		// if database has same oid item, skip it and show warning message
		filter := bson.M{"oid": json["oid"]}
		result, err := mongo.Client.Count(constant.MongoIndex.Profile, filter)
		if err != nil {
			fmt.Println("Error when trying to find a profile", err)
			os.Exit(1)
		}
		if result > 0 {
			fmt.Println("Warning: skip importing this row, profile exist, row=", i, " id=", json["oid"])
			continue
		}

		// generate cid for item
		json["cuid"] = cuid.New()

		// save to MongoDB, return url to post index
		profile, err := mongo.Client.InsertOne(constant.MongoIndex.Profile, json)
		if err != nil {
			fmt.Println("Error when trying to save a profile", err)
			os.Exit(1)
		}
		fmt.Println(profile)
		fmt.Println(json["cuid"])
	}

	// turn off connection with MongoDB

	// delete the local file
	err = os.Remove(fileName)
	if err != nil {
		fmt.Println("Error when deleting file: ", err)
	}
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

func validate(schema string, profile map[string]interface{}) (bool, string, error) {
	for k, v := range profile {
		if v == "" {
			delete(profile, k)
			continue
		}
		num, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			profile[k] = v
		} else {
			profile[k] = num
		}
	}
	// add schema
	s := make([]string, 1)
	s[0] = schema
	profile["linked_schemas"] = s
	profileJson, err := json.Marshal(profile)
	if err != nil {
		return false, "", err
	}
	// validate from index service
	validateUrl := config.Conf.Index.URL + "/v2/validate"
	res, err := http.Post(validateUrl, "application/json", bytes.NewBuffer(profileJson))
	if err != nil {
		return false, "", err
	}
	if res.StatusCode != 200 {
		err = fmt.Errorf("The status code is " + strconv.Itoa(res.StatusCode))
		return false, "", err
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
		return false, "Failed without reasons", nil
	}
	return true, "", nil
}
