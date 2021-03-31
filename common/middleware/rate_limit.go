package middleware

import (
	"net/http"
	"strconv"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

type respond struct {
	Data interface{} `json:"data,omitempty"`
}

func RateLimit(period string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(period)
	if err != nil {
		logger.Panic("Error when trying to parse rate limit period", err)
		return nil
	}
	store := memory.NewStore()
	ipRateLimiter := limiter.New(store, rate)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		context, err := ipRateLimiter.Get(c, ip)
		if err != nil {
			logger.Error("Error when trying to get ipRateLimiter context", err)
			c.JSON(http.StatusInternalServerError, nil)
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		if context.Reached {
			c.JSON(http.StatusTooManyRequests, resterr.NewTooManyRequestsError("Rate limit exceeded"))
			c.Abort()
			return
		}

		c.Next()
	}
}
