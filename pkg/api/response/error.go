package response

import (
	"encoding/json"
	"errors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"net/http"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}
type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error - Single point handling of HTTP errors
func Error(w http.ResponseWriter, err error) {
	if err == nil {
		writeErrorResponse(w, http.StatusInternalServerError, ErrorDetail{
			Code:    apperrors.CodeInternal,
			Message: "Unexpected nil error",
		})
		return
	}

	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		detail := ErrorDetail{
			Code:    appErr.Code,
			Message: appErr.Message,
		}

		// Only include details if they are not empty
		if len(appErr.Details) > 0 {
			detail.Details = appErr.Details
		}

		writeErrorResponse(w, appErr.HTTPStatus, detail)
		return
	}

	// Untyped error - fallback
	writeErrorResponse(w, http.StatusInternalServerError, ErrorDetail{
		Code:    apperrors.CodeInternal,
		Message: "internal server error",
	})
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, detail ErrorDetail) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{Error: detail}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Fallback if encoding fails
		http.Error(w, `{"error":{"code":"ENCODING_ERROR, "message":"failed to encode response"}}`, http.StatusInternalServerError)
	}
}
