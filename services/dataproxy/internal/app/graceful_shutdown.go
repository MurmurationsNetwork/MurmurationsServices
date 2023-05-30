package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
)

func waitForShutdown(server *http.Server, closed chan struct{}) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cleanup()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		config.Conf.Server.TimeoutIdle,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Data-Proxy service shutdown failure", err)
	}

	close(closed)
}
