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

var (
	router = gin.Default()
	server = getServer()
)

func init() {
	config.Init()
	elasticsearch.Init()
	mongodb.Init()
	nats.Init()
}

func StartApplication() {
	mapUrls()
	go listen()
	waitForShutdown()
	logger.Info("the server exited successfully")
}

func getServer() *http.Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Conf.Server.Port),
		Handler:      router,
		ReadTimeout:  config.Conf.Server.TimeoutRead,
		WriteTimeout: config.Conf.Server.TimeoutWrite,
		IdleTimeout:  config.Conf.Server.TimeoutIdle,
	}
	return srv
}

func listen() {
	listenToEvents()
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Panic("error when trying to start the app", err)
	}
}
