package http

import (
	"fmt"
	"net/http"

	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			// You can switch based on error types from internal/domain/errors.go
			// For generic fallback:
			msg := fmt.Sprintf("Request failed: %v", err)
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, msg))
		}
	}
}
