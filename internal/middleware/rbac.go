package middleware

import (
	"net/http"

	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

func RoleBasedAccessControl(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error(http.StatusForbidden, "Role context missing"))
			return
		}

		userRole := role.(string)
		permitted := false
		for _, r := range allowedRoles {
			if userRole == r {
				permitted = true
				break
			}
		}

		if !permitted {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error(http.StatusForbidden, "Insufficient permissions"))
			return
		}

		c.Next()
	}
}
