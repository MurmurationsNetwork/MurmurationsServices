package node

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/service"
	"github.com/gin-gonic/gin"
)

func getNodeId(params gin.Params) (string, resterr.RestErr) {
	nodeId, found := params.Get("node_id")
	if !found {
		return "", resterr.NewBadRequestError("invalid node id")
	}
	return nodeId, nil
}

func Add(c *gin.Context) {
	var node node.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		restErr := resterr.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, err := service.NodeService.AddNode(node)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result.Marshall())
}

func Search(c *gin.Context) {
	var query node.NodeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		restErr := resterr.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}

	searchRes, err := service.NodeService.SearchNode(&query)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, searchRes.Marshall())
}

func Delete(c *gin.Context) {
	nodeId, err := getNodeId(c.Params)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	err = service.NodeService.DeleteNode(nodeId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
