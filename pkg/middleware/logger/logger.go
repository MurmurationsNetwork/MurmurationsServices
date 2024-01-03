package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
)

// defaultSkipPaths defines paths that should be skipped in logging.
var defaultSkipPaths = []string{"/ping", "/v1/ping", "/v2/ping"}

// NewLogger creates a Gin middleware logger with default configuration.
func NewLogger() gin.HandlerFunc {
	return NewLoggerWithConfig(gin.LoggerConfig{})
}

// NewLoggerWithConfig creates a Gin middleware logger with custom configuration.
func NewLoggerWithConfig(conf gin.LoggerConfig) gin.HandlerFunc {
	skipPaths := mergeSkipPaths(conf.SkipPaths, defaultSkipPaths)
	skipMap := createSkipMap(skipPaths)

	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path

		c.Next()

		if _, shouldSkip := skipMap[path]; !shouldSkip {
			logRequest(c, startTime, path)
		}
	}
}

// mergeSkipPaths combines custom skip paths with default skip paths.
func mergeSkipPaths(userPaths, defaultPaths []string) []string {
	return append(userPaths, defaultPaths...)
}

// createSkipMap generates a map of paths to be skipped.
func createSkipMap(paths []string) map[string]struct{} {
	skipMap := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		skipMap[path] = struct{}{}
	}
	return skipMap
}

// logRequest logs the details of each HTTP request.
func logRequest(c *gin.Context, startTime time.Time, path string) {
	param := gin.LogFormatterParams{
		Request:    c.Request,
		Keys:       c.Keys,
		TimeStamp:  time.Now(),
		Latency:    time.Since(startTime),
		ClientIP:   c.ClientIP(),
		Method:     c.Request.Method,
		StatusCode: c.Writer.Status(),
		BodySize:   c.Writer.Size(),
		Path:       path,
	}

	geoInfo := getGeoInfo(param.ClientIP)

	logger.Info(
		"Log Entry",
		zap.Int("status", param.StatusCode),
		zap.String("latency", fmt.Sprintf("%v", param.Latency)),
		zap.String("method", param.Method),
		zap.String("path", param.Path),
		zap.String("ip", param.ClientIP),
		zap.String("city", geoInfo.City),
		zap.String("country", geoInfo.Country),
		zap.Float64("lat", geoInfo.Lat),
		zap.Float64("lon", geoInfo.Lon),
	)
}
