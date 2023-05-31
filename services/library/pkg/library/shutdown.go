package library

import (
	"os"
	"os/signal"
	"syscall"
)

// InstallShutdownHandler sets up a signal listener and calls the provided
// shutdown function when an interrupt or terminate signal is received.
func InstallShutdownHandler(shutdown func()) {
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		shutdown()
	}()
}
