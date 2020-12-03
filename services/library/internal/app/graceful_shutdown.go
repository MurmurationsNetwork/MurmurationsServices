package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/config"
)

func waitForShutdown(server *http.Server, closed chan struct{}) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down library service")

	cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Server.TimeoutIdle)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Library service shutdown failure", err)
	}

	close(closed)
}
