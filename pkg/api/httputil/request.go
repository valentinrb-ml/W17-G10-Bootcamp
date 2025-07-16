package httputil

import (
	"encoding/json"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"net/http"
)

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, apperrors.NewAppError(apperrors.CodeMethodNotAllowed, "method not allowed"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, apperrors.NewAppError(apperrors.CodeNotFound, "endpoint not found"))
}

func DecodeJSON(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return apperrors.NewAppError(apperrors.CodeBadRequest, "request body is required")
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperrors.NewAppError(apperrors.CodeBadRequest, "invalid JSON format")
	}
	return nil
}
