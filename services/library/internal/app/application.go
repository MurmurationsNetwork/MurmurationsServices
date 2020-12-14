package app

import (
	"fmt"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/adapter/mongodb"
	"github.com/gin-gonic/gin"
)

func init() {
	config.Init()
	mongodb.Init()
}

func StartApplication() {
	router := gin.Default()

	router.Use(CORS())

	mapUrls(router)

	server := getServer(router)

	closed := make(chan struct{})
	go waitForShutdown(server, closed)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Panic("Error when trying to start the server", err)
	}

	<-closed
	logger.Info("The service exited successfully")
}

func getServer(router *gin.Engine) *http.Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Conf.Server.Port),
		Handler:      router,
		ReadTimeout:  config.Conf.Server.TimeoutRead,
		WriteTimeout: config.Conf.Server.TimeoutWrite,
		IdleTimeout:  config.Conf.Server.TimeoutIdle,
	}
	return srv
}
