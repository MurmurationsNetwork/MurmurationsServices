package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/controller/http"
)

func mapUrls() {
	router.GET("/ping", http.PingController.Ping)

	router.POST("/nodes", http.NodeController.Add)
	router.GET("/nodes/:nodeId", http.NodeController.Get)
	router.GET("/nodes", http.NodeController.Search)
	router.DELETE("/nodes/:nodeId", http.NodeController.Delete)
}
