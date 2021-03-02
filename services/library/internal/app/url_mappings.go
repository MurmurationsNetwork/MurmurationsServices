package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/controller/http"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/service"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine) {
	schemaHandler := http.NewSchemaHandler(service.NewSchemaService(db.NewSchemaRepo()))
	router.GET("/schemas", schemaHandler.Search)

	pingHandler := http.NewPingHandler()
	router.GET("/ping", pingHandler.Ping)
}
