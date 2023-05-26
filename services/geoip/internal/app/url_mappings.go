package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/internal/controller/http"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/internal/service"
	"github.com/gin-gonic/gin"
)

func mapURLs(router *gin.Engine) {
	gepIPHandler := http.NewGepIPHandler(service.NewGeoIPService())
	router.GET("/city/:ip", gepIPHandler.GetCity)

	pingHandler := http.NewPingHandler()
	router.GET("/ping", pingHandler.Ping)
}
