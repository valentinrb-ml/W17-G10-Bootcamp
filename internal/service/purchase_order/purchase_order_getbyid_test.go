package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/purchase_order"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/purchase_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestPurchaseOrderService_GetByID(t *testing.T) {
	t.Run("Successfully get purchase order by ID", func(t *testing.T) {
		// Setup
		expectedPO := testhelpers.CreateTestPurchaseOrder(1)

		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncGetByID: func(ctx context.Context, id int) (*models.PurchaseOrder, error) {
				return &expectedPO, nil
			},
		}
		service := service.NewPurchaseOrderService(repoMock)

		// Execute
		result, err := service.GetByID(context.Background(), 1)

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, expectedPO.ID, result.ID)
		assert.Equal(t, expectedPO.OrderNumber, result.OrderNumber)
	})

	t.Run("Return not found error when purchase order doesn't exist", func(t *testing.T) {
		// Setup
		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncGetByID: func(ctx context.Context, id int) (*models.PurchaseOrder, error) {
				return nil, apperrors.NewAppError(apperrors.CodeNotFound, "purchase order not found")
			},
		}
		service := service.NewPurchaseOrderService(repoMock)

		// Execute
		_, err := service.GetByID(context.Background(), 999)

		// Verify
		assert.Error(t, err)
		var appErr *apperrors.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperrors.CodeNotFound, appErr.Code)
	})

	t.Run("Return error when repository fails", func(t *testing.T) {
		// Setup
		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncGetByID: func(ctx context.Context, id int) (*models.PurchaseOrder, error) {
				return nil, errors.New("repository error")
			},
		}
		service := service.NewPurchaseOrderService(repoMock)
		service.SetLogger(testhelpers.NewTestLogger())

		// Execute
		_, err := service.GetByID(context.Background(), 1)

		// Verify
		assert.Error(t, err)
	})
}
