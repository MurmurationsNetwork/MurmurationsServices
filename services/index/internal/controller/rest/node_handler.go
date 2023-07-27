package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/schemavalidator"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/index"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/repository/es"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/service"
)

type NodeHandler interface {
	Add(c *gin.Context)
	Get(c *gin.Context)
	Search(c *gin.Context)
	Delete(c *gin.Context)
	AddSync(c *gin.Context)
	Validate(c *gin.Context)
	Export(c *gin.Context)
	GetNodes(c *gin.Context)
}

type nodeHandler struct {
	svc service.NodeService
}

func NewNodeHandler(nodeService service.NodeService) NodeHandler {
	return &nodeHandler{
		svc: nodeService,
	}
}

var validationFields = []string{
	"name",
	"schema",
	"last_updated",
	"lat",
	"lon",
	"range",
	"locality",
	"region",
	"country",
	"status",
	"tags",
	"tags_filter",
	"tags_exact",
	"primary_url",
	"page",
	"page_size",
}

func (handler *nodeHandler) getNodeID(
	params gin.Params,
) (string, []jsonapi.Error) {
	nodeID, found := params.Get("nodeID")
	if !found {
		return "", jsonapi.NewError(
			[]string{"Invalid Node Id"},
			[]string{"The `node_id` is invalid."},
			nil,
			[]int{http.StatusBadRequest},
		)
	}
	return nodeID, nil
}

