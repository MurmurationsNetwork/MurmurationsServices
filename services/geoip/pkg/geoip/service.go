package geoip

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tevino/abool"
	"go.uber.org/zap/zapcore"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/core"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/internal/controller/rest"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/internal/service"
)

func init() {
	global.Init()
}

// Service represents the geoip service.
type Service struct {
	// HTTP server
	server *http.Server
	// Atomic boolean to manage service state
	run *abool.AtomicBool
	// HTTP router for the geoip service
	router *gin.Engine
	// Ensures cleanup is only run once
	runCleanup sync.Once
	// Context for shutdown
	shutdownCtx context.Context
	// Cancel function for shutdown context
	shutdownCancelCtx context.CancelFunc
}

// NewService initializes a new geoip service.
func NewService() *Service {
	svc := &Service{
		run: abool.New(),
	}

	svc.setupServer()
	core.InstallShutdownHandler(svc.Shutdown)

	return svc
}

// setupServer configures and initializes the HTTP server.
func (s *Service) setupServer() {
	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.router.Use(s.middlewares()...)
	s.registerRoutes()

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Conf.Server.Port),
		Handler:      s.router,
		ReadTimeout:  config.Conf.Server.TimeoutRead,
		WriteTimeout: config.Conf.Server.TimeoutWrite,
		IdleTimeout:  config.Conf.Server.TimeoutIdle,
	}

	s.shutdownCtx, s.shutdownCancelCtx = context.WithCancel(
		context.Background(),
	)
}

func (s *Service) middlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Recovery(),
	}
}

// registerRoutes sets up the routes for the HTTP server.
func (s *Service) registerRoutes() {
	geoIPHandler := rest.NewGeoIPHandler(service.NewGeoIPService())
	s.router.GET("/city/:ip", geoIPHandler.GetCity)

	pingHandler := rest.NewPingHandler()
	s.router.GET("/ping", pingHandler.Ping)
}

// panic performs a cleanup and then emits the supplied message as the panic value.
func (s *Service) panic(msg string, err error, logFields ...zapcore.Field) {
	s.cleanup()
	logger.Panic(msg, err, logFields...)
}

// Run starts the geoip service and will block until the service is shutdown.
func (s *Service) Run() {
	s.run.Set()
	if err := s.server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		s.panic("Error when trying to start the server", err)
	}
}

// WaitUntilUp returns a channel which blocks until the geoip service is up.
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
			logger.Info(
				"Ping failed, waiting for service to finish starting...",
			)
			time.Sleep(time.Second)
		}
	}()
	return initialized
}

// Shutdown stops the geoip service.
func (s *Service) Shutdown() {
	if s.run.IsSet() {
		if err := s.server.Shutdown(s.shutdownCtx); err != nil {
			logger.Error("GeoIP service shutdown failure", err)
		}
	}
	s.cleanup()
}

// cleanup will clean up the non-server resources associated with the service.
func (s *Service) cleanup() {
	s.runCleanup.Do(func() {
		s.shutdownCancelCtx()
		global.DB.Close()
		logger.Info("GeoIP service stopped gracefully")
	})
}
