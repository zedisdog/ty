package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/zedisdog/ty/errx"
	"github.com/zedisdog/ty/sdk/net/http/response"
)

func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.ErrorAndAbort(c, errx.NewWithCode(http.StatusForbidden, "rate limit"))
			return
		}
		c.Next()
	}
}
