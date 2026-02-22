package middleware

import (
	"net/http"
	"strings"

	"honda-leasing-api/configs"
	"honda-leasing-api/pkg/crypto"
	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

// Auth is the middleware that validates JWT token
func Auth(cfg configs.JwtConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Missing Authorization header"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Invalid Authorization header format"))
			return
		}

		tokenString := parts[1]

		claims, err := crypto.ValidateToken(tokenString, cfg.Secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Invalid or expired token"))
			return
		}

		// Set user identity context for downstream handlers
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.RoleName)

		c.Next()
	}
}
