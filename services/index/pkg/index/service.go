package index

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tevino/abool/v2"
	"go.uber.org/zap/zapcore"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/core"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/handler"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/middleware/limiter"
	midlogger "github.com/MurmurationsNetwork/MurmurationsServices/pkg/middleware/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/natsclient"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/controller/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/controller/rest"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/repository/es"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/service"
)

// Service represents the index service.
type Service struct {
	// HTTP server
	server *http.Server
	// Node event handler
	nodeHandler event.NodeHandler
	// Atomic boolean to manage service state
	run *abool.AtomicBool
	// HTTP router for the index service
	router *gin.Engine
	// Ensures cleanup is only run once
	runCleanup sync.Once
	// Context for shutdown
	shutdownCtx context.Context
	// Cancel function for shutdown context
	shutdownCancelCtx context.CancelFunc
}

// NewService initializes a new index service.
func NewService() *Service {
	svc := &Service{
		run: abool.New(),
	}

	svc.setupNATS()

	svc.setupServer()
	svc.nodeHandler = event.NewNodeHandler(
		service.NewNodeService(
			mongo.NewNodeRepository(),
			es.NewNodeRepository(),
		),
	)
	core.InstallShutdownHandler(svc.Shutdown)

	return svc
}

// setupNATS initializes Nats service.
func (s *Service) setupNATS() {
	err := natsclient.Initialize(config.Values.Nats.URL)
	if err != nil {
		logger.Error("Failed to create Nats client", err)
		os.Exit(1)
	}
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
	return nil
}

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
		// CORS for all origins, allowing:
		// - GET, POST and DELETE methods
		// - Origin, Authorization and Content-Type header
		// - Credentials share
		// - Preflight requests cached for 12 hours
		cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "DELETE"},
			AllowHeaders: []string{
				"Origin",
				"Authorization",
				"Content-Type",
			},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	}
}

// registerRoutes sets up the routes for the HTTP server.
func (s *Service) registerRoutes() {
	nodeHandler := rest.NewNodeHandler(
		service.NewNodeService(
			mongo.NewNodeRepository(),
			es.NewNodeRepository(),
		),
	)

	s.setupV1Routes()
	s.setupV2Routes(nodeHandler)
}

// setupV1Routes configures routes for API version 1.
func (s *Service) setupV1Routes() {
	v1 := s.router.Group("/v1")
	v1.Any("/*any", handler.NewDeprecationHandler("Index"))
}

// setupV2Routes configures routes for API version 2.
func (s *Service) setupV2Routes(nodeHandler rest.NodeHandler) {
	v2 := s.router.Group("/v2")
	v2.GET("/ping", handler.PingHandler)
	v2.PUT(
		"/feature-flag/:flagName",
		AllowInNonProductionMiddleware(),
		s.ToggleFeatureFlag,
	)

	// Node-related routes
	v2.POST("/nodes", nodeHandler.Add)
	v2.GET("/nodes/:nodeID", nodeHandler.Get)
	v2.GET("/nodes", nodeHandler.Search)
	v2.DELETE("/nodes", nodeHandler.Delete)
	v2.DELETE("/nodes/:nodeID", nodeHandler.Delete)
	v2.POST("/validate", nodeHandler.Validate)
	v2.POST("/nodes-sync", nodeHandler.AddSync)
	v2.POST("/export", nodeHandler.Export)
	v2.GET("/get-nodes", nodeHandler.GetNodes)
}

// panic performs a cleanup and then emits the supplied message as the panic value.
func (s *Service) panic(msg string, err error, logFields ...zapcore.Field) {
	s.cleanup()
	logger.Error(msg, err, logFields...)
	os.Exit(1)
}

// Run starts the index service and will block until the service is shutdown.
func (s *Service) Run() {
	s.run.Set()
	if err := s.nodeHandler.Validated(); err != nil &&
		err != http.ErrServerClosed {
		s.panic("Error when trying to listen events", err)
	}
	if err := s.nodeHandler.ValidationFailed(); err != nil &&
		err != http.ErrServerClosed {
		s.panic("Error when trying to listen events", err)
	}
	if err := s.server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		s.panic("Error when trying to start the server", err)
	}
}

// WaitUntilUp returns a channel which blocks until the index service is up.
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
						"http://localhost:%s/v2/ping",
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

// Shutdown stops the index service.
func (s *Service) Shutdown() {
	if s.run.IsSet() {
		if err := s.server.Shutdown(s.shutdownCtx); err != nil {
			logger.Error("Index service shutdown failure", err)
		}
	}
	s.cleanup()
}

// cleanup will clean up the non-server resources associated with the service.
func (s *Service) cleanup() {
	s.runCleanup.Do(func() {
		var errOccurred bool

		// Shutdown the context.
		s.shutdownCancelCtx()

		// Disconnect from MongoDB.
		mongodb.Client.Disconnect()

		// Disconnect from NATS.
		if err := natsclient.GetInstance().Disconnect(); err != nil {
			logger.Error("Error disconnecting from NATS: %v", err)
			errOccurred = true
		}

		// Log based on whether an error occurred.
		if errOccurred {
			logger.Info("Index service stopped with errors.")
		} else {
			logger.Info("Index service stopped gracefully.")
		}
	})
}

// ToggleFeatureFlag updates the state of a given feature flag.
func (s *Service) ToggleFeatureFlag(c *gin.Context) {
	flagName := c.Param("flagName")
	enable := c.Query("enable")

	// Convert newState to boolean.
	isEnabled, err := strconv.ParseBool(enable)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state value"})
		return
	}

	// Capture the previous state of the feature flag.
	prevState, flagExists := config.Values.FeatureToggles[flagName]

	// Update the feature flag state.
	config.Values.FeatureToggles[flagName] = isEnabled

	// Log the change for audit purposes.
	log.Printf(
		"Feature flag '%s' changed from '%t' to '%t'",
		flagName,
		prevState,
		isEnabled,
	)

	// Constructing a detailed response.
	response := gin.H{
		"flagName":      flagName,
		"previousState": prevState,
		"currentState":  isEnabled,
		"flagExists":    flagExists,
		"timestamp":     time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}
