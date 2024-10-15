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
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/core"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/profile/profilevalidator"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/index"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/repository/es"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/service"
)

type NodeHandler interface {
	// Add creates a new node.
	Add(c *gin.Context)
	// AddSync creates a new node synchronously.
	AddSync(c *gin.Context)
	// Get retrieves a specific node.
	Get(c *gin.Context)
	// GetNodes retrieves multiple nodes.
	GetNodes(c *gin.Context)
	// Search finds nodes that match certain criteria.
	Search(c *gin.Context)
	// Delete removes a node.
	Delete(c *gin.Context)
	// Validate validates a node.
	Validate(c *gin.Context)
	// Export exports nodes.
	Export(c *gin.Context)
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
	"expires",
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
		handleAddNodeErrors(c, err)
		return
	}

	res := jsonapi.Response(ToAddNodeResponse(result), nil, nil, nil)
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
		handleAddNodeErrors(c, err)
		return
	}

	// SERVER_TIMEOUT_WRITE is 15 seconds, so we can't wait for more than that.
	// Try the 1st time in 1 second, 2nd time in 2 seconds, 3rd in 4
	waitInterval := 1 * time.Second
	retries := 3

	for retries != 0 {
		nodeInfo, err := handler.svc.GetNode(result.ID)
		if err != nil {
			handleGetNodeErrors(c, err, &result.ID)
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

func (handler *nodeHandler) Get(c *gin.Context) {
	nodeID, jsonErr := handler.getNodeID(c.Params)
	if jsonErr != nil {
		res := jsonapi.Response(nil, jsonErr, nil, nil)
		c.JSON(jsonErr[0].Status, res)
		return
	}

	node, err := handler.svc.GetNode(nodeID)
	if err != nil {
		handleGetNodeErrors(c, err, &nodeID)
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
				[]string{
					"An unexpected error occurred. Please try again later.",
				},
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
		handleDeleteNodeErrors(c, err, nodeID, profileURL)
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
	linkedSchemas = append(linkedSchemas, "default-v2.1.0")

	// Validate against schemes specify inside the profile data.
	validator, err := profilevalidator.NewBuilder().
		WithURLSchemas(config.Values.Library.InternalURL, linkedSchemas).
		WithStrProfile(string(jsonString)).
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

	// Validate the "expires" field
	if expires, ok := node.(map[string]interface{})["expires"].(float64); ok {
		expireTime := time.Unix(int64(expires), 0)
		if expireTime.Before(time.Now()) {
			errors := jsonapi.NewError(
				[]string{"Invalid Expires Field"},
				[]string{"The `expires` field has already passed."},
				nil,
				[]int{http.StatusBadRequest},
			)
			res := jsonapi.Response(nil, errors, nil, nil)
			c.JSON(errors[0].Status, res)
			return
		}
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
		logger.Error("Failed to export nodes", err)

		var databaseError index.DatabaseError
		var jsonErr []jsonapi.Error

		switch {
		case errors.As(err, &databaseError):
			jsonErr = jsonapi.NewError(
				[]string{databaseError.Message},
				[]string{"Error while trying to export nodes."},
				nil,
				[]int{http.StatusInternalServerError},
			)
		default:
			jsonErr = jsonapi.NewError(
				[]string{"Unknown Error"},
				[]string{
					"An unexpected error occurred. Please try again later.",
				},
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
		handleGetNodeErrors(c, err, nil)
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

func handleAddNodeErrors(c *gin.Context, err error) {
	var validationError index.ValidationError
	var profileFetchError core.ProfileFetchError
	var jsonErr []jsonapi.Error

	switch {
	case errors.As(err, &validationError):
		jsonErr = jsonapi.NewError(
			[]string{"Validation Error"},
			[]string{validationError.Reason},
			nil,
			[]int{http.StatusBadRequest},
		)

	case errors.As(err, &profileFetchError):
		jsonErr = jsonapi.NewError(
			[]string{"Profile Fetch Error"},
			[]string{profileFetchError.Reason},
			nil,
			[]int{http.StatusNotFound},
		)

	default:
		jsonErr = jsonapi.NewError(
			[]string{"Unknown Error"},
			[]string{"An unexpected error occurred. Please try again later."},
			nil,
			[]int{http.StatusInternalServerError},
		)
	}

	res := jsonapi.Response(nil, jsonErr, nil, nil)
	c.JSON(jsonErr[0].Status, res)
}

func handleGetNodeErrors(c *gin.Context, err error, nodeID *string) {
	var notFoundError index.NotFoundError
	var databaseError index.DatabaseError
	var jsonErr []jsonapi.Error

	if errors.As(err, &notFoundError) {
		nodeIDMsg := "a node"
		if nodeID != nil {
			nodeIDMsg = fmt.Sprintf(
				"the following node_id in the Index: %s",
				*nodeID,
			)
		}
		jsonErr = jsonapi.NewError(
			[]string{"Node Not Found"},
			[]string{fmt.Sprintf("Could not locate %s", nodeIDMsg)},
			nil,
			[]int{http.StatusNotFound},
		)
	} else {
		logger.Error("Failed to get a node", err)

		switch {
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
				[]string{"An unexpected error occurred. Please try again later."},
				nil,
				[]int{http.StatusInternalServerError},
			)
		}
	}

	res := jsonapi.Response(nil, jsonErr, nil, nil)
	c.JSON(jsonErr[0].Status, res)
}

func handleDeleteNodeErrors(
	c *gin.Context,
	err error,
	nodeID, profileURL string,
) {
	var (
		notFoundError   index.NotFoundError
		databaseError   index.DatabaseError
		deleteNodeError index.DeleteNodeError
		jsonErr         []jsonapi.Error
		errMsg          string
		detailMsg       string
		statusCode      int
		logAsError      = true
	)

	switch {
	case errors.As(err, &notFoundError):
		errMsg = "Node Not Found"
		statusCode = http.StatusNotFound
		logAsError = false

	case errors.As(err, &deleteNodeError):
		statusCode = http.StatusBadRequest
		errMsg = deleteNodeError.Message
		detailMsg = deleteNodeError.Detail
		switch deleteNodeError.ErrorCode {
		case index.ErrorHTTPRequestFailed,
			index.ErrorProfileURLCheckFail,
			index.ErrorProfileStillExists:
			logAsError = false
		default:
		}

	case errors.As(err, &databaseError):
		errMsg = databaseError.Message
		detailMsg = "Error while trying to delete a node."
		statusCode = http.StatusInternalServerError

	default:
		errMsg = "Unknown Error"
		statusCode = http.StatusInternalServerError
	}

	if logAsError {
		logger.Error("Failed to delete a node", err)
	} else {
		logger.Info(fmt.Sprintf("Info on node deletion: %s - %s", errMsg, detailMsg))
	}

	details := []string{}
	if detailMsg != "" {
		details = append(details, detailMsg)
	}
	jsonErr = jsonapi.NewError(
		[]string{errMsg},
		details,
		nil,
		[]int{statusCode},
	)

	meta := jsonapi.NewMeta("", nodeID, profileURL)
	res := jsonapi.Response(nil, jsonErr, nil, meta)
	c.JSON(statusCode, res)
}
