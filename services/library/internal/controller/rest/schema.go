package rest

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/library"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/service"
)

// SchemaHandler defines the actions that can be performed with a Schema.
type SchemaHandler interface {
	Get(c *gin.Context)
	Search(c *gin.Context)
}

type schemaHandler struct {
	svc service.SchemaService
}

// NewSchemaHandler returns a new schemaHandler with the provided service.
func NewSchemaHandler(svc service.SchemaService) SchemaHandler {
	return &schemaHandler{
		svc: svc,
	}
}

// Get fetches a schema with a specific name.
func (handler *schemaHandler) Get(c *gin.Context) {
	schemaName, found := c.Params.Get("schemaName")
	// This normally won't happen, as if the user doesn't provide the name,
	// it will be considered a different API call.
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
		var schemaNotFoundError library.SchemaNotFoundError
		var dbError library.DatabaseError

		switch {
		case errors.As(err, &schemaNotFoundError):
			errors := jsonapi.NewError(
				[]string{"Schema Not Found"},
				[]string{schemaNotFoundError.Error()},
				nil,
				[]int{http.StatusNotFound},
			)
			res := jsonapi.Response(nil, errors, nil, nil)
			c.JSON(http.StatusNotFound, res)
		case errors.As(err, &dbError):
			errors := jsonapi.NewError(
				[]string{"Database Error"},
				[]string{dbError.Error()},
				nil,
				[]int{http.StatusInternalServerError},
			)
			res := jsonapi.Response(nil, errors, nil, nil)
			c.JSON(http.StatusInternalServerError, res)
		default:
			errors := jsonapi.NewError(
				[]string{"Unknown Error"},
				[]string{"An unexpected error has occurred."},
				nil,
				[]int{http.StatusInternalServerError},
			)
			res := jsonapi.Response(nil, errors, nil, nil)
			c.JSON(http.StatusInternalServerError, res)
		}
		return
	}

	c.JSON(http.StatusOK, schema)
}

// Search fetches all schemas that match the search criteria.
func (handler *schemaHandler) Search(c *gin.Context) {
	searchRes, err := handler.svc.Search()
	if err != nil {
		logger.Error("Error when trying to find schemas", err)
		errors := jsonapi.NewError(
			[]string{"Database Error"},
			[]string{"Error when trying to find schemas."},
			nil,
			[]int{http.StatusInternalServerError},
		)
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	res := jsonapi.Response(searchRes.Marshall(), nil, nil, nil)
	c.JSON(http.StatusOK, res)
}
