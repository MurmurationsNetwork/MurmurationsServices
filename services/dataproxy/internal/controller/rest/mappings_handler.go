package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/repository/mongo"
)

type MappingsHandler interface {
	Create(c *gin.Context)
}

type mappingsHandler struct {
	mappingRepository mongo.MappingRepository
}

func NewMappingsHandler(
	mappingRepository mongo.MappingRepository,
) MappingsHandler {
	return &mappingsHandler{
		mappingRepository: mappingRepository,
	}
}

func (handler *mappingsHandler) Create(c *gin.Context) {
	var mappings map[string]interface{}
	if err := c.BindJSON(&mappings); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body.")
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	schema := mappings["schema"]
	if schema == nil {
		restErr := resterr.NewBadRequestError("Schema is required.")
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	// download mapping from cdn
	url := fmt.Sprintf(
		"%s/v2/schemas/%s",
		config.Values.Library.InternalURL,
		schema,
	)
	//var schemas map[string]interface{}
	bytes, err := httputil.GetByte(url)
	if err != nil {
		restErr := resterr.NewInternalServerError(
			"Library retrieved Failed: ",
			err,
		)
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	var schemas map[string]interface{}
	err = json.Unmarshal(bytes, &schemas)
	if err != nil {
		restErr := resterr.NewInternalServerError(
			"Unmarshal Schema Failed: ",
			err,
		)
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	if schemas["properties"] == nil {
		restErr := resterr.NewBadRequestError(
			"Schema's Properties can't be found.",
		)
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	properties := make(map[string]string)
	schemaProperties := schemas["properties"].(map[string]interface{})
	for index := range schemaProperties {
		if index == "linked_schemas" {
			continue
		}
		properties[index] = ""
	}

	for index, value := range mappings {
		// check name is in the properties or not
		if index == "schema" {
			properties["schema"] = value.(string)
			continue
		}
		if _, ok := properties[index]; ok {
			properties[index] = value.(string)
		} else {
			restErr := resterr.NewBadRequestError("The property " + index + " can't be found")
			c.JSON(restErr.StatusCode(), restErr)
			return
		}
	}

	for index, value := range properties {
		if value == "" {
			restErr := resterr.NewBadRequestError(
				"The property " + index + " can't be blank.",
			)
			c.JSON(restErr.StatusCode(), restErr)
			return
		}
	}

	restErr := handler.mappingRepository.Save(properties)
	if restErr != nil {
		c.JSON(restErr.StatusCode(), restErr)
		return
	}

	c.String(http.StatusOK, "mapping is saved.")
}