func (handler *nodeHandler) Add(c *gin.Context) {
	var req NodeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := jsonapi.NewError(
			[]string{"JSON Error"},
			[]string{"The JSON document submitted could not be parsed."},
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(err[0].Status, jsonapi.Response(nil, err, nil, nil))
		return
	}

	result, err := handler.svc.AddNode(&model.Node{
		ProfileURL: req.ProfileURL,
	})
	if err != nil {
		logger.Error("Failed to add node", err)

		var validationError index.ValidationError
		var jsonErr []jsonapi.Error

		switch {
		case errors.As(err, &validationError):
			jsonErr = jsonapi.NewError(
				[]string{},
				[]string{validationError.Reason},
				nil,
				[]int{http.StatusBadRequest},
			)
		default:
			jsonErr = jsonapi.NewError(
				[]string{"Database Error"},
				[]string{"Error when trying to add a node."},
				nil,
				[]int{http.StatusInternalServerError},
			)
		}

		res := jsonapi.Response(nil, jsonErr, nil, nil)
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := jsonapi.Response(ToAddNodeResponse(result), nil, nil, nil)
	c.JSON(http.StatusOK, res)
}

func (handler *nodeHandler) Get(c *gin.Context) {
	nodeID, jsonErr := handler.getNodeID(c.Params)
	if jsonErr != nil {
		res := jsonapi.Response(nil, jsonErr, nil, nil)
		c.JSON(jsonErr[0].Status, res)
		return
	}

	node, err := handler.svc.GetNode(nodeID)
	if err != nil {
		logger.Error("Failed to get a node", err)

		var notFoundError index.NotFoundError
		var databaseError index.DatabaseError
		var jsonErr []jsonapi.Error

		switch {
		case errors.As(err, &notFoundError):
			jsonErr = jsonapi.NewError(
				[]string{"Node Not Found"},
				[]string{
					fmt.Sprintf(
						"Could not locate the following node_id in the Index: %s",
						nodeID,
					),
				},
				nil,
				[]int{http.StatusNotFound},
			)
		case errors.As(err, &databaseError):
			jsonErr = jsonapi.NewError(
				[]string{databaseError.Message},
				[]string{"Error while trying to delete a node."},
				nil,
				[]int{http.StatusNotFound},
			)
		default:
			jsonErr = jsonapi.NewError(
				[]string{"Unknown Error"},
				[]string{},
				nil,
				[]int{http.StatusInternalServerError},
			)
		}

		res := jsonapi.Response(nil, jsonErr, nil, nil)
		c.JSON(jsonErr[0].Status, res)
		return
	}

	if node.Status == constant.NodeStatus.PostFailed {
		meta := jsonapi.NewMeta(
			"The system will automatically re-post the node, please check back in a minute.",
			"",
			"",
		)
		res := jsonapi.Response(ToGetNodeResponse(node), nil, nil, meta)
		c.JSON(http.StatusOK, res)
		return
	}

	if node.Status == constant.NodeStatus.ValidationFailed {
		meta := jsonapi.NewMeta("", node.ID, node.ProfileURL)
		errors := *node.FailureReasons
		res := jsonapi.Response(nil, errors, nil, meta)
		c.JSON(errors[0].Status, res)
		return
	}

	res := jsonapi.Response(ToGetNodeResponse(node), nil, nil, nil)
	c.JSON(http.StatusOK, res)
}

func (handler *nodeHandler) Search(c *gin.Context) {
	errs := checkInputIsValid(c, validationFields, "GET")
	if errs != nil {
		res := jsonapi.Response(nil, errs, nil, nil)
		c.JSON(errs[0].Status, res)
		return
	}

	var esQuery es.Query
	if err := c.ShouldBindQuery(&esQuery); err != nil {
		errs = jsonapi.NewError(
			[]string{"JSON Error"},
			[]string{"The JSON document submitted could not be parsed."},
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errs, nil, nil)
		c.JSON(errs[0].Status, res)
		return
	}

	if esQuery.Page*esQuery.PageSize > 10000 {
		errMsgs := []string{"Max Results Exceeded"}
		detailMsgs := []string{
			"No more than 10,000 results can be returned. " +
				"Refine your query so it will return less " +
				"but more relevant results.",
		}
		errs = jsonapi.NewError(
			errMsgs,
			detailMsgs,
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errs, nil, nil)
		c.JSON(errs[0].Status, res)
		return
	}

	searchResult, err := handler.svc.Search(&esQuery)
	if err != nil {
		logger.Error("Failed to search a node", err)

		var databaseError index.DatabaseError
		var jsonErr []jsonapi.Error

		switch {
		case errors.As(err, &databaseError):
			jsonErr = jsonapi.NewError(
				[]string{databaseError.Message},
				[]string{"Error while trying to search a node."},
				nil,
				[]int{http.StatusNotFound},
			)
		default:
			jsonErr = jsonapi.NewError(
				[]string{"Unknown Error"},
				[]string{},
				nil,
				[]int{http.StatusInternalServerError},
			)
		}

		res := jsonapi.Response(nil, jsonErr, nil, nil)
		c.JSON(jsonErr[0].Status, res)
		return
	}

	// restrict the last page to the page of 10,000 results (ES limitation)
	totalPage := 10000 / esQuery.PageSize
	message := "No more than 10,000 results can be returned. " +
		"Refine your query so it will return less " +
		"but more relevant results."
	if totalPage >= searchResult.TotalPages {
		totalPage = searchResult.TotalPages
		message = ""
	}
	// edge case: page = 0 or larger than total page - response no data
	if searchResult.TotalPages == 0 || esQuery.Page > searchResult.TotalPages {
		res := jsonapi.Response(searchResult.Result, nil, nil, nil)
		c.JSON(http.StatusOK, res)
		return
	}
	meta := jsonapi.NewSearchMeta(
		message,
		searchResult.NumberOfResults,
		searchResult.TotalPages,
	)
	links := jsonapi.NewLinks(c, esQuery.Page, totalPage)
	res := jsonapi.Response(searchResult.Result, nil, links, meta)
	c.JSON(http.StatusOK, res)
}

func (handler *nodeHandler) Delete(c *gin.Context) {
	if c.Params.ByName("nodeID") == "" {
		errors := jsonapi.NewError(
			[]string{"Missing Path Parameter"},
			[]string{"The `node_id` path parameter is missing."},
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	nodeID, jsonErr := handler.getNodeID(c.Params)
	if jsonErr != nil {
		res := jsonapi.Response(nil, jsonErr, nil, nil)
		c.JSON(jsonErr[0].Status, res)
		return
	}

	profileURL, err := handler.svc.Delete(nodeID)
	if err != nil {
		logger.Error("Failed to delete a node", err)

		var deleteNodeError index.DeleteNodeError
		var databaseError index.DatabaseError
		var jsonErr []jsonapi.Error

		switch {
		case errors.As(err, &deleteNodeError):
			jsonErr = jsonapi.NewError(
				[]string{deleteNodeError.Message},
				[]string{deleteNodeError.Detail},
				nil,
				[]int{http.StatusBadRequest},
			)
		case errors.As(err, &databaseError):
			jsonErr = jsonapi.NewError(
				[]string{databaseError.Message},
				[]string{"Error while trying to delete a node."},
				nil,
				[]int{http.StatusInternalServerError},
			)
		default:
			jsonErr = jsonapi.NewError(
				[]string{"Unknown Error"},
				[]string{},
				nil,
				[]int{http.StatusInternalServerError},
			)
		}

		meta := jsonapi.NewMeta("", nodeID, profileURL)
		res := jsonapi.Response(nil, jsonErr, nil, meta)
		c.JSON(jsonErr[0].Status, res)
		return
	}

	deleteTTL := dateutil.FormatSeconds(config.Values.TTL.DeletedTTL)

	meta := jsonapi.NewMeta(
		fmt.Sprintf(
			"The Index has recorded as deleted the profile that was previously "+
				"posted at: %s -- It will be completely removed from the index in %s.",
			profileURL,
			deleteTTL,
		),
		"",
		"",
	)
	res := jsonapi.Response(nil, nil, nil, meta)
	c.JSON(http.StatusOK, res)
}

func (handler *nodeHandler) AddSync(c *gin.Context) {
	var req NodeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := jsonapi.NewError(
			[]string{"JSON Error"},
			[]string{"The JSON document submitted could not be parsed."},
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	if err := req.Validate(); err != nil {
		res := jsonapi.Response(nil, err, nil, nil)
		c.JSON(err[0].Status, res)
		return
	}

	result, err := handler.svc.AddNode(&model.Node{
		ProfileURL: req.ProfileURL,
	})
	if err != nil {
		logger.Error("Failed to add node", err)

		var validationError index.ValidationError
		var jsonErr []jsonapi.Error

		switch {
		case errors.As(err, &validationError):
			jsonErr = jsonapi.NewError(
				[]string{"Missing Required Property"},
				[]string{validationError.Reason},
				nil,
				[]int{http.StatusBadRequest},
			)
		default:
			jsonErr = jsonapi.NewError(
				[]string{"Database Error"},
				[]string{"Error when trying to add a node."},
				nil,
				[]int{http.StatusInternalServerError},
			)
		}

		res := jsonapi.Response(nil, jsonErr, nil, nil)
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// try the 1st time in 1 second, 2nd time in 2 seconds, 3rd in 4, 4th in 8, 5th in 16 seconds.
	waitInterval := 1 * time.Second
	retries := 5

	for retries != 0 {
		nodeInfo, err := handler.svc.GetNode(result.ID)
		if err != nil {
			logger.Error("Failed to get a node", err)

			var notFoundError index.NotFoundError
			var databaseError index.DatabaseError
			var jsonErr []jsonapi.Error

			switch {
			case errors.As(err, &notFoundError):
				jsonErr = jsonapi.NewError(
					[]string{"Node Not Found"},
					[]string{
						fmt.Sprintf(
							"Could not locate the following node_id in the Index: %s",
							result.ID,
						),
					},
					nil,
					[]int{http.StatusNotFound},
				)
			case errors.As(err, &databaseError):
				jsonErr = jsonapi.NewError(
					[]string{databaseError.Message},
					[]string{"Error while trying to delete a node."},
					nil,
					[]int{http.StatusNotFound},
				)
			default:
				jsonErr = jsonapi.NewError(
					[]string{"Unknown Error"},
					[]string{},
					nil,
					[]int{http.StatusInternalServerError},
				)
			}

			res := jsonapi.Response(nil, jsonErr, nil, nil)
			c.JSON(jsonErr[0].Status, res)
			return
		}

		if nodeInfo.Status == constant.NodeStatus.PostFailed {
			meta := jsonapi.NewMeta(
				"The system will automatically re-post the node, please check back in a minute.",
				"",
				"",
			)
			res := jsonapi.Response(ToGetNodeResponse(result), nil, nil, meta)
			c.JSON(http.StatusOK, res)
			return
		}

		if nodeInfo.Status == constant.NodeStatus.ValidationFailed {
			meta := jsonapi.NewMeta("", nodeInfo.ID, nodeInfo.ProfileURL)
			errors := *nodeInfo.FailureReasons
			res := jsonapi.Response(nil, errors, nil, meta)
			c.JSON(errors[0].Status, res)
			return
		}

		if nodeInfo.Status == constant.NodeStatus.Posted ||
			nodeInfo.Status == constant.NodeStatus.Deleted {
			res := jsonapi.Response(
				ToGetNodeResponse(nodeInfo),
				nil,
				nil,
				nil,
			)
			c.JSON(http.StatusOK, res)
			return
		}

		time.Sleep(waitInterval)
		waitInterval *= 2
		retries--
	}

	// If server can't get the node with posted or failed information, return
	// the node id for user to get the node in the future.
	res := jsonapi.Response(ToAddNodeResponse(result), nil, nil, nil)
	c.JSON(http.StatusOK, res)
}

func (handler *nodeHandler) Validate(c *gin.Context) {
	var node interface{}

	if err := c.ShouldBindJSON(&node); err != nil {
		errors := jsonapi.NewError(
			[]string{"JSON Error"},
			[]string{"The JSON document submitted could not be parsed."},
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	jsonString, err := json.Marshal(node)
	if err != nil {
		errors := jsonapi.NewError(
			[]string{"JSON Error"},
			[]string{"The body of the JSON document submitted is malformed."},
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	linkedSchemas, ok := getLinkedSchemas(node)
	if !ok {
		errors := jsonapi.NewError(
			[]string{"Missing Required Property"},
			[]string{"The `linked_schemas` property is required."},
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	// Validate against the default schema.
	linkedSchemas = append(linkedSchemas, "default-v2.0.0")

	// Validate against schemes specify inside the profile data.
	validator, err := schemavalidator.NewBuilder().
		WithURLSchemas(config.Values.Library.InternalURL, linkedSchemas).
		WithStrProfileLoader(string(jsonString)).
		Build()
	if err != nil {
		// Log the error for internal debugging and auditing.
		logger.Error("Failed to build schema validator", err)

		errors := jsonapi.NewError(
			[]string{"Internal Server Error"},
			[]string{
				"An error occurred while validating the profile data. Please try again later.",
			},
			nil,
			[]int{http.StatusInternalServerError},
		)
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	result := validator.Validate()
	if !result.Valid {
		message := "Failed to validate against schemas: " + strings.Join(
			result.ErrorMessages,
			" ",
		)
		logger.Info(message)
		errors := jsonapi.NewError(
			result.ErrorMessages,
			result.Details,
			result.Sources,
			result.ErrorStatus,
		)
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	meta := jsonapi.NewMeta(
		"The submitted profile was validated successfully to its linked schemas.",
		"",
		"",
	)
	res := jsonapi.Response(nil, nil, nil, meta)
	c.JSON(http.StatusOK, res)
}

func (handler *nodeHandler) Export(c *gin.Context) {
	// return error if there is an invalid query
	// get the fields from query.EsQuery
	fields := []string{"schema", "page_size", "search_after"}
	errs := checkInputIsValid(c, fields, "POST")
	if errs != nil {
		res := jsonapi.Response(nil, errs, nil, nil)
		c.JSON(errs[0].Status, res)
		return
	}

	var esQuery es.BlockQuery
	if err := c.ShouldBindJSON(&esQuery); err != nil {
		fmt.Println(err)
		errs = jsonapi.NewError(
			[]string{"JSON Error"},
			[]string{"The JSON document submitted could not be parsed."},
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errs, nil, nil)
		c.JSON(errs[0].Status, res)
		return
	}

	// set default page_size for esQuery
	if esQuery.PageSize == 0 {
		esQuery.PageSize = 100
	}

	searchResult, err := handler.svc.Export(&esQuery)
	if err != nil {
		logger.Error("Failed to export a node", err)

		var databaseError index.DatabaseError
		var jsonErr []jsonapi.Error

		switch {
		case errors.As(err, &databaseError):
			jsonErr = jsonapi.NewError(
				[]string{databaseError.Message},
				[]string{"Error while trying to delete a node."},
				nil,
				[]int{http.StatusInternalServerError},
			)
		default:
			jsonErr = jsonapi.NewError(
				[]string{"Unknown Error"},
				[]string{},
				nil,
				[]int{http.StatusInternalServerError},
			)
		}

		res := jsonapi.Response(nil, jsonErr, nil, nil)
		c.JSON(jsonErr[0].Status, res)
		return
	}

	meta := jsonapi.NewBlockSearchMeta(searchResult.Sort)
	res := jsonapi.Response(searchResult.Result, nil, nil, meta)
	c.JSON(http.StatusOK, res)
}

func (handler *nodeHandler) GetNodes(c *gin.Context) {
	errs := checkInputIsValid(c, validationFields, "GET")
	if errs != nil {
		res := jsonapi.Response(nil, errs, nil, nil)
		c.JSON(errs[0].Status, res)
		return
	}

	var esQuery es.Query
	if err := c.ShouldBindQuery(&esQuery); err != nil {
		errs = jsonapi.NewError(
			[]string{"JSON Error"},
			[]string{"The JSON document submitted could not be parsed."},
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errs, nil, nil)
		c.JSON(errs[0].Status, res)
		return
	}

	if esQuery.Page*esQuery.PageSize > 10000 {
		msg := "No more than 10,000 results can be returned. " +
			"Refine your query so it will return less " +
			"but more relevant results."
		errs = jsonapi.NewError(
			[]string{"Max Results Exceeded"},
			[]string{msg},
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errs, nil, nil)
		c.JSON(errs[0].Status, res)
		return
	}

	searchResult, err := handler.svc.GetNodes(&esQuery)
	if err != nil {
		logger.Error("Failed to get a node", err)

		var notFoundError index.NotFoundError
		var databaseError index.DatabaseError
		var jsonErr []jsonapi.Error

		switch {
		case errors.As(err, &notFoundError):
			jsonErr = jsonapi.NewError(
				[]string{"Node Not Found"},
				[]string{},
				nil,
				[]int{http.StatusNotFound},
			)
		case errors.As(err, &databaseError):
			jsonErr = jsonapi.NewError(
				[]string{databaseError.Message},
				[]string{"Error while trying to delete a node."},
				nil,
				[]int{http.StatusNotFound},
			)
		default:
			jsonErr = jsonapi.NewError(
				[]string{"Unknown Error"},
				[]string{},
				nil,
				[]int{http.StatusInternalServerError},
			)
		}

		res := jsonapi.Response(nil, jsonErr, nil, nil)
		c.JSON(jsonErr[0].Status, res)
		return
	}

	// restrict the last page to the page of 10,000 results (ES limitation)
	totalPage := 10000 / esQuery.PageSize
	message := "No more than 10,000 results can be returned. " +
		"Refine your query so it will return less " +
		"but more relevant results."
	if totalPage >= searchResult.TotalPages {
		totalPage = searchResult.TotalPages
		message = ""
	}
	// edge case: page = 0 or larger than total page - response no data
	if searchResult.TotalPages == 0 || esQuery.Page > searchResult.TotalPages {
		res := jsonapi.Response(searchResult.Result, nil, nil, nil)
		c.JSON(http.StatusOK, res)
		return
	}

	meta := jsonapi.NewSearchMeta(
		message,
		searchResult.NumberOfResults,
		searchResult.TotalPages,
	)
	links := jsonapi.NewLinks(c, esQuery.Page, totalPage)
	res := jsonapi.Response(searchResult.Result, nil, links, meta)
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

func checkInputIsValid(
	c *gin.Context,
	fields []string,
	requestType string,
) []jsonapi.Error {
	// return error if there is an invalid query
	// get the fields from query.EsQuery
	var queryFields map[string]interface{}
	if requestType == "POST" {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			errors := jsonapi.NewError(
				[]string{"JSON Error"},
				[]string{"The JSON document submitted could not be parsed."},
				nil,
				[]int{http.StatusBadRequest},
			)
			return errors
		}

		err = json.Unmarshal(body, &queryFields)
		if err != nil {
			errors := jsonapi.NewError(
				[]string{"JSON Error"},
				[]string{"The JSON document submitted could not be unmarshal."},
				nil,
				[]int{http.StatusBadRequest},
			)
			return errors
		}

		// restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	} else {
		queryMap := c.Request.URL.Query()
		queryFields = make(map[string]interface{})

		for key, value := range queryMap {
			if len(value) == 1 {
				queryFields[key] = value[0]
			} else {
				queryFields[key] = value
			}
		}
	}

	var (
		invalidQueryTitles, invalidQueryDetails []string
		invalidQuerySources                     [][]string
		invalidQueryStatus                      []int
	)
	for fieldName := range queryFields {
		found := false
		for _, validFieldName := range fields {
			if fieldName == validFieldName {
				found = true
				break
			}
		}
		if !found {
			invalidQueryTitles = append(
				invalidQueryTitles,
				"Invalid Query Parameter",
			)
			invalidQueryDetails = append(
				invalidQueryDetails,
				fmt.Sprintf(
					"The following query parameter is not valid: %v",
					fieldName,
				),
			)
			invalidQuerySources = append(
				invalidQuerySources,
				[]string{"parameter", fieldName},
			)
			invalidQueryStatus = append(
				invalidQueryStatus,
				http.StatusBadRequest,
			)
		}
	}

	if len(invalidQueryTitles) != 0 {
		errors := jsonapi.NewError(
			invalidQueryTitles,
			invalidQueryDetails,
			invalidQuerySources,
			invalidQueryStatus,
		)
		return errors
	}

	return nil
}
