package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product"
	productmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// GetAll() – happy path and error wrapped.
func TestProductService_GetAll(t *testing.T) {
	t.Parallel()

	happy := []models.Product{testhelpers.BuildProduct(1)}

	tests := []struct {
		name    string
		mockFn  func(*productmock.MockRepository)
		wantErr bool
		appCode string
	}{
		{
			name: "success",
			mockFn: func(r *productmock.MockRepository) {
				r.On("GetAll", mock.Anything).Return(happy, nil).Once()
			},
		},
		{
			name:    "repo error wrapped",
			wantErr: true,
			appCode: apperrors.CodeInternal,
			mockFn: func(r *productmock.MockRepository) {
				// return a null slice type so as not to break the type-assert
				var nilSlice []models.Product
				r.On("GetAll", mock.Anything).Return(nilSlice, errors.New("db down")).Once()
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &productmock.MockRepository{}
			tc.mockFn(repo)

			svc := service.NewProductService(repo)
			_, err := svc.GetAll(context.Background())

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
