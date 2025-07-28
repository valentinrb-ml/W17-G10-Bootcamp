package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/geography"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestGeographyHandler_Create(t *testing.T) {
	type args struct {
		requestBody any
	}

	tests := []struct {
		name             string
		args             args
		mockService      func() *mocks.GeographyServiceMock
		wantStatus       int
		wantResponseBody any
		wantErrorCode    string
		wantErrorMsgSub  string
	}{
		{
			name: "success",
			args: args{
				requestBody: testhelpers.DummyRequestGeography(),
			},
			mockService: func() *mocks.GeographyServiceMock {
				mock := &mocks.GeographyServiceMock{}
				mock.CreateFn = func(ctx context.Context, req models.RequestGeography) (*models.ResponseGeography, error) {
					expected := testhelpers.DummyResponseGeography()
					return &expected, nil
				}
				return mock
			},
			wantStatus:       http.StatusCreated,
			wantResponseBody: testhelpers.DummyResponseGeography(),
		},
		{
			name: "error - invalid request payload",
			args: args{
				requestBody: `{invalid json}`,
			},
			mockService: func() *mocks.GeographyServiceMock {
				return &mocks.GeographyServiceMock{}
			},
			wantStatus:      http.StatusBadRequest,
			wantErrorCode:   apperrors.CodeBadRequest,
			wantErrorMsgSub: "invalid JSON",
		},
		{
			name: "error - validation error",
			args: args{
				requestBody: models.RequestGeography{},
			},
			mockService: func() *mocks.GeographyServiceMock {
				return &mocks.GeographyServiceMock{}
			},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorCode:   apperrors.CodeValidationError,
			wantErrorMsgSub: "required",
		},
		{
			name: "error - service layer returns error",
			args: args{
				requestBody: testhelpers.DummyRequestGeography(),
			},
			mockService: func() *mocks.GeographyServiceMock {
				mock := &mocks.GeographyServiceMock{}
				mock.CreateFn = func(ctx context.Context, req models.RequestGeography) (*models.ResponseGeography, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "unexpected service error")
				}
				return mock
			},
			wantStatus:      http.StatusInternalServerError,
			wantErrorCode:   apperrors.CodeInternal,
			wantErrorMsgSub: "unexpected service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBodyBytes []byte
			switch val := tt.args.requestBody.(type) {
			case string:
				requestBodyBytes = []byte(val)
			default:
				b, err := json.Marshal(tt.args.requestBody)
				require.NoError(t, err)
				requestBodyBytes = b
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/geography", bytes.NewReader(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h := handler.NewGeographyHandler(tt.mockService())

			h.Create(rec, req)
			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusCreated {
				var responseEnvelope struct {
					Data models.ResponseGeography `json:"data"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &responseEnvelope)
				require.NoError(t, err)
				require.Equal(t, tt.wantResponseBody, responseEnvelope.Data)
			} else {
				var body struct {
					Error struct {
						Code    string         `json:"code"`
						Message string         `json:"message"`
						Details map[string]any `json:"details"`
					} `json:"error"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &body)
				require.NoError(t, err)
				require.Equal(t, tt.wantErrorCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.wantErrorMsgSub)
			}
		})
	}
}
