package http

import (
	"encoding/json"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
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
		restErr := resterr.NewBadRequestError("Invalid JSON body.")
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	jsonString, err := json.Marshal(node)
	if err != nil {
		restErr := resterr.NewBadRequestError("Cannot parse JSON body.")
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	linkedSchemas, ok := getLinkedSchemas(node)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"failure_reasons": []string{"The submitted profile does not contain the linked_schemas property."},
			"status":          http.StatusBadRequest,
		})
		return
	}

	// Validate against the default schema.
	linkedSchemas = append(linkedSchemas, "default-v2.0.0")

	// Validate against schemes specify inside the profile data.
	failureReasons, errorStatus := handler.validateAgainstSchemas(linkedSchemas, string(jsonString))
	if len(failureReasons) != 0 {
		message := "Failed to validate against schemas: " + strings.Join(failureReasons, " ")
		logger.Info(message)
		c.JSON(http.StatusOK, gin.H{
			"failure_reasons": failureReasons,
			"status":          errorStatus,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
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

func (handler *nodeHandler) validateAgainstSchemas(linkedSchemas []string, validateData string) ([]string, int) {
	FailureReasons := []string{}
	errorStatus := 0

	for _, linkedSchema := range linkedSchemas {
		schemaURL := getSchemaURL(linkedSchema)

		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(schemaURL))
		if err != nil {
			FailureReasons = append(FailureReasons, fmt.Sprintf("Error when trying to read from schema %s: %s", schemaURL, err.Error()))
			if errorStatus == 0 {
				errorStatus = http.StatusNotFound
			}
			continue
		}

		result, err := schema.Validate(gojsonschema.NewStringLoader(validateData))
		if err != nil {
			FailureReasons = append(FailureReasons, "Error when trying to validate document: ", err.Error())
			if errorStatus == 0 {
				errorStatus = http.StatusBadRequest
			}
			continue
		}

		if !result.Valid() {
			FailureReasons = append(FailureReasons, handler.parseValidateError(linkedSchema, result.Errors())...)
			if errorStatus == 0 {
				errorStatus = http.StatusBadRequest
			}
		}
	}

	return FailureReasons, errorStatus
}

func getSchemaURL(linkedSchema string) string {
	return config.Conf.Library.URL + "/schemas/" + linkedSchema + ".json"
}

func (handler *nodeHandler) parseValidateError(schema string, resultErrors []gojsonschema.ResultError) []string {
	FailureReasons := make([]string, 0)
	for _, desc := range resultErrors {
		// Output string: "demo-v1.(root): url is required"
		FailureReasons = append(FailureReasons, schema+"."+desc.String())
	}
	return FailureReasons
}
