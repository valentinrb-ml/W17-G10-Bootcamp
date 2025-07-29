package handler_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product"
	productmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestProductHandler_GetAll(t *testing.T) {
	t.Parallel()

	hResp := []models.ProductResponse{{ID: 1, ProductData: models.ProductData{ProductCode: "X"}}}

	tests := []struct {
		name      string
		mockSetup func(*productmock.MockService)
		status    int
		appCode   string
	}{
		{
			name: "success",
			mockSetup: func(s *productmock.MockService) {
				s.On("GetAll", mock.Anything).Return(hResp, nil).Once()
			},
			status: http.StatusOK,
		},
		{
			name: "service error",
			mockSetup: func(s *productmock.MockService) {
				var nilSlice []models.ProductResponse
				s.On("GetAll", mock.Anything).
					Return(nilSlice, apperrors.NewAppError(apperrors.CodeInternal, "db")).Once()
			},
			status:  http.StatusInternalServerError,
			appCode: apperrors.CodeInternal,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := &productmock.MockService{}
			tc.mockSetup(svc)

			h := handler.NewProductHandler(svc)
			rec := testhelpers.DoRequest(t, http.MethodGet, "/products", nil,
				h.GetAll)

			require.Equal(t, tc.status, rec.Code)
			if tc.appCode != "" {
				appErr, _ := testhelpers.DecodeAppErr(rec.Body)
				testhelpers.RequireAppErr(t, appErr, tc.appCode)
			}
			svc.AssertExpectations(t)
		})
	}
}
