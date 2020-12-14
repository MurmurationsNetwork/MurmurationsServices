package http

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/service"
	"github.com/gin-gonic/gin"
)

type SchemaHandler interface {
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

func (handler *schemaHandler) Search(c *gin.Context) {
	searchRes, err := handler.svc.Search()
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, searchRes.Marshall())
}
