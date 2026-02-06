package appvalidator

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator implements Echo's Validator interface
type CustomValidator struct {
	validate *validator.Validate
}

// NewCustomValidator creates a new instance of CustomValidator
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{validate: validator.New()}
}

// Validate validates the request body
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}
