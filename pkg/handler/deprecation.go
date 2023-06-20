package handler

import (
	"fmt"
	"net/http"
	"unicode"

	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
)

// NewDeprecationHandler handles requests made to deprecated API versions.
func NewDeprecationHandler(service string) gin.HandlerFunc {
	service = ensureFirstUpper(service)

	return func(c *gin.Context) {
		errorMessage := fmt.Sprintf(
			"The v1 API has been deprecated. "+
				"Please use the v2 API instead: "+
				"https://app.swaggerhub.com/apis-docs/MurmurationsNetwork/%sAPI/2.0.0",
			service,
		)

		apiErrors := jsonapi.NewError(
			[]string{"Gone"},
			[]string{errorMessage},
			nil,
			[]int{http.StatusGone},
		)

		response := jsonapi.Response(nil, apiErrors, nil, nil)
		c.JSON(http.StatusGone, response)
	}
}

// ensureFirstUpper ensures that the first letter of a string is uppercase and
// the rest are lowercase.
func ensureFirstUpper(input string) string {
	r := []rune(input)
	r[0] = unicode.ToUpper(r[0])
	for i := 1; i < len(r); i++ {
		r[i] = unicode.ToLower(r[i])
	}
	return string(r)
}
