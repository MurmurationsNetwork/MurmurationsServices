package http

import (
	"encoding/json"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/xeipuuv/gojsonschema"
	"net/http"
	"strings"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity/query"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/usecase"
	"github.com/gin-gonic/gin"
)

type NodeHandler interface {
	Add(c *gin.Context)
	Get(c *gin.Context)
	Search(c *gin.Context)
	Delete(c *gin.Context)
	AddSync(c *gin.Context)
	Validate(c *gin.Context)
}

type nodeHandler struct {
	nodeUsecase usecase.NodeUsecase
}

func NewNodeHandler(nodeService usecase.NodeUsecase) NodeHandler {
	return &nodeHandler{
		nodeUsecase: nodeService,
	}
}

func (handler *nodeHandler) getNodeId(params gin.Params) (string, resterr.RestErr) {
	nodeId, found := params.Get("nodeId")
	if !found {
		return "", resterr.NewBadRequestError("Invalid node_id.")
	}
	return nodeId, nil
}

func (handler *nodeHandler) Add(c *gin.Context) {
	var node nodeDTO
	if err := c.ShouldBindJSON(&node); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body.")
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	if err := node.Validate(); err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	result, err := handler.nodeUsecase.AddNode(node.toEntity())
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, handler.toAddNodeVO(result))
}

func (handler *nodeHandler) Get(c *gin.Context) {
	nodeId, err := handler.getNodeId(c.Params)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	node, err := handler.nodeUsecase.GetNode(nodeId)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, handler.toGetNodeVO(node))
}

func (handler *nodeHandler) Search(c *gin.Context) {
	var query query.EsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body.")
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	if query.Page*query.PageSize > 10000 {
		restErr := resterr.NewBadRequestError("No more than 10,000 results can be returned. Refine your query so it will return less but more relevant results.")
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	searchResult, err := handler.nodeUsecase.Search(&query)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, searchResult.ToVO())
}

func (handler *nodeHandler) Delete(c *gin.Context) {
	if c.Params.ByName("nodeId") == "" {
		restErr := resterr.NewBadRequestError("The node_id path parameter is missing.")
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	nodeId, err := handler.getNodeId(c.Params)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	profileUrl, err := handler.nodeUsecase.Delete(nodeId)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Removed profile at " + profileUrl,
		"status":  http.StatusOK,
	})
}

func (handler *nodeHandler) AddSync(c *gin.Context) {
	var node nodeDTO
	if err := c.ShouldBindJSON(&node); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body.")
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	if err := node.Validate(); err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	result, err := handler.nodeUsecase.AddNode(node.toEntity())
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	// try the 1st time in 1 second, 2nd time in 2 seconds, 3rd in 4, 4th in 8, 5th in 16 seconds.
	waitInterval := 1 * time.Second
	retries := 5

	for retries != 0 {
		nodeInfo, err := handler.nodeUsecase.GetNode(result.ID)
		if err != nil {
			c.JSON(err.StatusCode(), err)
			return
		}

		if nodeInfo.Status == constant.NodeStatus.Posted || nodeInfo.Status == constant.NodeStatus.ValidationFailed || nodeInfo.Status == constant.NodeStatus.PostFailed {
			c.JSON(http.StatusOK, handler.toGetNodeVO(nodeInfo))
			return
		}

		time.Sleep(waitInterval)
		waitInterval *= 2
		retries--
	}

	// if server can't get the node with posted or failed information, return the node id for user to get the node in the future.
	c.JSON(http.StatusOK, handler.toAddNodeVO(result))
}

