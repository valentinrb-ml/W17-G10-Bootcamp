package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/seller"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/seller"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestSellerHandler_FindById(t *testing.T) {
	tests := []struct {
		name            string
		routeID         string
		mockService     func() *mocks.SellerServiceMock
		wantStatus      int
		wantResponse    *models.ResponseSeller
		wantErrorCode   string
		wantErrorMsgSub string
	}{
		{
			name:    "success",
			routeID: "1",
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.FindByIdFn = func(ctx context.Context, id int) (*models.ResponseSeller, error) {
					resp := testhelpers.DummyResponseSeller()
					return &resp, nil
				}
				return mock
			},
			wantStatus:   http.StatusOK,
			wantResponse: testhelpers.Ptr(testhelpers.DummyResponseSeller()),
		},
		{
			name:    "error - invalid id param (missing)",
			routeID: "",
			mockService: func() *mocks.SellerServiceMock {
				return &mocks.SellerServiceMock{}
			},
			wantStatus:      http.StatusBadRequest,
			wantErrorCode:   apperrors.CodeBadRequest,
			wantErrorMsgSub: "id parameter is required",
		},
		{
			name:    "error - not found",
			routeID: "999",
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.FindByIdFn = func(ctx context.Context, id int) (*models.ResponseSeller, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The seller you are looking for does not exist.")
				}
				return mock
			},
			wantStatus:      http.StatusNotFound,
			wantErrorCode:   apperrors.CodeNotFound,
			wantErrorMsgSub: "does not exist",
		},
		{
			name:    "error - internal error from repo",
			routeID: "1",
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.FindByIdFn = func(ctx context.Context, id int) (*models.ResponseSeller, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "internal server error")
				}
				return mock
			},
			wantStatus:      http.StatusInternalServerError,
			wantErrorCode:   apperrors.CodeInternal,
			wantErrorMsgSub: "internal server error",
		},
		{
			name:    "error - unknown error fallback",
			routeID: "2",
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.FindByIdFn = func(ctx context.Context, id int) (*models.ResponseSeller, error) {
					return nil, errors.New("unknown error happened")
				}
				return mock
			},
			wantStatus:      http.StatusInternalServerError,
			wantErrorCode:   apperrors.CodeInternal,
			wantErrorMsgSub: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiPath := "/api/v1/sellers"
			if tt.routeID != "" {
				apiPath += "/" + tt.routeID
			}
			req, err := http.NewRequest(http.MethodGet, apiPath, nil)
			require.NoError(t, err)

			// Set chi RouteContext for id param if present
			routeCtx := chi.NewRouteContext()
			if tt.routeID != "" {
				routeCtx.URLParams.Add("id", tt.routeID)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			rec := httptest.NewRecorder()
			h := handler.NewSellerHandler(tt.mockService())
			h.SetLogger(testhelpers.NewTestLogger())

			h.FindById(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				var envelope struct {
					Data models.ResponseSeller `json:"data"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &envelope)
				require.NoError(t, err)
				require.Equal(t, *tt.wantResponse, envelope.Data)
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
				require.Equal(t, tt.wantErrorCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.wantErrorMsgSub)
			}
		})
	}
}
