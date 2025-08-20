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

func TestProductService_GetByID(t *testing.T) {
	t.Parallel()

	okProd := testhelpers.BuildProduct(42)

	tests := []struct {
		name    string
		mockFn  func(*productmock.MockRepository)
		wantErr bool
	}{
		{
			name: "success",
			mockFn: func(r *productmock.MockRepository) {
				r.On("GetByID", mock.Anything, 42).Return(okProd, nil).Once()
			},
		},
		{
			name:    "not found",
			wantErr: true,
			mockFn: func(r *productmock.MockRepository) {
				r.On("GetByID", mock.Anything, 42).
					Return(models.Product{}, apperrors.NewAppError(apperrors.CodeNotFound, "x")).Once()
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
			svc.SetLogger(testhelpers.NewTestLogger())

			_, err := svc.GetByID(context.Background(), 42)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}
