package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using go-playground/validator tags.
func ValidateStruct(obj interface{}) error {
	return validate.Struct(obj)
}
