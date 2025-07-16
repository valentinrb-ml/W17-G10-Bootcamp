package httputil

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
)

var (
	ErrParamNotProvided = apperrors.NewAppError("BAD_REQUEST", "parameter not provided")
)

// ParseIntQueryParam parsea un parámetro de query opcional como entero
func ParseIntQueryParam(r *http.Request, name string) (*int, error) {
	valueStr := r.URL.Query().Get(name)
	if valueStr == "" {
		return nil, ErrParamNotProvided
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, name+" must be a valid integer")
	}

	return &value, nil
}

// getParamValue gets the value of a parameter (URL param or query param)
func getParamValue(r *http.Request, name string) string {
	// Search first in URL params
	value := chi.URLParam(r, name)

	// If it is not in URL, search in query params
	if value == "" {
		value = r.URL.Query().Get(name)
	}

	return value
}

// validateIntValue converts and validates that a string is a valid integer
func validateIntValue(valueStr, paramName string) (int, error) {
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, apperrors.NewAppError(apperrors.CodeBadRequest, paramName+" must be a valid integer")
	}
	return value, nil
}

// ParseIntParam parses an integer parameter REQUIRED
func ParseIntParam(r *http.Request, name string) (int, error) {
	valueStr := getParamValue(r, name)

	if valueStr == "" {
		return 0, apperrors.NewAppError(apperrors.CodeBadRequest, name+" parameter is required")
	}

	return validateIntValue(valueStr, name)
}

// ParseOptionalIntParam parse an integer parameter OPTIONAL
func ParseOptionalIntParam(r *http.Request, name string) (int, error) {
	valueStr := getParamValue(r, name)

	if valueStr == "" {
		return 0, nil // Returns 0 if not present
	}

	return validateIntValue(valueStr, name)
}

func HandleError(w http.ResponseWriter, err error) bool {
	if err != nil {
		response.Error(w, ConvertServiceErrorToAppError(err))
		return true
	}
	return false
}

func ConvertServiceErrorToAppError(err error) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *apperrors.AppError:
		return e
	case *api.ServiceError:
		// Mapear el código de ServiceError a AppError
		var code string
		switch e.ResponseCode {
		case http.StatusBadRequest:
			code = apperrors.CodeBadRequest
		case http.StatusUnauthorized:
			code = apperrors.CodeUnauthorized
		case http.StatusForbidden:
			code = apperrors.CodeForbidden
		case http.StatusNotFound:
			code = apperrors.CodeNotFound
		case http.StatusConflict:
			code = apperrors.CodeConflict
		case http.StatusUnprocessableEntity:
			code = apperrors.CodeValidationError
		default:
			code = apperrors.CodeInternal
		}
		return apperrors.NewAppError(code, e.Message)
	default:
		return apperrors.Wrap(err, "internal server error")
	}
}
