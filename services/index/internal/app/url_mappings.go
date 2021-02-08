package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/controller/http"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/usecase"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine) {
	nodeHandler := http.NewNodeHandler(usecase.NewNodeService(db.NewRepository()))
	router.POST("/nodes", nodeHandler.Add)
	router.GET("/nodes/:nodeId", nodeHandler.Get)
	router.GET("/nodes", nodeHandler.Search)
	router.DELETE("/nodes/:nodeId", nodeHandler.Delete)

	pingHandler := http.NewPingHandler()
	router.GET("/ping", pingHandler.Ping)
}
