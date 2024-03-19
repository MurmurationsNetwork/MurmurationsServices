package index

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AllowInNonLiveMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Default to restrictive behavior (suitable for live)
		// Only allow the operation in non-live environments if explicitly specified.
		if os.Getenv("APP_ENV") != "dev" &&
			os.Getenv("APP_ENV") != "live-test" &&
			os.Getenv("APP_ENV") != "ci" {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "Operation not allowed in this environment"},
			)
			return
		}
		c.Next()
	}
}
