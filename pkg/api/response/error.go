package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
)

type ErrorResponse struct {
	Error     ErrorDetail `json:"error"`
	RequestID string      `json:"request_id,omitempty"`
}
type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error - Single point handling of HTTP errors
func Error(w http.ResponseWriter, err error) {
	ErrorWithRequest(w, nil, err)
}

// ErrorWithRequest - Handle HTTP errors with request context for Request ID
func ErrorWithRequest(w http.ResponseWriter, r *http.Request, err error) {
	var requestID string
	if r != nil {
		if id, ok := r.Context().Value("request_id").(string); ok {
			requestID = id
		}
	}

	if err == nil {
		writeErrorResponse(w, http.StatusInternalServerError, ErrorDetail{
			Code:    apperrors.CodeInternal,
			Message: "Unexpected nil error",
		}, requestID)
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

		writeErrorResponse(w, appErr.HTTPStatus, detail, requestID)
		return
	}

	// Untyped error - fallback
	writeErrorResponse(w, http.StatusInternalServerError, ErrorDetail{
		Code:    apperrors.CodeInternal,
		Message: "internal server error",
	}, requestID)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, detail ErrorDetail, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:     detail,
		RequestID: requestID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Fallback if encoding fails
		http.Error(w, `{"error":{"code":"ENCODING_ERROR", "message":"failed to encode response"}}`, http.StatusInternalServerError)
	}
}
