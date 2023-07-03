package app

import (
	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/controller/http"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/service"
)

func mapURLs(router *gin.Engine) {
	pingHandler := http.NewPingHandler()
	mappingsHandler := http.NewMappingsHandler(mongo.NewMappingRepository())
	profilesHandler := http.NewProfilesHandler(mongo.NewProfileRepository())
	updatesHandler := http.NewUpdatesHandler(mongo.NewUpdateRepository())
	batchesHandler := http.NewBatchesHandler(
		service.NewBatchService(mongo.NewBatchRepository()),
	)

	v1 := router.Group("/v1")
	{
		v1.GET("/ping", pingHandler.Ping)
		v1.POST("/mappings", mappingsHandler.Create)
		v1.GET("/profiles/:profileID", profilesHandler.Get)
		v1.GET("/health/:schemaName", updatesHandler.Get)

		// for csv batch import
		v1.GET("/batch/user", batchesHandler.GetBatchesByUserID)
		v1.POST("/batch/validate", batchesHandler.Validate)
		v1.POST("/batch/import", batchesHandler.Import)
		v1.PUT("/batch/import", batchesHandler.Edit)
		v1.DELETE("/batch/import", batchesHandler.Delete)
	}
}
