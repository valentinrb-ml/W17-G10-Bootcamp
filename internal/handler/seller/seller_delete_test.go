package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/seller"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/seller"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
)

func TestSellerHandler_Delete(t *testing.T) {
	tests := []struct {
		name            string
		routeID         string
		mockService     func() *mocks.SellerServiceMock
		wantStatus      int
		wantErrorCode   string
		wantErrorMsgSub string
	}{
		{
			name:    "success",
			routeID: "1",
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.DeleteFn = func(ctx context.Context, id int) error {
					return nil
				}
				return mock
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:    "error - missing id param",
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
				mock.DeleteFn = func(ctx context.Context, id int) error {
					return apperrors.NewAppError(apperrors.CodeNotFound, "The seller you are trying to delete does not exist")
				}
				return mock
			},
			wantStatus:      http.StatusNotFound,
			wantErrorCode:   apperrors.CodeNotFound,
			wantErrorMsgSub: "does not exist",
		},
		{
			name:    "error - conflict (associated products)",
			routeID: "777",
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.DeleteFn = func(ctx context.Context, id int) error {
					return apperrors.NewAppError(apperrors.CodeConflict, "Cannot delete seller: there are products associated with this seller.")
				}
				return mock
			},
			wantStatus:      http.StatusConflict,
			wantErrorCode:   apperrors.CodeConflict,
			wantErrorMsgSub: "associated with this seller",
		},
		{
			name:    "error - internal error",
			routeID: "2",
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.DeleteFn = func(ctx context.Context, id int) error {
					return apperrors.NewAppError(apperrors.CodeInternal, "unexpected internal error")
				}
				return mock
			},
			wantStatus:      http.StatusInternalServerError,
			wantErrorCode:   apperrors.CodeInternal,
			wantErrorMsgSub: "unexpected internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Prepare request
			apiPath := "/api/v1/sellers"
			if tt.routeID != "" {
				apiPath += "/" + tt.routeID
			}
			req, err := http.NewRequest(http.MethodDelete, apiPath, nil)
			require.NoError(t, err)

			// Add chi RouteContext with "id"
			routeCtx := chi.NewRouteContext()
			if tt.routeID != "" {
				routeCtx.URLParams.Add("id", tt.routeID)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			rec := httptest.NewRecorder()
			h := handler.NewSellerHandler(tt.mockService())

			h.Delete(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus != http.StatusNoContent {
				var body struct {
					Error struct {
						Code    string         `json:"code"`
						Message string         `json:"message"`
						Details map[string]any `json:"details"`
					} `json:"error"`
				}
				err := json.NewDecoder(rec.Body).Decode(&body)
				require.NoError(t, err)
				require.Equal(t, tt.wantErrorCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.wantErrorMsgSub)
			} else {
				require.Empty(t, rec.Body.String())
			}
		})
	}
}
