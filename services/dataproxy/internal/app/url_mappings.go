package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/controller/http"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine) {
	pingHandler := http.NewPingHandler()

	v1 := router.Group("/v1")
	{
		v1.GET("/ping", pingHandler.Ping)
	}
}
