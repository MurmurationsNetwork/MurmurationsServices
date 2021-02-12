package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/adapter/controller/http"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine) {
	pingHandler := http.NewPingHandler()
	router.GET("/ping", pingHandler.Ping)
}
