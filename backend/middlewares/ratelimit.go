package middlewares

import (
	"github.com/gin-gonic/gin"
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement rate limiter middleware
		c.Next()
	}
}
