package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product"
	productmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// invalidProduct builds a product that surely breaks business rules:
//   - negative dimensions
//   - nil SellerID
//   - ProductType == 0
func invalidProduct() models.Product {
	p := testhelpers.BuildProduct(0)
	p.Dimensions.Width = -1
	p.Dimensions.Height = -1
	p.Dimensions.Length = -1
	p.SellerID = nil
	p.ProductType = 0
	return p
}

// Create() – success, repository error, validation error.
func TestProductService_Create(t *testing.T) {
	t.Parallel()

	saved := testhelpers.BuildProduct(10)

	tests := []struct {
		name    string
		input   models.Product
		mockFn  func(*productmock.MockRepository)
		wantErr bool
		appCode string
	}{
		{
			name:  "success",
			input: testhelpers.BuildProduct(0),
			mockFn: func(r *productmock.MockRepository) {
				r.On("Save", mock.Anything, testhelpers.BuildProduct(0)).
					Return(saved, nil).Once()
			},
		},
		{
			name:    "repo conflict",
			input:   testhelpers.BuildProduct(0),
			wantErr: true,
			appCode: apperrors.CodeConflict,
			mockFn: func(r *productmock.MockRepository) {
				r.On("Save", mock.Anything, mock.Anything).
					Return(models.Product{}, apperrors.NewAppError(apperrors.CodeConflict, "dup")).Once()
			},
		},
		{
			name:    "validation error",
			input:   invalidProduct(),
			wantErr: true,
			appCode: apperrors.CodeValidationError,
			mockFn:  func(r *productmock.MockRepository) {}, // repo should not be called
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &productmock.MockRepository{}
			tc.mockFn(repo)

			svc := service.NewProductService(repo)
			_, err := svc.Create(context.Background(), tc.input)

			if tc.wantErr {
				require.Error(t, err)
				testhelpers.RequireAppErr(t, err, tc.appCode)
			} else {
				require.NoError(t, err)
			}
			
			repo.AssertExpectations(t)
		})
	}
}
