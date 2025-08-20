package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/purchase_order"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/purchase_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestPurchaseOrderService_Create(t *testing.T) {
	t.Run("Successfully create purchase order", func(t *testing.T) {
		// Setup
		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncCreate: func(ctx context.Context, po models.PurchaseOrder) (*models.PurchaseOrder, error) {
				return &po, nil
			},
			FuncExistsOrderNumber: func(ctx context.Context, orderNumber string) bool {
				return false
			},
		}
		service := service.NewPurchaseOrderService(repoMock)

		req := models.RequestPurchaseOrder{
			OrderNumber:     "PO-001",
			OrderDate:       "2023-01-01",
			TrackingCode:    "TRACK001",
			BuyerID:         101,
			ProductRecordID: 201,
		}

		// Execute
		result, err := service.Create(context.Background(), req)

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, req.OrderNumber, result.OrderNumber)
		assert.Equal(t, req.BuyerID, result.BuyerID)
	})

	t.Run("Fail when order number already exists", func(t *testing.T) {
		// Setup
		repoMock := &mocks.PurchaseOrderRepositoryMock{
			FuncExistsOrderNumber: func(ctx context.Context, orderNumber string) bool {
				return true
			},
			// Necesitamos implementar Create aunque falle, para evitar el nil pointer
			FuncCreate: func(ctx context.Context, po models.PurchaseOrder) (*models.PurchaseOrder, error) {
				return nil, apperrors.NewAppError(apperrors.CodeConflict, "order number already exists")
			},
		}
		service := service.NewPurchaseOrderService(repoMock)

		req := models.RequestPurchaseOrder{
			OrderNumber:     "PO-001",
			OrderDate:       "2023-01-01",
			TrackingCode:    "TRACK001",
			BuyerID:         101,
			ProductRecordID: 201,
		}

		// Execute
		_, err := service.Create(context.Background(), req)

		// Verify
		assert.Error(t, err)
		var appErr *apperrors.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperrors.CodeConflict, appErr.Code)
	})

	t.Run("Fail when date is invalid", func(t *testing.T) {
		// Setup
		repoMock := &mocks.PurchaseOrderRepositoryMock{}
		service := service.NewPurchaseOrderService(repoMock)

		req := models.RequestPurchaseOrder{
			OrderNumber:     "PO-001",
			OrderDate:       "invalid-date",
			TrackingCode:    "TRACK001",
			BuyerID:         101,
			ProductRecordID: 201,
		}

		// Execute
		_, err := service.Create(context.Background(), req)

		// Verify
		assert.Error(t, err)
		var appErr *apperrors.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperrors.CodeValidationError, appErr.Code)
	})

	t.Run("Fail when date is in the future", func(t *testing.T) {
		// Setup
		repoMock := &mocks.PurchaseOrderRepositoryMock{}
		service := service.NewPurchaseOrderService(repoMock)
		service.SetLogger(testhelpers.NewTestLogger())

		futureDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		req := models.RequestPurchaseOrder{
			OrderNumber:     "PO-001",
			OrderDate:       futureDate,
			TrackingCode:    "TRACK001",
			BuyerID:         101,
			ProductRecordID: 201,
		}

		// Execute
		_, err := service.Create(context.Background(), req)

		// Verify
		assert.Error(t, err)
		var appErr *apperrors.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperrors.CodeValidationError, appErr.Code)
	})
}
