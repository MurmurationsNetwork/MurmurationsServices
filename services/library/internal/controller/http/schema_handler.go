package http

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/service"
	"github.com/gin-gonic/gin"
)

type SchemaHandler interface {
	Get(c *gin.Context)
	Search(c *gin.Context)
}

type schemaHandler struct {
	svc service.SchemaService
}

func NewSchemaHandler(svc service.SchemaService) SchemaHandler {
	return &schemaHandler{
		svc: svc,
	}
}

func (handler *schemaHandler) Get(c *gin.Context) {
	schemaName, found := c.Params.Get("schemaName")
	if !found {
		c.JSON(http.StatusBadRequest, "Invalid schemaName.")
		return
	}
	schema, err := handler.svc.Get(schemaName)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, schema)
}

func (handler *schemaHandler) Search(c *gin.Context) {
	searchRes, err := handler.svc.Search()
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, searchRes.Marshall())
}
