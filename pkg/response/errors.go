package response

import (
	"errors"
	"net/http"

	"honda-leasing-api/internal/domain"
)

// MapDomainError maps a domain error to an appropriate HTTP status code and message.
// Returns the HTTP status code and a safe, client-facing error message.
func MapDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "Resource not found"
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest, "Invalid input data"
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, "Unauthorized access"
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, "Forbidden action"
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict, "Resource already exists"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
