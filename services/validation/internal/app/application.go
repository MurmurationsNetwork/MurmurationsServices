package app

import (
	"fmt"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/middleware"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/global"
	"github.com/gin-gonic/gin"
)

func init() {
	global.Init()
}

func StartApplication() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), middleware.Logger())

	mapUrls(router)

	server := getServer(router)

	closed := make(chan struct{})
	go waitForShutdown(server, closed)

	if err := listenToEvents(); err != nil && err != http.ErrServerClosed {
		logger.Panic("Error when trying to listen events", err)
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Panic("Error when trying to start the server", err)
	}

	<-closed
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
