package mocks

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type PurchaseOrderServiceMock struct {
	CreateFn                  func(ctx context.Context, po models.PurchaseOrder) (*models.PurchaseOrder, error)
	GetAllFn                  func(ctx context.Context) ([]models.PurchaseOrder, error)
	GetByIDFn                 func(ctx context.Context, id int) (*models.PurchaseOrder, error)
	GetCountByBuyerFn         func(ctx context.Context, buyerID int) ([]models.BuyerWithPurchaseCount, error)
	GetAllWithPurchaseCountFn func(ctx context.Context) ([]models.BuyerWithPurchaseCount, error)
}

func (m *PurchaseOrderServiceMock) Create(ctx context.Context, po models.PurchaseOrder) (*models.PurchaseOrder, error) {
	return m.CreateFn(ctx, po)
}

func (m *PurchaseOrderServiceMock) GetAll(ctx context.Context) ([]models.PurchaseOrder, error) {
	return m.GetAllFn(ctx)
}

func (m *PurchaseOrderServiceMock) GetByID(ctx context.Context, id int) (*models.PurchaseOrder, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *PurchaseOrderServiceMock) GetCountByBuyer(ctx context.Context, buyerID int) ([]models.BuyerWithPurchaseCount, error) {
	return m.GetCountByBuyerFn(ctx, buyerID)
}

func (m *PurchaseOrderServiceMock) GetAllWithPurchaseCount(ctx context.Context) ([]models.BuyerWithPurchaseCount, error) {
	return m.GetAllWithPurchaseCountFn(ctx)
}
