package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product_record"
	productrecordmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// fixture returned by the “success” mock
var sampleReport = []models.ProductRecordReport{
	{ProductID: 1, Description: "p1", RecordsCount: 3},
}

func TestProductRecordService_GetRecordsReport(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        int // input param
		mockSetup func(r *productrecordmock.MockProductRecordRepository)
		wantErr   string // expected AppError.Code ("" = no error)
	}{
		{
			name: "success_all",
			id:   0,
			mockSetup: func(r *productrecordmock.MockProductRecordRepository) {
				r.On("GetRecordsReport", mock.Anything, 0).Return(sampleReport, nil).Once()
			},
		},
		{
			name: "repo_internal_error",
			id:   0,
			mockSetup: func(r *productrecordmock.MockProductRecordRepository) {
				// repo fails → service must propagate INTERNAL_ERROR
				r.On("GetRecordsReport", mock.Anything, 0).
					Return([]models.ProductRecordReport{}, apperrors.NewAppError(apperrors.CodeInternal, "db")).Once()
			},
			wantErr: apperrors.CodeInternal,
		},
	}

	for _, tc := range tests {
		tc := tc // capture range var
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// 1. mock repository
			repoMock := &productrecordmock.MockProductRecordRepository{}
			tc.mockSetup(repoMock)

			// 2. service under test
			svc := service.NewProductRecordService(repoMock)
			svc.SetLogger(testhelpers.NewTestLogger())

			// 3. execute
			_, err := svc.GetRecordsReport(context.Background(), tc.id)

			// 4. assertions
			if tc.wantErr != "" {
				testhelpers.RequireAppErr(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			repoMock.AssertExpectations(t) // all mocked calls were hit
		})
	}
}
