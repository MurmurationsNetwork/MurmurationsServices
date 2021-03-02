package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler interface {
	Ping(c *gin.Context)
}

type pingHandler struct {
}

func NewPingHandler() PingHandler {
	return &pingHandler{}
}

func (handler *pingHandler) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
