package http

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"net/http"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
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

		if nodeInfo.Status == constant.NodeStatus.ValidationFailed || nodeInfo.Status == constant.NodeStatus.PostFailed || nodeInfo.Status == constant.NodeStatus.Posted {
			c.JSON(http.StatusOK, handler.toGetNodeVO(nodeInfo))
			return
		}

		time.Sleep(waitInterval)
		retries--
	}

	// if server can't get the node with posted or failed information, return the node id for user to get the node in the future.
	c.JSON(http.StatusOK, handler.toAddNodeVO(result))
}
