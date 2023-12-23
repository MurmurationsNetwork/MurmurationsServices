package validation

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tevino/abool/v2"
	"go.uber.org/zap/zapcore"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/core"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/handler"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	midlogger "github.com/MurmurationsNetwork/MurmurationsServices/pkg/middleware/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/controller/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/service"
)

// Service represents the validation service.
type Service struct {
	// HTTP server
	server *http.Server
	// Atomic boolean to manage service state
	isRunning *abool.AtomicBool
	// Node event handler
	nodeHandler event.NodeHandler
	// Ensures cleanup is only run once
	runCleanup sync.Once
	// Context for shutdown
	shutdownCtx context.Context
	// Cancel function for shutdown context
	shutdownCancelCtx context.CancelFunc
}

// NewService initializes a new validation service.
func NewService() *Service {
	svc := &Service{
		isRunning: abool.New(),
	}

	svc.setupServer()
	svc.nodeHandler = event.NewNodeHandler(service.NewValidationService())
	core.InstallShutdownHandler(svc.Shutdown)

	return svc
}

// setupServer configures and initializes the HTTP server.
func (s *Service) setupServer() {
	err := nats.NewClient(config.Values.NATS.URL)
	if err != nil {
		logger.Panic("failed to connect to NATS", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), midlogger.NewLogger())

	router.GET("/ping", handler.PingHandler)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Values.Server.Port),
		Handler:      router,
		ReadTimeout:  config.Values.Server.TimeoutRead,
		WriteTimeout: config.Values.Server.TimeoutWrite,
		IdleTimeout:  config.Values.Server.TimeoutIdle,
	}

	s.shutdownCtx, s.shutdownCancelCtx = context.WithCancel(
		context.Background(),
	)
}

// panic performs a cleanup and then emits the supplied message as the panic value.
func (s *Service) panic(msg string, err error, logFields ...zapcore.Field) {
	s.cleanup()
	logger.Panic(msg, err, logFields...)
}

// Run starts the validation service and will block until the service is shutdown.
func (s *Service) Run() {
	s.isRunning.Set()
	if err := s.nodeHandler.NewNodeCreatedListener(); err != nil &&
		err != http.ErrServerClosed {
		s.panic("Error when trying to listen events", err)
	}
	if err := s.server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		s.panic("Error when trying to start the server", err)
	}
}

// WaitUntilUp returns a channel which blocks until the validation service is up.
func (s *Service) WaitUntilUp() <-chan struct{} {
	initialized := make(chan struct{})
	go func() {
		for {
			resp, err := http.Get(
				fmt.Sprintf(
					"http://localhost:%s/ping",
					config.Values.Server.Port,
				),
			)
			if err == nil && resp.StatusCode == http.StatusOK {
				close(initialized)
				return
			}
			logger.Info(
				"Ping failed, waiting for service to finish starting...",
			)
			time.Sleep(time.Second)
		}
	}()
	return initialized
}

// Shutdown stops the validation service.
func (s *Service) Shutdown() {
	if s.isRunning.IsSet() {
		if err := s.server.Shutdown(s.shutdownCtx); err != nil {
			logger.Error("Validation service shutdown failure", err)
		}
	}
	s.cleanup()
}

// cleanup will clean up the non-server resources associated with the service.
func (s *Service) cleanup() {
	s.runCleanup.Do(func() {
		s.shutdownCancelCtx()
		nats.Client.Disconnect()
		logger.Info("Validation service stopped gracefully")
	})
}
