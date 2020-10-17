package nodes

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/indexer/internal/domain/nodes"
	"github.com/MurmurationsNetwork/MurmurationsServices/indexer/internal/services"
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/rest_errors"
	"github.com/gin-gonic/gin"
)

func getNodeId(params gin.Params) (string, rest_errors.RestErr) {
	nodeId, found := params.Get("node_id")
	if !found {
		return "", rest_errors.NewBadRequestError("invalid node id")
	}
	return nodeId, nil
}

func Add(c *gin.Context) {
	var node nodes.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, err := services.NodeService.AddNode(node)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result.AddNodeRespond())
}

func Get(c *gin.Context) {
	nodeId, err := getNodeId(c.Params)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	node, err := services.NodeService.GetNode(nodeId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, node)
}

func Search(c *gin.Context) {
	var query nodes.NodeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		restErr := rest_errors.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, err := services.NodeService.SearchNode(&query)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func Delete(c *gin.Context) {
	nodeId, err := getNodeId(c.Params)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	err = services.NodeService.DeleteNode(nodeId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
