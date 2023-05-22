package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type options struct {
	Headers []string
	Methods []string
}

func WithAllowHeaders(fields []string) func(*options) {
	return func(o *options) {
		o.Headers = fields
	}
}

func GenCros(opts ...func(*options)) gin.HandlerFunc {
	options := &options{}
	for _, opt := range opts {
		opt(options)
	}
	appendHeaders := ""
	if len(options.Headers) > 0 {
		appendHeaders += ", "
		appendHeaders += strings.Join(options.Headers, ", ")
	}
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control, X-Requested-With"+appendHeaders)
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// Deprecated: Use GenCros instead.
func Cros(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	c.Next()
}