func (handler *nodeHandler) Validate(c *gin.Context) {
	var node interface{}

	if err := c.ShouldBindJSON(&node); err != nil {
		errors := jsonapi.NewError([]string{"JSON Error"}, []string{"Invalid JSON body."}, nil, []int{http.StatusBadRequest})
		res := jsonapi.Response(errors, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	jsonString, err := json.Marshal(node)
	if err != nil {
		errors := jsonapi.NewError([]string{"JSON Error"}, []string{"Cannot parse JSON body."}, nil, []int{http.StatusBadRequest})
		res := jsonapi.Response(errors, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	linkedSchemas, ok := getLinkedSchemas(node)
	if !ok {
		errors := jsonapi.NewError([]string{"Missing Required Property"}, []string{"The `linked_schemas` property is required."}, nil, []int{http.StatusBadRequest})
		res := jsonapi.Response(errors, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	// Validate against the default schema.
	linkedSchemas = append(linkedSchemas, "default-v2.0.0")

	// Validate against schemes specify inside the profile data.
	titles, details, sources, errorStatus := handler.validateAgainstSchemas(linkedSchemas, string(jsonString))
	if len(titles) != 0 {
		message := "Failed to validate against schemas: " + strings.Join(titles, " ")
		logger.Info(message)
		errors := jsonapi.NewError(titles, details, sources, errorStatus)
		res := jsonapi.Response(errors, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	meta := jsonapi.NewMeta("The submitted profile was validated successfully to its linked schemas.")
	res := jsonapi.Response(nil, meta)
	c.JSON(http.StatusOK, res)
}

func getLinkedSchemas(data interface{}) ([]string, bool) {
	json, ok := data.(map[string]interface{})
	if !ok {
		return nil, false
	}
	_, ok = json["linked_schemas"]
	if !ok {
		return nil, false
	}
	arrInterface, ok := json["linked_schemas"].([]interface{})
	if !ok {
		return nil, false
	}

	var linkedSchemas = make([]string, 0)

	for _, data := range arrInterface {
		linkedSchema, ok := data.(string)
		if !ok {
			return nil, false
		}
		linkedSchemas = append(linkedSchemas, linkedSchema)
	}

	return linkedSchemas, true
}

func (handler *nodeHandler) validateAgainstSchemas(linkedSchemas []string, validateData string) ([]string, []string, []string, []int) {
	var (
		titles, details, sources []string
		errorStatus              []int
	)

	for _, linkedSchema := range linkedSchemas {
		schemaURL := getSchemaURL(linkedSchema)

		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(schemaURL))
		if err != nil {
			titles = append(titles, []string{"Schema Not Found"}...)
			details = append(details, []string{"Could not locate the following schema in the library: " + linkedSchema}...)
			sources = append(sources, []string{"/linked_schemas"}...)
			errorStatus = append(errorStatus, http.StatusNotFound)
			continue
		}

		result, err := schema.Validate(gojsonschema.NewStringLoader(validateData))
		if err != nil {
			titles = append(titles, "Cannot Validate Document")
			details = append(details, []string{"Error when trying to validate document: ", err.Error()}...)
			errorStatus = append(errorStatus, http.StatusBadRequest)
			continue
		}

		if !result.Valid() {
			failedTitles, failedDetails, failedSources := handler.parseValidateError(linkedSchema, result.Errors())
			titles = append(titles, failedTitles...)
			details = append(details, failedDetails...)
			sources = append(sources, failedSources...)
			for i := 0; i < len(titles); i++ {
				errorStatus = append(errorStatus, http.StatusBadRequest)
			}
		}
	}

	return titles, details, sources, errorStatus
}

func getSchemaURL(linkedSchema string) string {
	return config.Conf.Library.InternalURL + "/v1/schema/" + linkedSchema
}

func (handler *nodeHandler) parseValidateError(schema string, resultErrors []gojsonschema.ResultError) ([]string, []string, []string) {
	var failedTitles, failedDetails, failedSources []string
	for _, desc := range resultErrors {
		// title
		failedType := desc.Type()

		// details
		var expected, given, min, max, property, failedDetail, failedField string
		for index, value := range desc.Details() {
			if index == "expected" {
				expected = value.(string)
			} else if index == "given" {
				given = value.(string)
			} else if index == "min" {
				min = fmt.Sprint(value)
			} else if index == "max" {
				max = fmt.Sprint(value)
			} else if index == "property" {
				property = value.(string)
			}
		}

		if failedType == "invalid_type" {
			failedType = "Invalid Type"
			failedDetail = "Expected: " + expected + " - Given: " + given + " - Schema: " + schema
		} else if failedType == "number_gte" {
			failedType = "Invalid Amount"
			failedDetail = "Amount must be greater than or equal to " + min + " - Schema: " + schema
		} else if failedType == "number_lte" {
			failedType = "Invalid Amount"
			failedDetail = "Amount must be less than or equal to " + max + " - Schema: " + schema
		} else if failedType == "required" {
			failedType = "Missing Required Property"
			failedDetail = "The `/" + desc.Field() + "/" + property + "` property is required."
		}

		// append title and detail
		failedTitles = append(failedTitles, failedType)
		failedDetails = append(failedDetails, failedDetail)

		// sources
		if property != "" {
			failedField = "/" + strings.Replace(desc.Field(), ".", "/", -1) + "/" + property
		} else {
			failedField = "/" + strings.Replace(desc.Field(), ".", "/", -1)
		}
		failedSources = append(failedSources, failedField)
	}

	return failedTitles, failedDetails, failedSources
}
