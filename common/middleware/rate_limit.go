package middleware

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimit(period string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(period)
	if err != nil {
		logger.Panic("Error when trying to parse rate limit period", err)
		return nil
	}
	store := memory.NewStore()
	return mgin.NewMiddleware(limiter.New(store, rate))
}
