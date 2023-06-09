package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler defines the interface for a ping handler.
type PingHandler interface {
	Ping(c *gin.Context)
}

// pingHandler implements the PingHandler interface.
type pingHandler struct{}

// NewPingHandler creates a new ping handler.
func NewPingHandler() PingHandler {
	return &pingHandler{}
}

// Ping responds with a "pong" message.
func (handler *pingHandler) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
