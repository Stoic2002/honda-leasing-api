package middleware

import (
	"github.com/gin-gonic/gin"
)

// RateLimiter acts as a placeholder for a real redis/memory rate limiter
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// e.g. check client IP in Redis, if exceeded abort:
		// c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Error(http.StatusTooManyRequests, "Rate limit exceeded"))
		// return
		c.Next()
	}
}
