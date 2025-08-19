package response_test

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"net/http"
	"net/http/httptest"
	"testing"
)

func decodeErrorBody(t *testing.T, rec *httptest.ResponseRecorder) response.ErrorResponse {
	var errResp response.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	return errResp
}

func TestError(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		wantCode    int
		wantAppCode string
		wantMsgSub  string
		wantDetails map[string]interface{}
		wantRaw     string // Solo para casos de encoder fallback (si lo quieres explorar)
	}{
		{
			name:        "nil error gives internal server error",
			inputErr:    nil,
			wantCode:    http.StatusInternalServerError,
			wantAppCode: apperrors.CodeInternal,
			wantMsgSub:  "Unexpected nil error",
		},
		{
			name:        "typed AppError sets code, message, status",
			inputErr:    apperrors.NewAppError(apperrors.CodeBadRequest, "bad request"),
			wantCode:    http.StatusBadRequest,
			wantAppCode: apperrors.CodeBadRequest,
			wantMsgSub:  "bad request",
		},
		{
			name: "typed AppError with details sets details",
			inputErr: func() error {
				e := apperrors.NewAppError(apperrors.CodeValidationError, "validation error")
				e.Details = map[string]interface{}{"field": "name", "msg": "required"}
				return e
			}(),
			wantCode:    http.StatusUnprocessableEntity,
			wantAppCode: apperrors.CodeValidationError,
			wantMsgSub:  "validation error",
			wantDetails: map[string]interface{}{"field": "name", "msg": "required"},
		},
		{
			name:        "normal error gives 500",
			inputErr:    errors.New("some string error"),
			wantCode:    http.StatusInternalServerError,
			wantAppCode: apperrors.CodeInternal,
			wantMsgSub:  "internal server error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			response.Error(rec, tc.inputErr)

			require.Equal(t, tc.wantCode, rec.Code)
			if tc.wantRaw != "" {
				require.Contains(t, rec.Body.String(), tc.wantRaw)
				return
			}
			errResp := decodeErrorBody(t, rec)

			require.Equal(t, tc.wantAppCode, errResp.Error.Code)
			require.Contains(t, errResp.Error.Message, tc.wantMsgSub)
			if len(tc.wantDetails) > 0 {
				require.Equal(t, tc.wantDetails, errResp.Error.Details)
			}
		})
	}
}
