package httputil

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"net/http"
	"strconv"
)

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
