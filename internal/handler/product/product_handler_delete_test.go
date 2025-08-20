package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product"
	productmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func withID(r *http.Request, id string) *http.Request {
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}

func TestProductHandler_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		param     string
		mockSetup func(*productmock.MockService)
		status    int
		appCode   string
	}{
		{
			name:  "success",
			param: "7",
			mockSetup: func(s *productmock.MockService) {
				s.On("Delete", mock.Anything, 7).Return(nil).Once()
			},
			status: http.StatusNoContent,
		},
		{
			name:      "bad id param",
			param:     "abc",
			status:    http.StatusBadRequest,
			appCode:   apperrors.CodeBadRequest,
			mockSetup: func(s *productmock.MockService) {},
		},
		{
			name:    "not found",
			param:   "22",
			status:  http.StatusNotFound,
			appCode: apperrors.CodeNotFound,
			mockSetup: func(s *productmock.MockService) {
				s.On("Delete", mock.Anything, 22).
					Return(apperrors.NewAppError(apperrors.CodeNotFound, "x")).Once()
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := &productmock.MockService{}
			tc.mockSetup(svc)

			h := handler.NewProductHandler(svc)
			h.SetLogger(testhelpers.NewTestLogger())
			req := testhelpers.NewRequest(t, http.MethodDelete, "/products/"+tc.param, nil)
			req = withID(req, tc.param)

			rec := testhelpers.DoRawRequest(t, req, http.HandlerFunc(h.Delete))

			require.Equal(t, tc.status, rec.Code)
			if tc.appCode != "" {
				app, _ := testhelpers.DecodeAppErr(rec.Body)
				testhelpers.RequireAppErr(t, app, tc.appCode)
			}
			svc.AssertExpectations(t)
		})
	}
}
