package handler_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_record"
	productrecordmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

var sampleReport = []models.ProductRecordReport{
	{ProductID: 1, Description: "p1", RecordsCount: 2},
}

/*
Test for GET /productRecords/report handler.

Covered scenarios
  - 200 OK  (all products)
  - 200 OK  (single product)
  - 400 BAD_REQUEST     – invalid / negative id param
  - 404 NOT_FOUND       – service returns not-found
  - 500 INTERNAL_ERROR  – service failure
*/
func TestProductRecordHandler_GetRecordsReport(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		url       string
		mockSetup func(*productrecordmock.MockProductRecordService)
		status    int
		appErr    string // expected AppError.Code ("" == success)
	}{
		{
			name: "all_products_success",
			url:  "/productRecords/report",
			mockSetup: func(m *productrecordmock.MockProductRecordService) {
				m.On("GetRecordsReport", mock.Anything, 0).Return(sampleReport, nil).Once()
			},
			status: http.StatusOK,
		},
		{
			name: "single_product_success",
			url:  "/productRecords/report?id=1",
			mockSetup: func(m *productrecordmock.MockProductRecordService) {
				m.On("GetRecordsReport", mock.Anything, 1).Return(sampleReport, nil).Once()
			},
			status: http.StatusOK,
		},
		{
			name:      "invalid_id_param",
			url:       "/productRecords/report?id=abc", // non-numeric
			mockSetup: func(_ *productrecordmock.MockProductRecordService) {},
			status:    http.StatusBadRequest,
			appErr:    apperrors.CodeBadRequest,
		},
		{
			name: "service_not_found",
			url:  "/productRecords/report?id=99",
			mockSetup: func(m *productrecordmock.MockProductRecordService) {
				m.On("GetRecordsReport", mock.Anything, 99).
					Return([]models.ProductRecordReport{},
						apperrors.NewAppError(apperrors.CodeNotFound, "nf")).Once()
			},
			status: http.StatusNotFound,
			appErr: apperrors.CodeNotFound,
		},
		{
			name:      "negative_id_param",
			url:       "/productRecords/report?id=-1",
			mockSetup: func(_ *productrecordmock.MockProductRecordService) {},
			status:    http.StatusBadRequest,
			appErr:    apperrors.CodeBadRequest,
		},
		{
			name: "service_internal_error",
			url:  "/productRecords/report",
			mockSetup: func(m *productrecordmock.MockProductRecordService) {
				m.On("GetRecordsReport", mock.Anything, 0).
					Return([]models.ProductRecordReport{}, apperrors.NewAppError(apperrors.CodeInternal, "db down")).Once()
			},
			status: http.StatusInternalServerError,
			appErr: apperrors.CodeInternal,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			svc := &productrecordmock.MockProductRecordService{}
			tc.mockSetup(svc)
			h := handler.NewProductRecordHandler(svc)

			// act
			rec := testhelpers.DoRequest(t, http.MethodGet, tc.url, nil, http.HandlerFunc(h.GetRecordsReport))

			// assert
			require.Equal(t, tc.status, rec.Code)
			if tc.appErr != "" {
				app, derr := testhelpers.DecodeAppErr(rec.Body)
				require.NoError(t, derr)
				testhelpers.RequireAppErr(t, app, tc.appErr)
			}
			svc.AssertExpectations(t)
		})
	}
}
