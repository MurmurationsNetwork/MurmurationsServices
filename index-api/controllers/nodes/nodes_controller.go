package nodes

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/index-api/domain/nodes"
	"github.com/MurmurationsNetwork/MurmurationsServices/index-api/services"
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/rest_errors"
	"github.com/gin-gonic/gin"
)

func AddNode(c *gin.Context) {
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

	c.JSON(http.StatusOK, result)
}

func GetNode(c *gin.Context) {
	nodeId, found := c.Params.Get("node_id")
	if !found {
		restErr := rest_errors.NewBadRequestError("invalid node id")
		c.JSON(restErr.Status(), restErr)
		return
	}

	node, err := services.NodeService.GetNode(nodeId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, node)
}

func DeleteNode(c *gin.Context) {
	c.String(http.StatusNotImplemented, "TODO")
}
