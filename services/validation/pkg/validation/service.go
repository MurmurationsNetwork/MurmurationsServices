package validation

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tevino/abool"
	"go.uber.org/zap/zapcore"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	midlogger "github.com/MurmurationsNetwork/MurmurationsServices/pkg/middleware/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/adapter/controller/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/adapter/controller/rest"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/service"
)

func init() {
	global.Init()
}

// Service represents the validation service.
type Service struct {
	// HTTP server
	server *http.Server
	// Atomic boolean to manage service state
	run *abool.AtomicBool
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
		run: abool.New(),
	}

	svc.setupServer()
	svc.nodeHandler = event.NewNodeHandler(service.NewValidationService())
	InstallShutdownHandler(svc.Shutdown)

	return svc
}

// setupServer configures and initializes the HTTP server.
func (s *Service) setupServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), midlogger.NewLogger())

	pingHandler := rest.NewPingHandler()
	router.GET("/ping", pingHandler.Ping)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Conf.Server.Port),
		Handler:      router,
		ReadTimeout:  config.Conf.Server.TimeoutRead,
		WriteTimeout: config.Conf.Server.TimeoutWrite,
		IdleTimeout:  config.Conf.Server.TimeoutIdle,
	}

	s.shutdownCtx, s.shutdownCancelCtx = context.WithCancel(
		context.Background(),
	)
}

// panic performs a cleanup and then emit the supplied message as the panic value.
func (s *Service) panic(msg string, err error, logFields ...zapcore.Field) {
	s.cleanup()
	logger.Panic(msg, err, logFields...)
}

// Run starts the validation service and will block until the service is shutdown.
func (s *Service) Run() {
	s.run.Set()
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
					config.Conf.Server.Port,
				),
			)
			if err == nil && resp.StatusCode == http.StatusOK {
				close(initialized)
				return
			}
			logger.Info("Ping failed, waiting for service to go up...")
			time.Sleep(time.Second)
		}
	}()
	return initialized
}

// Shutdown stops the validation service.
func (s *Service) Shutdown() {
	if s.run.IsSet() {
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
		logger.Info("validation service stopped gracefully")
	})
}
