package apperrors

import (
	"errors"
	"fmt"
)

type AppError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	HTTPStatus int                    `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Main Constructor
func NewAppError(code, message string) *AppError {
	status, exists := codeToStatus[code]
	if !exists {
		// Warning for unregistered codes
		fmt.Printf("WARNING: Unknown error code '%s', using default status\n", code)
		status = 500
	}

	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
		Details:    make(map[string]interface{}),
	}
}

// Constructors for common cases
func BadRequest(message string) *AppError {
	return NewAppError(CodeBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return NewAppError(CodeUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return NewAppError(CodeForbidden, message)
}

func NotFound(message string) *AppError {
	return NewAppError(CodeNotFound, message)
}

func Conflict(message string) *AppError {
	return NewAppError(CodeConflict, message)
}

func Internal(message string) *AppError {
	return NewAppError(CodeInternal, message)
}

// Wrap for external errors
func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return Internal(message)
}

// WithDetail - Thread-safe for correct functionality in asynchronism
func (e *AppError) WithDetail(key string, value interface{}) *AppError {
	// Create new map (does not modify the original)
	details := make(map[string]interface{}, len(e.Details)+1)
	for k, v := range e.Details {
		details[k] = v
	}
	details[key] = value

	// Return new instance
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		HTTPStatus: e.HTTPStatus,
		Details:    details,
	}
}
