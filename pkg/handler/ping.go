package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler responds to ping requests with "pong!".
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong!")
}
