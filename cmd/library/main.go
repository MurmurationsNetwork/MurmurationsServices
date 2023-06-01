// This is a Go program that defines the main function for the library service.
// The import statements import two packages: logger and library. The logger
// package is used for logging messages, while the library package contains the
// implementation of the library service.

package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/pkg/library"
)

func main() {
	// Start by logging a message to indicate that the service is starting up.
	logger.Info("Library service starting")

	// Create a new instance of the library.Service struct using the
	// library.NewService function. This function initializes the service and
	// returns a pointer to the Service struct.
	s := library.NewService()

	// Start a new goroutine that waits for the service to start up using the
	// WaitUntilUp method of the Service struct. This method returns a channel
	// that is closed when the service is up and running. The goroutine waits
	// for the channel to be closed and then logs a message to indicate that the
	// service has started.
	go func() {
		<-s.WaitUntilUp()
		logger.Info("Library service started")
	}()

	// Call the Run method of the Service struct to start the service. This
	// method blocks until the service is shut down.
	s.Run()
}
