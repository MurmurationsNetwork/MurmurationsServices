package limiter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
)

const (
	defaultPeriod = "20-M"
	defaultMethod = "*"
)

type RateLimitOptions struct {
	Period string
	Method string
}

// NewStore creates a new instance of ratel imit with defaults.
func NewRateLimit() gin.HandlerFunc {
	return NewRateLimitWithOptions(RateLimitOptions{
		Period: defaultPeriod,
		Method: defaultMethod,
	})
}

func NewRateLimitWithOptions(options RateLimitOptions) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(options.Period)
	if err != nil {
		logger.Panic("Error when trying to parse rate limit period", err)
		return nil
	}
	store := memory.NewStore()
	ipRateLimiter := limiter.New(store, rate)

	return func(c *gin.Context) {
		if c.Request.Method == options.Method ||
			options.Method == defaultMethod {
			ip := c.ClientIP()

			context, err := ipRateLimiter.Get(c, ip)
			if err != nil {
				logger.Error(
					"Error when trying to get ipRateLimiter context",
					err,
				)
				errors := jsonapi.NewError(
					[]string{"Internal Server Error"},
					[]string{
						"An internal server error was triggered and has been logged. Please try your request again later.",
					},
					nil,
					[]int{http.StatusInternalServerError},
				)
				res := jsonapi.Response(nil, errors, nil, nil)
				c.JSON(errors[0].Status, res)
				c.Abort()
				return
			}

			c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
			c.Header(
				"X-RateLimit-Remaining",
				strconv.FormatInt(context.Remaining, 10),
			)
			c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

			if context.Reached {
				errors := jsonapi.NewError(
					[]string{"Too Many Requests"},
					[]string{
						"You have exceeded the maximum number of requests per minute/hour. Please try again later. For more information see: https://docs.murmurations.network/developers/rate-limits.html",
					},
					nil,
					[]int{http.StatusTooManyRequests},
				)
				res := jsonapi.Response(nil, errors, nil, nil)
				c.JSON(errors[0].Status, res)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
