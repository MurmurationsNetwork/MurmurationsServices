package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/adapter/controller/http"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/adapter/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/usecase"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine) {
	deprecationHandler := http.NewDeprecationHandler()
	pingHandler := http.NewPingHandler()
	nodeHandler := http.NewNodeHandler(usecase.NewNodeService(db.NewRepository()))

	v1 := router.Group("/v1")
	{
		v1.Any("/*any", deprecationHandler.DeprecationV1)
	}

	v2 := router.Group("/v2")
	{
		v2.GET("/ping", pingHandler.Ping)

		v2.POST("/nodes", nodeHandler.Add)
		v2.GET("/nodes/:nodeId", nodeHandler.Get)
		v2.GET("/nodes", nodeHandler.Search)
		v2.DELETE("/nodes/:nodeId", nodeHandler.Delete)
		v2.POST("/validate", nodeHandler.Validate)

		// synchronously response
		v2.POST("/nodes-sync", nodeHandler.AddSync)
	}
}
