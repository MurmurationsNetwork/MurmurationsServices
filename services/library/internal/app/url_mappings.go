package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/controller/http"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/service"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine) {
	pingHandler := http.NewPingHandler()
	schemaHandler := http.NewSchemaHandler(service.NewSchemaService(db.NewSchemaRepo()))

	v1 := router.Group("/v1")
	{
		v1.GET("/ping", pingHandler.Ping)
		v1.GET("/schemas", schemaHandler.Search)
	}
}
