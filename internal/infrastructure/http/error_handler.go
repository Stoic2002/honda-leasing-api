package http

import (
	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			statusCode, msg := response.MapDomainError(err)
			c.JSON(statusCode, response.Error(statusCode, msg))
		}
	}
}
