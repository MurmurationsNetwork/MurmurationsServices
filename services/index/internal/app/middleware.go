package app

import (
	"time"

	"github.com/gin-gonic/gin"

	corslib "github.com/gin-contrib/cors"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/middleware/limiter"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/middleware/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

func getMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Recovery(),
		limiter.NewRateLimitWithOptions(limiter.RateLimitOptions{
			Period: config.Conf.Server.PostRateLimitPeriod,
			Method: "POST",
		}),
		limiter.NewRateLimitWithOptions(limiter.RateLimitOptions{
			Period: config.Conf.Server.GetRateLimitPeriod,
			Method: "GET",
		}),
		logger.NewLogger(),
		cors(),
	}
}

func cors() gin.HandlerFunc {
	// CORS for all origins, allowing:
	// - GET, POST and DELETE methods
	// - Origin, Authorization and Content-Type header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	return corslib.New(corslib.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
