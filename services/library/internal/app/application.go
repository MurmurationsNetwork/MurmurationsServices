package app

import (
	"fmt"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/middleware"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/global"
	"github.com/gin-gonic/gin"
)

func init() {
	global.Init()
}

func StartApplication() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), middleware.RateLimit(config.Conf.Server.RateLimitPeriod), middleware.Logger(), CORS())

	mapUrls(router)

	server := getServer(router)

	closed := make(chan struct{})
	go waitForShutdown(server, closed)

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
