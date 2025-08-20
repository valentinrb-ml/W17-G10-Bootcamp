package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/buyer"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestBuyerHandler_Update(t *testing.T) {
	type args struct {
		requestBody any
		routeID     string
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
				routeID:     "1",
			},
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.UpdateFn = func(ctx context.Context, id int, req models.RequestBuyer) (*models.ResponseBuyer, error) {
					expected := testhelpers.DummyResponseBuyer()
					return &expected, nil
				}
				return mock
			},
			wantStatus:   http.StatusOK,
			responseBody: testhelpers.DummyResponseBuyer(),
		},
		{
			name: "error - invalid id param (missing)",
			args: args{
				requestBody: testhelpers.DummyRequestBuyer(),
				routeID:     "",
			},
			mockService: func() *mocks.BuyerServiceMock {
				return &mocks.BuyerServiceMock{}
			},
			wantStatus:   http.StatusBadRequest,
			errorCode:    apperrors.CodeBadRequest,
			errorMessage: "Invalid ID parameter",
		},
		{
			name: "error - invalid request payload",
			args: args{
				requestBody: `{invalid json`,
				routeID:     "1",
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
				routeID:     "1",
			},
			mockService: func() *mocks.BuyerServiceMock {
				return &mocks.BuyerServiceMock{}
			},
			wantStatus:   http.StatusUnprocessableEntity,
			errorCode:    apperrors.CodeValidationError,
			errorMessage: "required",
		},
		{
			name: "error - not found",
			args: args{
				requestBody: testhelpers.DummyRequestBuyer(),
				routeID:     "999",
			},
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.UpdateFn = func(ctx context.Context, id int, req models.RequestBuyer) (*models.ResponseBuyer, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The buyer you are trying to update does not exist")
				}
				return mock
			},
			wantStatus:   http.StatusNotFound,
			errorCode:    apperrors.CodeNotFound,
			errorMessage: "does not exist",
		},
		{
			name: "error - service layer returns error",
			args: args{
				requestBody: testhelpers.DummyRequestBuyer(),
				routeID:     "1",
			},
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.UpdateFn = func(ctx context.Context, id int, req models.RequestBuyer) (*models.ResponseBuyer, error) {
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

			req, err := http.NewRequest(http.MethodPatch, "/api/v1/buyers", bytes.NewReader(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			routeCtx := chi.NewRouteContext()
			if tt.args.routeID != "" {
				routeCtx.URLParams.Add("id", tt.args.routeID)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			rec := httptest.NewRecorder()
			h := handler.NewBuyerHandler(tt.mockService())
			h.SetLogger(testhelpers.NewTestLogger())

			h.Update(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
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
