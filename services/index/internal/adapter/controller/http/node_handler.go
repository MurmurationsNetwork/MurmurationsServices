package http

import (
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
		c.JSON(restErr.Status(), restErr)
		return
	}

	if err := node.Validate(); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	result, err := handler.nodeUsecase.AddNode(node.toEntity())
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, handler.toAddNodeVO(result))
}

func (handler *nodeHandler) Get(c *gin.Context) {
	nodeId, err := handler.getNodeId(c.Params)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	node, err := handler.nodeUsecase.GetNode(nodeId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, handler.toGetNodeVO(node))
}

func (handler *nodeHandler) Search(c *gin.Context) {
	var query query.EsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body.")
		c.JSON(restErr.Status(), restErr)
		return
	}

	searchResult, err := handler.nodeUsecase.Search(&query)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, searchResult.ToVO())
}

func (handler *nodeHandler) Delete(c *gin.Context) {
	nodeId, err := handler.getNodeId(c.Params)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	err = handler.nodeUsecase.Delete(nodeId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (handler *nodeHandler) AddSync(c *gin.Context) {
	var node nodeDTO
	if err := c.ShouldBindJSON(&node); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body.")
		c.JSON(restErr.Status(), restErr)
		return
	}

	if err := node.Validate(); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	result, err := handler.nodeUsecase.AddNode(node.toEntity())
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	// default time interval is 5 seconds
	waitInterval := 5 * time.Second
	retries := 5

	for retries != 0 {
		nodeInfo, err := handler.nodeUsecase.GetNode(result.ID)
		if err != nil {
			c.JSON(err.Status(), err)
			return
		}

		if nodeInfo.Status == constant.NodeStatus.Posted {
			c.JSON(http.StatusOK, handler.toGetNodeVO(nodeInfo))
			return
		}

		if nodeInfo.Status == constant.NodeStatus.ValidationFailed || nodeInfo.Status == constant.NodeStatus.PostFailed {
			c.JSON(http.StatusBadRequest, handler.toGetNodeVO(nodeInfo))
			return
		}

		time.Sleep(waitInterval)
		retries--
	}

	// if server can't get the node with posted or failed information, return the node id for user to get the node in the future.
	c.JSON(http.StatusOK, handler.toAddNodeVO(result))
}

func (handler *nodeHandler) Validate(c *gin.Context) {
	var node interface{}

	if err := c.ShouldBindJSON(&node); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body.")
		c.JSON(restErr.Status(), restErr)
		return
	}

	linkedSchemas, ok := getLinkedSchemas(node)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"status":          "validation_failed",
				"failure_reasons": "Could not read linked_schemas",
			},
		})
		return
	}

	// Validate against schemes specify inside the profile data.
	failureReasons := handler.validateAgainstSchemas(linkedSchemas, node)
	if len(failureReasons) != 0 {
		message := "Failed to validate against schemas: " + strings.Join(failureReasons, " ")
		logger.Info(message)
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"status":          "validation_failed",
				"failure_reasons": message,
			},
		})
		return
	}

	c.Status(http.StatusOK)
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

func (handler *nodeHandler) validateAgainstSchemas(linkedSchemas []string, validateData interface{}) []string {
	FailureReasons := []string{}

	for _, linkedSchema := range linkedSchemas {
		schemaURL := getSchemaURL(linkedSchema)

		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(schemaURL))
		if err != nil {
			FailureReasons = append(FailureReasons, fmt.Sprintf("Error when trying to read from schema %s: %s", schemaURL, err.Error()))
			continue
		}

		result, err := schema.Validate(gojsonschema.NewRawLoader(validateData))
		if err != nil {
			FailureReasons = append(FailureReasons, "Error when trying to validate document: ", err.Error())
			continue
		}

		if !result.Valid() {
			FailureReasons = append(FailureReasons, handler.parseValidateError(linkedSchema, result.Errors())...)
		}
	}

	return FailureReasons
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
