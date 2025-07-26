package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/buyer"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestBuyerHandler_Create(t *testing.T) {
	type args struct {
		requestBody any
	}

	tests := []struct {
		name         string
		args         args
		mockService  func() *mocks.BuyerServiceMock
		wantStatus   int
		responseBody any
		errorCode    string
		errorMessage string
	}{
		{
			name: "success",
			args: args{
				requestBody: testhelpers.DummyRequestBuyer(),
			},
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.CreateFn = func(ctx context.Context, req models.RequestBuyer) (*models.ResponseBuyer, error) {
					expected := testhelpers.DummyResponseBuyer()
					return &expected, nil
				}
				return mock
			},
			wantStatus:   http.StatusCreated,
			responseBody: testhelpers.DummyResponseBuyer(),
		},
		{
			name: "error - invalid request payload",
			args: args{
				requestBody: `{invalid json}`,
			},
			mockService: func() *mocks.BuyerServiceMock {
				return &mocks.BuyerServiceMock{}
			},
			wantStatus:   http.StatusBadRequest,
			errorCode:    apperrors.CodeBadRequest,
			errorMessage: "Invalid request body",
		},
		{
			name: "error - invalid Buyer validation",
			args: args{
				requestBody: models.RequestBuyer{},
			},
			mockService: func() *mocks.BuyerServiceMock {
				return &mocks.BuyerServiceMock{}
			},
			wantStatus:   http.StatusUnprocessableEntity,
			errorCode:    apperrors.CodeValidationError,
			errorMessage: "required",
		},
		{
			name: "error - service layer returns error",
			args: args{
				requestBody: testhelpers.DummyRequestBuyer(),
			},
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.CreateFn = func(ctx context.Context, req models.RequestBuyer) (*models.ResponseBuyer, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "unexpected error")
				}
				return mock
			},
			wantStatus:   http.StatusInternalServerError,
			errorCode:    apperrors.CodeInternal,
			errorMessage: "unexpected error",
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

			req, err := http.NewRequest(http.MethodPost, "/api/v1/buyers", bytes.NewReader(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h := handler.NewBuyerHandler(tt.mockService())

			h.Create(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusCreated {
				var responseEnvelope struct {
					Data models.ResponseBuyer `json:"data"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &responseEnvelope)
				require.NoError(t, err)
				require.Equal(t, tt.responseBody, responseEnvelope.Data)
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
				require.Equal(t, tt.errorCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.errorMessage)
			}
		})
	}
}
