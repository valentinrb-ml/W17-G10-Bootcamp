package apperrors

import "net/http"

// Standard HTTP codes
const (
	// HTTP standard
	CodeBadRequest   = "BAD_REQUEST"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeNotFound     = "NOT_FOUND"
	CodeConflict     = "CONFLICT"
	CodeInternal     = "INTERNAL_ERROR"

	// Handler specific
	CodeMethodNotAllowed = "METHOD_NOT_ALLOWED"
	CodeValidationError  = "VALIDATION_ERROR"
)

// Mapping of codes to HTTP statuses
var codeToStatus = map[string]int{
	// HTTP standard
	CodeBadRequest:   http.StatusBadRequest,
	CodeUnauthorized: http.StatusUnauthorized,
	CodeForbidden:    http.StatusForbidden,
	CodeNotFound:     http.StatusNotFound,
	CodeConflict:     http.StatusConflict,
	CodeInternal:     http.StatusInternalServerError,

	// Handler specific
	CodeMethodNotAllowed: http.StatusMethodNotAllowed,    // 405
	CodeValidationError:  http.StatusUnprocessableEntity, // 422
}
