package http

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type CountryHandler interface {
	GetMap(c *gin.Context)
}

type countryHandler struct {
}

func NewCountryHandler() CountryHandler {
	return &countryHandler{}
}

func (handler *countryHandler) GetMap(c *gin.Context) {
	contents, err := os.ReadFile("static/countries/map.json")
	if err != nil {
		errors := jsonapi.NewError([]string{"Get countries map error"}, []string{"error:" + err.Error()}, nil, []int{http.StatusInternalServerError})
		res := jsonapi.Response(nil, errors, nil, nil)
		c.JSON(errors[0].Status, res)
		return
	}

	c.Header("Content-Type", "application/json")

	c.String(http.StatusOK, string(contents))
}
