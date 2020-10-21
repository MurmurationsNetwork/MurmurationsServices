package app

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
	server = getServer()
)

func StartApplication() {
	mapUrls()
	go listen(server)

	waitForShutdown()
	cleanup()
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
