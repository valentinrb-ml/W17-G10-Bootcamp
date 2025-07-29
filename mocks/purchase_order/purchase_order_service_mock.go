package mocks

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

// Mock del service.PurchaseOrderService
type PurchaseOrderServiceMock struct {
	CreateFn           func(ctx context.Context, req models.RequestPurchaseOrder) (*models.ResponsePurchaseOrder, error)
	GetAllFn           func(ctx context.Context) ([]models.ResponsePurchaseOrder, error)
	GetByIDFn          func(ctx context.Context, id int) (*models.ResponsePurchaseOrder, error)
	GetReportByBuyerFn func(ctx context.Context, buyerID *int) ([]models.BuyerWithPurchaseCount, error)
}

func (m *PurchaseOrderServiceMock) Create(ctx context.Context, req models.RequestPurchaseOrder) (*models.ResponsePurchaseOrder, error) {
	return m.CreateFn(ctx, req)
}

func (m *PurchaseOrderServiceMock) GetAll(ctx context.Context) ([]models.ResponsePurchaseOrder, error) {
	return m.GetAllFn(ctx)
}

func (m *PurchaseOrderServiceMock) GetByID(ctx context.Context, id int) (*models.ResponsePurchaseOrder, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *PurchaseOrderServiceMock) GetReportByBuyer(ctx context.Context, buyerID *int) ([]models.BuyerWithPurchaseCount, error) {
	return m.GetReportByBuyerFn(ctx, buyerID)
}
