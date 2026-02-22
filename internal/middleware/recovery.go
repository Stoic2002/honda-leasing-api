package middleware

import (
	"net/http"

	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

func GlobalRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Internal Server Error: "+err))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Internal Server Error - Panic Recovered"))
	})
}
