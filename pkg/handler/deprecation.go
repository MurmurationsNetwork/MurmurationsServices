package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
)

// DeprecationHandler handles requests made to deprecated API versions.
func DeprecationHandler(c *gin.Context) {
	errorMessage := "The v1 API has been deprecated. " +
		"Please use the v2 API instead: " +
		"https://app.swaggerhub.com/apis-docs/MurmurationsNetwork/LibraryAPI/2.0.0"

	errors := jsonapi.NewError(
		[]string{"Gone"},
		[]string{errorMessage},
		nil,
		[]int{http.StatusGone},
	)

	res := jsonapi.Response(nil, errors, nil, nil)
	c.JSON(errors[0].Status, res)
}
