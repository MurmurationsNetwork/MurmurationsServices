package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	PingController pingControllerInterface = &pingController{}
)

type pingControllerInterface interface {
	Ping(c *gin.Context)
}

type pingController struct{}

func (cont *pingController) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
