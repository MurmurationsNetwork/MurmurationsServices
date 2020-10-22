package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
)

func waitForShutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cleanup()

	logger.Info("trying to shut down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", err)
	}
}
