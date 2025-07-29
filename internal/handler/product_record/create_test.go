package handler_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_record"
	mappers "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers/product_record"
	productrecordmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

/*
Test for POST /productRecords handler.

Scenarios:
  - 201 Created  (happy path)
  - 409 Conflict  (service error)
  - 400 BadRequest    – JSON decode failure
  - 422 Validation    – business rule error
*/
func TestProductRecordHandler_Create(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	type want struct {
		status int
		appErr string
	}

	tests := []struct {
		name      string
		reqBody   any
		mockSetup func(*productrecordmock.MockProductRecordService, models.ProductRecord)
		want      want
	}{
		{
			name:    "success",
			reqBody: testhelpers.BuildProductRecordRequest(fixedTime, 10.2, 15.5, 3),
			mockSetup: func(m *productrecordmock.MockProductRecordService, r models.ProductRecord) {
				m.On("Create", mock.Anything, r).Return(r, nil).Once()
			},
			want: want{status: http.StatusCreated},
		},
		{
			name:    "service_conflict",
			reqBody: testhelpers.BuildProductRecordRequest(fixedTime, 10.2, 15.5, 99),
			mockSetup: func(m *productrecordmock.MockProductRecordService, r models.ProductRecord) {
				m.On("Create", mock.Anything, r).
					Return(models.ProductRecord{},
						apperrors.NewAppError(apperrors.CodeConflict, "fk")).Once()
			},
			want: want{status: http.StatusConflict, appErr: apperrors.CodeConflict},
		},
		{
			name:      "decode_error",
			reqBody:   "{{", // malformed JSON
			mockSetup: func(_ *productrecordmock.MockProductRecordService, _ models.ProductRecord) {},
			want:      want{status: http.StatusBadRequest, appErr: apperrors.CodeBadRequest},
		},
		{
			name:      "validation_error",
			reqBody:   testhelpers.BuildProductRecordRequest(fixedTime, -5, 0, 0),
			mockSetup: func(_ *productrecordmock.MockProductRecordService, _ models.ProductRecord) {},
			want:      want{status: http.StatusUnprocessableEntity, appErr: apperrors.CodeValidationError},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := &productrecordmock.MockProductRecordService{}

			// build domain model only if payload is valid JSON
			var domain models.ProductRecord
			if req, ok := tc.reqBody.(models.ProductRecordRequest); ok {
				domain = mappers.ProductRecordRequestToDomain(req)
			}
			tc.mockSetup(svc, domain)

			h := handler.NewProductRecordHandler(svc)
			rec := testhelpers.DoRequest(t, http.MethodPost, "/productRecords", tc.reqBody, h.Create)

			require.Equal(t, tc.want.status, rec.Code)
			if tc.want.appErr != "" {
				app, derr := testhelpers.DecodeAppErr(rec.Body)
				require.NoError(t, derr)
				testhelpers.RequireAppErr(t, app, tc.want.appErr)
			}
			svc.AssertExpectations(t)
		})
	}
}
