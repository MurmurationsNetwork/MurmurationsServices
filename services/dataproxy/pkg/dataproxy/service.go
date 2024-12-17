package dataproxy

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	corslib "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tevino/abool"
	"go.uber.org/zap/zapcore"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/core"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/middleware/limiter"
	midlogger "github.com/MurmurationsNetwork/MurmurationsServices/pkg/middleware/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/controller/rest"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/service"
)

// Service represents the dataproxy service.
type Service struct {
	// HTTP server
	server *http.Server
	// Atomic boolean to manage service state
	run *abool.AtomicBool
	// HTTP router for the dataproxy service
	router *gin.Engine
	// Ensures cleanup is only run once
	runCleanup sync.Once
	// Context for shutdown
	shutdownCtx context.Context
	// Cancel function for shutdown context
	shutdownCancelCtx context.CancelFunc
}

// NewService initializes a new dataproxy service.
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

	if err := s.connectToMongoDB(); err != nil {
		s.panic("error when trying to connect to MongoDB", err)
	}

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Values.Server.Port),
		Handler:      s.router,
		ReadTimeout:  config.Values.Server.TimeoutRead,
		WriteTimeout: config.Values.Server.TimeoutWrite,
		IdleTimeout:  config.Values.Server.TimeoutIdle,
	}

	s.shutdownCtx, s.shutdownCancelCtx = context.WithCancel(
		context.Background(),
	)
}

// connectToMongoDB establishes a connection to MongoDB.
func (s *Service) connectToMongoDB() error {
	uri := mongodb.GetURI(
		config.Values.Mongo.USERNAME,
		config.Values.Mongo.PASSWORD,
		config.Values.Mongo.HOST,
	)
	err := mongodb.NewClient(uri, config.Values.Mongo.DBName)
	if err != nil {
		return err
	}
	err = mongodb.Client.Ping()
	if err != nil {
		return err
	}
	// Create an index on the `cuid` field in the `profiles` collection
	err = mongodb.Client.CreateUniqueIndex("profiles", "cuid")
	if err != nil {
		logger.Error("error when trying to create index on `cuid` field in `profiles` collection", err)
		os.Exit(1)
	}
	return nil
}

// middlewares returns the list of middlewares to be used by the HTTP server.
func (s *Service) middlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Recovery(),
		limiter.NewRateLimitWithOptions(limiter.RateLimitOptions{
			Period: config.Values.Server.PostRateLimitPeriod,
			Method: "POST",
		}),
		limiter.NewRateLimitWithOptions(limiter.RateLimitOptions{
			Period: config.Values.Server.GetRateLimitPeriod,
			Method: "GET",
		}),
		midlogger.NewLogger(),
		s.cors(),
	}
}

// cors returns the CORS middleware configuration.
func (s *Service) cors() gin.HandlerFunc {
	// CORS for all origins, allowing:
	// - GET and POST methods
	// - Origin, Authorization and Content-Type header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	return corslib.New(corslib.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// registerRoutes sets up the routes for the HTTP server.
func (s *Service) registerRoutes() {
	pingHandler := rest.NewPingHandler()
	mappingsHandler := rest.NewMappingsHandler(mongo.NewMappingRepository())
	profilesHandler := rest.NewProfilesHandler(mongo.NewProfileRepository())
	updatesHandler := rest.NewUpdatesHandler(mongo.NewUpdateRepository())
	batchesHandler := rest.NewBatchesHandler(
		service.NewBatchService(mongo.NewBatchRepository()),
	)

	v1 := s.router.Group("/v1")
	{
		v1.GET("/ping", pingHandler.Ping)
		v1.POST("/mappings", mappingsHandler.Create)
		v1.GET("/profiles/:profileID", profilesHandler.Get)
		v1.GET("/health/:schemaName", updatesHandler.Get)

		// for csv batch import
		v1.GET("/batch/user", batchesHandler.GetBatchesByUserID)
		v1.POST("/batch/validate", batchesHandler.Validate)
		v1.POST("/batch/import", batchesHandler.Import)
		v1.PUT("/batch/import", batchesHandler.Edit)
		v1.DELETE("/batch/import", batchesHandler.Delete)
	}
}

// panic performs a cleanup and then emits the supplied message as the panic value.
func (s *Service) panic(msg string, err error, logFields ...zapcore.Field) {
	s.cleanup()
	logger.Error(msg, err, logFields...)
	os.Exit(1)
}

// Run starts the dataproxy service and will block until the service is shutdown.
func (s *Service) Run() {
	s.run.Set()

	if err := s.server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		s.panic("Error when trying to start the server", err)
	}
}

// WaitUntilUp returns a channel which blocks until the dataproxy service is up.
func (s *Service) WaitUntilUp() <-chan struct{} {
	initialized := make(chan struct{})
	timeout := time.After(1 * time.Minute)
	go func() {
		for {
			select {
			case <-timeout:
				panic("Service startup timed out")
			default:
				resp, err := http.Get(
					fmt.Sprintf(
						"http://localhost:%s/v1/ping",
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
				time.Sleep(5 * time.Second)
			}
		}
	}()
	return initialized
}

// Shutdown stops the dataproxy service.
func (s *Service) Shutdown() {
	if s.run.IsSet() {
		if err := s.server.Shutdown(s.shutdownCtx); err != nil {
			logger.Error("Data-Proxy service shutdown failure", err)
		}
	}
	s.cleanup()
}

// cleanup will clean up the non-server resources associated with the service.
func (s *Service) cleanup() {
	s.runCleanup.Do(func() {
		// Shutdown the context.
		s.shutdownCancelCtx()
	})
}
