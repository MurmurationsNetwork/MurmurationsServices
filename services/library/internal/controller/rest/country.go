package rest

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/config"
)

// CountryHandler defines the methods that a country handler should implement.
type CountryHandler interface {
	GetMap(c *gin.Context)
}

type countryHandler struct {
}

// NewCountryHandler creates a new country handler.
func NewCountryHandler() CountryHandler {
	return &countryHandler{}
}

// GetMap handles the request to get the country map.
func (handler *countryHandler) GetMap(c *gin.Context) {
	contents, err := os.ReadFile(
		config.Values.Static.StaticFilePath + "/countries.json",
	)

	if err != nil {
		errors := jsonapi.NewError(
			[]string{"Get countries map error"},
			[]string{"error:" + err.Error()},
			nil,
			[]int{http.StatusInternalServerError},
		)
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	c.Header("Content-Type", "application/json")

	c.String(http.StatusOK, string(contents))
}
