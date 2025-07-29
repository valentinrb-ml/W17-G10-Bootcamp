package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product"
	productmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
)

func TestProductService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mockFn  func(*productmock.MockRepository)
		wantErr bool
	}{
		{
			name: "success",
			mockFn: func(r *productmock.MockRepository) {
				r.On("Delete", mock.Anything, 8).Return(nil).Once()
			},
		},
		{
			name:    "not found",
			wantErr: true,
			mockFn: func(r *productmock.MockRepository) {
				r.On("Delete", mock.Anything, 8).
					Return(apperrors.NewAppError(apperrors.CodeNotFound, "x")).Once()
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
			err := svc.Delete(context.Background(), 8)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}
