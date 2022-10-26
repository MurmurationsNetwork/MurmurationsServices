package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/controller/http"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/service"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine) {
	deprecationHandler := http.NewDeprecationHandler()
	pingHandler := http.NewPingHandler()
	schemaHandler := http.NewSchemaHandler(service.NewSchemaService(db.NewSchemaRepo()))

	v1 := router.Group("/v1")
	{
		v1.Any("/*any", deprecationHandler.DeprecationV1)
	}

	v2 := router.Group("/v2")
	{
		v2.GET("/ping", pingHandler.Ping)
		v2.GET("/schemas", schemaHandler.Search)
		v2.GET("/schemas/:schemaName", schemaHandler.Get)
	}
}
