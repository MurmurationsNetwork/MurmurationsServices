package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeprecationHandler interface {
	DeprecationV1(c *gin.Context)
}

type deprecationHandler struct {
}

func NewDeprecationHandler() DeprecationHandler {
	return &deprecationHandler{}
}

func (handler *deprecationHandler) DeprecationV1(c *gin.Context) {
	c.JSON(http.StatusGone, gin.H{"message": "The v1 API has been deprecated. Please use the v2 API instead: https://app.swaggerhub.com/apis-docs/MurmurationsNetwork/IndexAPI/2.0.0"})
}
