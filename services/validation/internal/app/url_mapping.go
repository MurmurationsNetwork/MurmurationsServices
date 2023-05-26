package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/adapter/controller/http"
	"github.com/gin-gonic/gin"
)

func mapURLs(router *gin.Engine) {
	pingHandler := http.NewPingHandler()
	router.GET("/ping", pingHandler.Ping)
}
