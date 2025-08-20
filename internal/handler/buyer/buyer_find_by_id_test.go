package handler_test

import (
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

func TestBuyerHandler_FindById(t *testing.T) {
	tests := []struct {
		name         string
		routeID      string
		mockService  func() *mocks.BuyerServiceMock
		wantStatus   int
		response     *models.ResponseBuyer
		errorCode    string
		errorMessage string
	}{
		{
			name:    "success",
			routeID: "1",
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.FindByIdFn = func(ctx context.Context, id int) (*models.ResponseBuyer, error) {
					resp := testhelpers.DummyResponseBuyer()
					return &resp, nil
				}
				return mock
			},
			wantStatus: http.StatusOK,
			response:   testhelpers.PtrBuyer(testhelpers.DummyResponseBuyer()),
		},
		{
			name:    "error - invalid id param (missing)",
			routeID: "",
			mockService: func() *mocks.BuyerServiceMock {
				return &mocks.BuyerServiceMock{}
			},
			wantStatus:   http.StatusBadRequest,
			errorCode:    apperrors.CodeBadRequest,
			errorMessage: "Invalid ID parameter",
		},
		{
			name:    "error - invalid id param (non-numeric)",
			routeID: "abc",
			mockService: func() *mocks.BuyerServiceMock {
				return &mocks.BuyerServiceMock{}
			},
			wantStatus:   http.StatusBadRequest,
			errorCode:    apperrors.CodeBadRequest,
			errorMessage: "Invalid ID parameter",
		},
		{
			name:    "error - not found",
			routeID: "999",
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.FindByIdFn = func(ctx context.Context, id int) (*models.ResponseBuyer, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The buyer you are looking for does not exist")
				}
				return mock
			},
			wantStatus:   http.StatusNotFound,
			errorCode:    apperrors.CodeNotFound,
			errorMessage: "does not exist",
		},
		{
			name:    "error - internal error",
			routeID: "1",
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.FindByIdFn = func(ctx context.Context, id int) (*models.ResponseBuyer, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "internal server error")
				}
				return mock
			},
			wantStatus:   http.StatusInternalServerError,
			errorCode:    apperrors.CodeInternal,
			errorMessage: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiPath := "/api/v1/buyers"
			if tt.routeID != "" {
				apiPath += "/" + tt.routeID
			}
			req, err := http.NewRequest(http.MethodGet, apiPath, nil)
			require.NoError(t, err)

			routeCtx := chi.NewRouteContext()
			if tt.routeID != "" {
				routeCtx.URLParams.Add("id", tt.routeID)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			rec := httptest.NewRecorder()
			h := handler.NewBuyerHandler(tt.mockService())
			h.SetLogger(testhelpers.NewTestLogger())

			h.FindById(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				var envelope struct {
					Data models.ResponseBuyer `json:"data"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &envelope)
				require.NoError(t, err)
				require.Equal(t, *tt.response, envelope.Data)
			} else {
				var body struct {
					Error struct {
						Code    string         `json:"code"`
						Message string         `json:"message"`
						Details map[string]any `json:"details"`
					} `json:"error"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				require.NoError(t, err)
				require.Equal(t, tt.errorCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.errorMessage)
			}
		})
	}
}
