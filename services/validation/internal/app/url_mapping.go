package app

import (
	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/adapter/controller/http"
)

func mapURLs(router *gin.Engine) {
	pingHandler := http.NewPingHandler()
	router.GET("/ping", pingHandler.Ping)
}
