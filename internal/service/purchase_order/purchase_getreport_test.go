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

func TestPurchaseOrderService_GetReportByBuyer(t *testing.T) {
	t.Run("Successfully get report for specific buyer", func(t *testing.T) {
		// Setup
		expectedReport := []models.BuyerWithPurchaseCount{
			testhelpers.BuyerWithPurchaseCountDummyMap[101],
		}

		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncGetCountByBuyer: func(ctx context.Context, buyerID int) ([]models.BuyerWithPurchaseCount, error) {
				return expectedReport, nil
			},
		}
		service := service.NewPurchaseOrderService(repoMock)

		buyerID := 101

		// Execute
		result, err := service.GetReportByBuyer(context.Background(), &buyerID)

		// Verify
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, buyerID, result[0].ID)
		assert.Equal(t, "John", result[0].FirstName)
	})

	t.Run("Successfully get report for all buyers", func(t *testing.T) {
		// Setup
		expectedReport := []models.BuyerWithPurchaseCount{
			testhelpers.BuyerWithPurchaseCountDummyMap[101],
			testhelpers.BuyerWithPurchaseCountDummyMap[102],
		}

		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncGetAllWithPurchaseCount: func(ctx context.Context) ([]models.BuyerWithPurchaseCount, error) {
				return expectedReport, nil
			},
		}
		service := service.NewPurchaseOrderService(repoMock)

		// Execute
		result, err := service.GetReportByBuyer(context.Background(), nil)

		// Verify
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "John", result[0].FirstName)
		assert.Equal(t, "Jane", result[1].FirstName)
	})

	t.Run("Return error when repository fails for specific buyer", func(t *testing.T) {
		// Setup
		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncGetCountByBuyer: func(ctx context.Context, buyerID int) ([]models.BuyerWithPurchaseCount, error) {
				return nil, errors.New("repository error")
			},
		}
		service := service.NewPurchaseOrderService(repoMock)

		buyerID := 101

		// Execute
		_, err := service.GetReportByBuyer(context.Background(), &buyerID)

		// Verify
		assert.Error(t, err)
	})

	t.Run("Return error when repository fails for all buyers", func(t *testing.T) {
		// Setup
		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncGetAllWithPurchaseCount: func(ctx context.Context) ([]models.BuyerWithPurchaseCount, error) {
				return nil, errors.New("repository error")
			},
		}
		service := service.NewPurchaseOrderService(repoMock)

		// Execute
		_, err := service.GetReportByBuyer(context.Background(), nil)

		// Verify
		assert.Error(t, err)
	})

	t.Run("Return not found when buyer has no orders", func(t *testing.T) {
		// Setup
		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncGetCountByBuyer: func(ctx context.Context, buyerID int) ([]models.BuyerWithPurchaseCount, error) {
				return nil, apperrors.NewAppError(apperrors.CodeNotFound, "buyer has no orders")
			},
		}
		service := service.NewPurchaseOrderService(repoMock)
		service.SetLogger(testhelpers.NewTestLogger())

		buyerID := 999

		// Execute
		_, err := service.GetReportByBuyer(context.Background(), &buyerID)

		// Verify
		assert.Error(t, err)
		var appErr *apperrors.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperrors.CodeNotFound, appErr.Code)
	})
}
