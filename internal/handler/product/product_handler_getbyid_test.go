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
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func withIDParam(r *http.Request, id string) *http.Request {
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}

func TestProductHandler_GetByID(t *testing.T) {
	t.Parallel()

	resp := models.ProductResponse{ID: 3}

	tests := []struct {
		name      string
		param     string
		mockSetup func(*productmock.MockService)
		status    int
		appCode   string
	}{
		{
			name:  "success",
			param: "3",
			mockSetup: func(s *productmock.MockService) {
				s.On("GetByID", mock.Anything, 3).Return(resp, nil).Once()
			},
			status: http.StatusOK,
		},
		{
			name:      "bad param",
			param:     "abc",
			status:    http.StatusBadRequest,
			appCode:   apperrors.CodeBadRequest,
			mockSetup: func(s *productmock.MockService) {},
		},
		{
			name:    "not found",
			param:   "99",
			status:  http.StatusNotFound,
			appCode: apperrors.CodeNotFound,
			mockSetup: func(s *productmock.MockService) {
				s.On("GetByID", mock.Anything, 99).
					Return(models.ProductResponse{},
						apperrors.NewAppError(apperrors.CodeNotFound, "no")).Once()
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
			req := testhelpers.NewRequest(t, http.MethodGet, "/products/"+tc.param, nil)
			req = withIDParam(req, tc.param)

			rec := testhelpers.DoRawRequest(t, req, http.HandlerFunc(h.GetByID))

			require.Equal(t, tc.status, rec.Code)
			if tc.appCode != "" {
				app, _ := testhelpers.DecodeAppErr(rec.Body)
				testhelpers.RequireAppErr(t, app, tc.appCode)
			}
			svc.AssertExpectations(t)
		})
	}
}
