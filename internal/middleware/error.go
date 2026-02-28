package middleware

import (
	"errors"
	"net/http"

	"honda-leasing-api/internal/domain"
	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a global middleware that translates internal domain errors into appropriate HTTP responses.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Proceed to next handler
		c.Next()

		// If there are errors added to the context by handlers...
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var code int

			// Map domain errors to HTTP status codes
			switch {
			case errors.Is(err, domain.ErrNotFound):
				code = http.StatusNotFound
			case errors.Is(err, domain.ErrInvalidInput):
				code = http.StatusBadRequest
			case errors.Is(err, domain.ErrUnauthorized):
				code = http.StatusUnauthorized
			case errors.Is(err, domain.ErrForbidden):
				code = http.StatusForbidden
			case errors.Is(err, domain.ErrConflict):
				code = http.StatusConflict
			case err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password":
				// Handle specific bcrypt error for wrong password
				code = http.StatusUnauthorized
				err = errors.New("invalid email or password")
			default:
				code = http.StatusInternalServerError
			}

			// If response was not already written, write the JSON error response
			if !c.Writer.Written() {
				c.JSON(code, response.Error(code, err.Error()))
			}
		}
	}
}
