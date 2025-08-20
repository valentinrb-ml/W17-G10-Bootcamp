package handler_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product"
	productMappers "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers/product"
	productmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestProductHandler_Create(t *testing.T) {
	t.Parallel()

	validReq := models.ProductRequest{ProductData: models.ProductData{
		ProductCode: "X", Description: "d", Width: 1, Height: 1, Length: 1,
		NetWeight: 10, ExpirationRate: 5, RecommendedFreezingTemperature: 2,
		FreezingRate: 6, ProductTypeID: 1, SellerID: func() *int { i := 2; return &i }(),
	}}
	validDomain := productMappers.ToDomain(validReq)
	validResp := models.ProductResponse{ID: 5, ProductData: validReq.ProductData}

	tests := []struct {
		name      string
		body      []byte
		mockSetup func(*productmock.MockService)
		status    int
		appCode   string
	}{
		{
			name: "success",
			body: testhelpers.MustJSON(validReq),
			mockSetup: func(s *productmock.MockService) {
				s.On("Create", mock.Anything, validDomain).Return(validResp, nil).Once()
			},
			status: http.StatusCreated,
		},
		{
			name:      "invalid json",
			body:      []byte(`{"malformed":`),
			status:    http.StatusBadRequest,
			appCode:   apperrors.CodeBadRequest,
			mockSetup: func(s *productmock.MockService) {},
		},
		{
			name:      "validation error",
			body:      testhelpers.MustJSON(models.ProductRequest{}),
			status:    http.StatusUnprocessableEntity,
			appCode:   "VALIDATION_ERROR",
			mockSetup: func(s *productmock.MockService) {},
		},
		{
			name: "conflict from service",
			body: testhelpers.MustJSON(validReq),
			mockSetup: func(s *productmock.MockService) {
				s.On("Create", mock.Anything, validDomain).
					Return(models.ProductResponse{},
						apperrors.NewAppError(apperrors.CodeConflict, "dup")).Once()
			},
			status:  http.StatusConflict,
			appCode: apperrors.CodeConflict,
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
			req := testhelpers.NewRequest(t, http.MethodPost, "/products", bytes.NewReader(tc.body))
			rec := testhelpers.DoRawRequest(t, req, http.HandlerFunc(h.Create))

			require.Equal(t, tc.status, rec.Code)
			if tc.appCode != "" {
				app, _ := testhelpers.DecodeAppErr(rec.Body)
				testhelpers.RequireAppErr(t, app, tc.appCode)
			}
			svc.AssertExpectations(t)
		})
	}
}
