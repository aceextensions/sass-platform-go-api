package apperrors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aceextension/core/logger"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AppError represents a custom application error
type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// MapDBError converts a database error into a friendly AppError
func MapDBError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return NewAppError(http.StatusConflict, "This record already exists")
		case "23503": // foreign_key_violation
			return NewAppError(http.StatusBadRequest, "Related record not found")
		case "23502": // not_null_violation
			return NewAppError(http.StatusBadRequest, "Missing required fields")
		}
	}

	// Default to internal server error if not handled
	return NewAppError(http.StatusInternalServerError, "An unexpected database error occurred")
}

// formatValidationError converts a technical field error into a human-readable message
func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters", e.Field(), e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", e.Field(), e.Param())
	case "alphanum":
		return fmt.Sprintf("%s can only contain letters and numbers", e.Field())
	}
	return fmt.Sprintf("%s failed on the '%s' tag", e.Field(), e.Tag())
}

// GlobalErrorHandler handles all Echo errors
func GlobalErrorHandler(err error, c echo.Context) {
	var appErr *AppError
	var echoErr *echo.HTTPError
	var validationErr validator.ValidationErrors

	code := http.StatusInternalServerError
	var response interface{}
	response = map[string]string{"error": "An internal server error occurred"}

	if errors.As(err, &appErr) {
		code = appErr.Code
		response = map[string]string{"error": appErr.Message}
	} else if errors.As(err, &validationErr) {
		code = http.StatusBadRequest
		errs := make([]string, 0, len(validationErr))
		for _, e := range validationErr {
			errs = append(errs, formatValidationError(e))
		}
		response = map[string]interface{}{"errors": errs}
	} else if errors.As(err, &echoErr) {
		code = echoErr.Code
		message := "An internal server error occurred"
		if m, ok := echoErr.Message.(string); ok {
			message = m
		}
		response = map[string]string{"error": message}
	}

	// Log technical details (will go to logs/error.log if it's an error)
	if code >= 500 {
		logger.Log.Error("HTTP Error",
			zap.Int("code", code),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Path()),
			zap.Error(err),
		)
	} else {
		logger.Log.Info("HTTP Request handled",
			zap.Int("code", code),
			zap.String("path", c.Path()),
		)
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, response)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
