package httputil

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
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
