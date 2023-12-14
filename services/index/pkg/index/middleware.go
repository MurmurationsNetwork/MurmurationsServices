package index

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AllowInNonProductionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Default to restrictive behavior (suitable for production)
		// Only allow the operation in non-production environments if explicitly specified.
		if os.Getenv("APP_ENV") != "development" &&
			os.Getenv("APP_ENV") != "staging" &&
			os.Getenv("APP_ENV") != "pretest" {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "Operation not allowed in this environment"},
			)
			return
		}
		c.Next()
	}
}
