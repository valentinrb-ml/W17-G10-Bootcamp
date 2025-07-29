package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/purchase_order"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/purchase_order"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestPurchaseOrderService_GetAll(t *testing.T) {
	t.Run("Successfully get all purchase orders", func(t *testing.T) {
		// Setup
		expectedPOs := []models.PurchaseOrder{
			testhelpers.CreateTestPurchaseOrder(1),
			testhelpers.CreateTestPurchaseOrder(2),
		}

		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncGetAll: func(ctx context.Context) ([]models.PurchaseOrder, error) {
				return expectedPOs, nil
			},
		}
		service := service.NewPurchaseOrderService(repoMock)

		// Execute
		result, err := service.GetAll(context.Background())

		// Verify
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "TEST-PO", result[0].OrderNumber)
		assert.Equal(t, 1, result[0].ID)
	})

	t.Run("Return empty slice when no purchase orders exist", func(t *testing.T) {
		// Setup
		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncGetAll: func(ctx context.Context) ([]models.PurchaseOrder, error) {
				return []models.PurchaseOrder{}, nil
			},
		}
		service := service.NewPurchaseOrderService(repoMock)

		// Execute
		result, err := service.GetAll(context.Background())

		// Verify
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Return error when repository fails", func(t *testing.T) {
		// Setup
		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncGetAll: func(ctx context.Context) ([]models.PurchaseOrder, error) {
				return nil, errors.New("repository error")
			},
		}
		service := service.NewPurchaseOrderService(repoMock)

		// Execute
		_, err := service.GetAll(context.Background())

		// Verify
		assert.Error(t, err)
	})
}
