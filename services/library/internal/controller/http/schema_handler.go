package http

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"

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
		errors := jsonapi.NewError(
			[]string{"Invalid Schema Name"},
			[]string{"The schema name is not valid."},
			nil,
			[]int{http.StatusBadRequest},
		)
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}
	schema, err := handler.svc.Get(schemaName)
	if err != nil {
		res := jsonapi.Response(nil, err, nil, nil)
		c.JSON(err[0].Status, res)
		return
	}

	c.JSON(http.StatusOK, schema)
}

func (handler *schemaHandler) Search(c *gin.Context) {
	searchRes, err := handler.svc.Search()
	if err != nil {
		res := jsonapi.Response(nil, err, nil, nil)
		c.JSON(err[0].Status, res)
		return
	}

	res := jsonapi.Response(searchRes.Marshall(), nil, nil, nil)
	c.JSON(http.StatusOK, res)
}
