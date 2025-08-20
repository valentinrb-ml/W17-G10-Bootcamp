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
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestBuyerHandler_Delete(t *testing.T) {
	tests := []struct {
		name         string
		routeID      string
		mockService  func() *mocks.BuyerServiceMock
		wantStatus   int
		errorCode    string
		errorMessage string
	}{
		{
			name:    "success",
			routeID: "1",
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
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
				mock.DeleteFn = func(ctx context.Context, id int) error {
					return apperrors.NewAppError(apperrors.CodeNotFound, "The buyer you are trying to delete does not exist")
				}
				return mock
			},
			wantStatus:   http.StatusNotFound,
			errorCode:    apperrors.CodeNotFound,
			errorMessage: "does not exist",
		},
		{
			name:    "error - conflict (associated purchases)",
			routeID: "2",
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.DeleteFn = func(ctx context.Context, id int) error {
					return apperrors.NewAppError(apperrors.CodeConflict, "Cannot delete buyer: there are purchases associated with this buyer")
				}
				return mock
			},
			wantStatus:   http.StatusConflict,
			errorCode:    apperrors.CodeConflict,
			errorMessage: "associated with this buyer",
		},
		{
			name:    "error - internal error",
			routeID: "3",
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.DeleteFn = func(ctx context.Context, id int) error {
					return apperrors.NewAppError(apperrors.CodeInternal, "unexpected internal error")
				}
				return mock
			},
			wantStatus:   http.StatusInternalServerError,
			errorCode:    apperrors.CodeInternal,
			errorMessage: "unexpected internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiPath := "/api/v1/buyers"
			if tt.routeID != "" {
				apiPath += "/" + tt.routeID
			}
			req, err := http.NewRequest(http.MethodDelete, apiPath, nil)
			require.NoError(t, err)

			routeCtx := chi.NewRouteContext()
			if tt.routeID != "" {
				routeCtx.URLParams.Add("id", tt.routeID)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			rec := httptest.NewRecorder()
			h := handler.NewBuyerHandler(tt.mockService())
			h.SetLogger(testhelpers.NewTestLogger())

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
				require.Equal(t, tt.errorCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.errorMessage)
			} else {
				require.Empty(t, rec.Body.String())
			}
		})
	}
}
