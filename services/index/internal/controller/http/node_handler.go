package http

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/query"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/service"
	"github.com/gin-gonic/gin"
)

type NodeHandler interface {
	Add(c *gin.Context)
	Get(c *gin.Context)
	Search(c *gin.Context)
	Delete(c *gin.Context)
}

type nodeHandler struct {
	nodeService service.NodesService
}

func NewNodeHandler(nodeService service.NodesService) NodeHandler {
	return &nodeHandler{
		nodeService: nodeService,
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
	var node node.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body.")
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, err := handler.nodeService.AddNode(&node)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result.AddNodeRespond())
}

func (handler *nodeHandler) Get(c *gin.Context) {
	nodeId, err := handler.getNodeId(c.Params)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	node, err := handler.nodeService.GetNode(nodeId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, node.GetNodeRespond())
}

func (handler *nodeHandler) Search(c *gin.Context) {
	var query query.EsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body.")
		c.JSON(restErr.Status(), restErr)
		return
	}

	searchRes, err := handler.nodeService.Search(&query)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, searchRes.Marshall())
}

func (handler *nodeHandler) Delete(c *gin.Context) {
	nodeId, err := handler.getNodeId(c.Params)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	err = handler.nodeService.Delete(nodeId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}
