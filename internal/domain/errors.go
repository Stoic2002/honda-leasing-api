package domain

import "errors"

var (
	ErrNotFound            = errors.New("resource not found")
	ErrInvalidInput        = errors.New("invalid input data")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrForbidden           = errors.New("forbidden action")
	ErrConflict            = errors.New("resource already exists")
	ErrInternalServerError = errors.New("internal server error")
)
