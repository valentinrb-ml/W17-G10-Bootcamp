package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product_record"
	productrecordmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

/*
Covered scenarios:
 1. Happy-path  -> repo returns record, service returns record
 2. Business rule violated  -> service stops before hitting repo
 3. Repo returns CONFLICT   -> service propagates the error
*/
func TestProductRecordService_Create(t *testing.T) {
	t.Parallel()

	// baseline valid record used in several cases
	valid := models.ProductRecord{
		ProductRecordCore: models.ProductRecordCore{
			LastUpdateDate: time.Now(),
			PurchasePrice:  10,
			SalePrice:      15,
			ProductID:      1,
		},
	}

	tests := []struct {
		name      string
		input     models.ProductRecord
		mockSetup func(r *productrecordmock.MockProductRecordRepository)
		wantErr   string // expected AppError.Code ("" means success)
	}{
		{
			name:  "success",
			input: valid,
			mockSetup: func(r *productrecordmock.MockProductRecordRepository) {
				// repo happy path: return the same record, no error
				r.On("Create", mock.Anything, valid).Return(valid, nil).Once()
			},
		},
		{
			name: "business_rule_validation",
			input: func() models.ProductRecord {
				rec := valid
				rec.PurchasePrice = -1 // invalid -> service must reject
				return rec
			}(),
			mockSetup: func(_ *productrecordmock.MockProductRecordRepository) {},
			wantErr:   apperrors.CodeBadRequest,
		},
		{
			name:  "repository_conflict_error",
			input: valid,
			mockSetup: func(r *productrecordmock.MockProductRecordRepository) {
				// repo signals UNIQUE/FK conflict
				r.On("Create", mock.Anything, valid).
					Return(models.ProductRecord{}, apperrors.NewAppError(apperrors.CodeConflict, "duplicate")).Once()
			},
			wantErr: apperrors.CodeConflict,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repoMock := &productrecordmock.MockProductRecordRepository{}
			tc.mockSetup(repoMock)
			svc := service.NewProductRecordService(repoMock)
			svc.SetLogger(testhelpers.NewTestLogger())

			got, err := svc.Create(context.Background(), tc.input)

			if tc.wantErr != "" {
				testhelpers.RequireAppErr(t, err, tc.wantErr)
				repoMock.AssertExpectations(t)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.input.ProductID, got.ProductID)
			repoMock.AssertExpectations(t)
		})
	}
}
