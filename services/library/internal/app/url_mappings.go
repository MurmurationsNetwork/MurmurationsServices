package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/controller/http"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine) {
	router.GET("/ping", http.PingController.Ping)
}
