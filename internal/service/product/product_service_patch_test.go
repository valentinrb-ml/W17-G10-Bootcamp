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

func TestProductService_Patch(t *testing.T) {
	t.Parallel()

	req := models.ProductPatchRequest{Description: func() *string { s := "upd"; return &s }()}
	updated := testhelpers.BuildProduct(21)

	tests := []struct {
		name    string
		mockFn  func(*productmock.MockRepository)
		wantErr bool
	}{
		{
			name: "success",
			mockFn: func(r *productmock.MockRepository) {
				r.On("Patch", mock.Anything, 21, req).Return(updated, nil).Once()
			},
		},
		{
			name:    "conflict",
			wantErr: true,
			mockFn: func(r *productmock.MockRepository) {
				r.On("Patch", mock.Anything, 21, req).
					Return(models.Product{}, apperrors.NewAppError(apperrors.CodeConflict, "dup")).Once()
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
			_, err := svc.Patch(context.Background(), 21, req)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}
