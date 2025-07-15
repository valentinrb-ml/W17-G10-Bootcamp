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
	valueStr := chi.URLParam(r, name)
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
		return apperrors.NewAppError(MapServiceErrorCode(e.Code), e.Message)
	default:
		return apperrors.Wrap(err, "internal server error")
	}
}

func MapServiceErrorCode(code int) string {
	switch code {
	case http.StatusBadRequest:
		return apperrors.CodeBadRequest
	case http.StatusUnauthorized:
		return apperrors.CodeUnauthorized
	case http.StatusForbidden:
		return apperrors.CodeForbidden
	case http.StatusNotFound:
		return apperrors.CodeNotFound
	case http.StatusConflict:
		return apperrors.CodeConflict
	case http.StatusUnprocessableEntity:
		return apperrors.CodeValidationError
	default:
		return apperrors.CodeBadRequest
	}
}
