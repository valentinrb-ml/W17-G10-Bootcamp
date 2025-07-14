package httputil

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"net/http"
	"strconv"
)

func DecodeJSON(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return apperrors.BadRequest("request body is required")
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperrors.BadRequest("invalid JSON format")
	}
	return nil
}

func ParseIDParam(r *http.Request, paramName string) (int, error) {
	idStr := chi.URLParam(r, paramName)
	if idStr == "" {
		return 0, apperrors.BadRequest(paramName + " parameter is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, apperrors.BadRequest(paramName + " must be a valid integer")
	}

	if id <= 0 {
		return 0, apperrors.BadRequest(paramName + " must be a positive integer")
	}

	return id, nil
}
