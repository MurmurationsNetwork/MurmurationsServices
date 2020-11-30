package app

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
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
	elasticsearch.Init()
	mongodb.Init()
	nats.Init()
}

func StartApplication() {
	mapUrls()
	go listen(server)
	listenToEvents()

	waitForShutdown()
	logger.Info("the server exited successfully")
}

func getServer() *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	return srv
}

func listen(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Panic("error when trying to start the app", err)
	}
}
