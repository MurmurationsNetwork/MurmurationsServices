package app

import (
	"fmt"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/adapter/elasticsearch"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/adapter/mongodb"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/adapter/nats"
	"github.com/gin-gonic/gin"
)

func init() {
	config.Init()
	elasticsearch.Init()
	mongodb.Init()
	nats.Init()
}

func StartApplication() {
	router := gin.Default()
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
