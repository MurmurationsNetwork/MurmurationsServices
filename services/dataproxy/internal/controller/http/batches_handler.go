package http

import (
	"encoding/csv"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/usecase"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type BatchesHandler interface {
	Validate(c *gin.Context)
	Import(c *gin.Context)
}

type batchesHandler struct {
	batchUsecase usecase.BatchUsecase
}

func NewBatchesHandler(batchService usecase.BatchUsecase) BatchesHandler {
	return &batchesHandler{
		batchUsecase: batchService,
	}
}

func (handler *batchesHandler) Validate(c *gin.Context) {
	file, schemas, errors := validateCsvInputs(c)
	if errors != nil {
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	records, errors := parseCsv(file)
	if errors != nil {
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	line, err := handler.batchUsecase.Validate(schemas, records)
	if err != nil {
		errors := jsonapi.NewError([]string{"CSV Validate Failed"}, []string{"Failed to validate line " + strconv.Itoa(line) + " Error: " + err.Error()}, nil, []int{http.StatusBadRequest})
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	meta := jsonapi.NewMeta("The submitted csv file was validated successfully to its linked schemas.", "", "")
	res := jsonapi.Response(nil, nil, nil, meta)
	c.JSON(http.StatusOK, res)
}

func (handler *batchesHandler) Import(c *gin.Context) {
	// todo: we might need to validate the excel before import
	// todo: we might need to validate user id from DB: but it costs more time to validate user id

	// cuid is 25 characters long
	userCuid := c.PostForm("cuid")
	if len(userCuid) != 25 {
		errors := jsonapi.NewError([]string{"Invalid User Cuid"}, []string{"User Cuid is not valid"}, nil, []int{http.StatusBadRequest})
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	file, schemas, errors := validateCsvInputs(c)
	if errors != nil {
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	records, errors := parseCsv(file)
	if errors != nil {
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	batchId, line, err := handler.batchUsecase.Import(schemas, records, userCuid)
	if err != nil {
		errors := jsonapi.NewError([]string{"CSV Import Failed"}, []string{"Failed to import line " + strconv.Itoa(line) + ", Batch Id: " + batchId + ", Error: " + err.Error()}, nil, []int{http.StatusBadRequest})
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	meta := jsonapi.NewBatchMeta("The submitted csv file was imported successfully to its linked schemas.", batchId)
	res := jsonapi.Response(nil, nil, nil, meta)
	c.JSON(http.StatusOK, res)
}

func validateCsvInputs(c *gin.Context) (*multipart.FileHeader, []string, []jsonapi.Error) {
	// get fields from the POST request
	file, err := c.FormFile("file")
	if err != nil {
		errors := jsonapi.NewError([]string{"Get File Error"}, []string{"The File document submitted could not be parsed."}, [][]string{{"parameter", "file"}}, []int{http.StatusBadRequest})
		return nil, nil, errors
	}

	// if the file is not csv, we cannot process it
	fileName := file.Filename
	if fileName[len(fileName)-4:] != ".csv" {
		errors := jsonapi.NewError([]string{"Get File Error"}, []string{"The File document submitted is not csv."}, [][]string{{"parameter", "file"}}, []int{http.StatusBadRequest})
		return nil, nil, errors
	}

	// get schemas from the POST request
	rawSchemas := c.PostForm("schemas")
	if rawSchemas == "" {
		errors := jsonapi.NewError([]string{"Invalid Query Parameter"}, []string{"The following query parameter is not valid: schemas."}, [][]string{{"parameter", "schemas"}}, []int{http.StatusBadRequest})
		return nil, nil, errors
	}

	// make schemas to []string
	rawSchemas = strings.ReplaceAll(rawSchemas, "\"", "")
	schemas := strings.Split(rawSchemas[1:len(rawSchemas)-1], ",")

	return file, schemas, nil
}

func parseCsv(file *multipart.FileHeader) ([][]string, []jsonapi.Error) {
	// parse csv and put all data in service
	rawFile, err := file.Open()
	if err != nil {
		errors := jsonapi.NewError([]string{"File Open Error"}, []string{"The file is corrupted and can't be opened."}, [][]string{{"parameter", "file"}}, []int{http.StatusBadRequest})
		return nil, errors
	}

	csvReader := csv.NewReader(rawFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		errors := jsonapi.NewError([]string{"File Open Error"}, []string{"Unable to parse file as csv."}, [][]string{{"parameter", "file"}}, []int{http.StatusBadRequest})
		return nil, errors
	}

	return records, nil
}
