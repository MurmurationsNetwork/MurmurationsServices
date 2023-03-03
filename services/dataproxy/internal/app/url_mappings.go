package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/controller/http"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/usecase"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine) {
	pingHandler := http.NewPingHandler()
	mappingsHandler := http.NewMappingsHandler(db.NewMappingRepository())
	profilesHandler := http.NewProfilesHandler(db.NewProfileRepository())
	updatesHandler := http.NewUpdatesHandler(db.NewUpdateRepository())
	batchesHandler := http.NewBatchesHandler(usecase.NewBatchService(db.NewBatchRepository()))

	v1 := router.Group("/v1")
	{
		v1.GET("/ping", pingHandler.Ping)
		v1.POST("/mappings", mappingsHandler.Create)
		v1.GET("/profiles/:profileId", profilesHandler.Get)
		v1.GET("/health/:schemaName", updatesHandler.Get)

		// for csv batch import
		v1.GET("/batch/user", batchesHandler.GetBatchesByUserID)
		v1.POST("/batch/validate", batchesHandler.Validate)
		v1.POST("/batch/import", batchesHandler.Import)
		v1.PUT("/batch/import", batchesHandler.Edit)
		v1.DELETE("/batch/import", batchesHandler.Delete)
	}
}
