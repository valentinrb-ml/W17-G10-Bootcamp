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

// ParseIntParam (ya existente, por si acaso)
func ParseIntParam(r *http.Request, name string) (int, error) {
	// Search first in URL params
	valueStr := chi.URLParam(r, name)

	// If it is not in URL, search in query params
	if valueStr == "" {
		valueStr = r.URL.Query().Get(name)
	}

	if valueStr == "" {
		return 0, apperrors.NewAppError(apperrors.CodeBadRequest, name+" parameter is required")
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, apperrors.NewAppError(apperrors.CodeBadRequest, name+" must be a valid integer")
	}

	return value, nil
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
