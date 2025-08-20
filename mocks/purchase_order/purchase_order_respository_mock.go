package mocks

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

// PurchaseOrderRepositoryMock implements PurchaseOrderRepository for testing
type PurchaseOrderRepositoryMock struct {
	FuncCreate                  func(ctx context.Context, po models.PurchaseOrder) (*models.PurchaseOrder, error)
	FuncGetAll                  func(ctx context.Context) ([]models.PurchaseOrder, error)
	FuncGetByID                 func(ctx context.Context, id int) (*models.PurchaseOrder, error)
	FuncExistsOrderNumber       func(ctx context.Context, orderNumber string) bool
	FuncGetCountByBuyer         func(ctx context.Context, buyerID int) ([]models.BuyerWithPurchaseCount, error)
	FuncGetAllWithPurchaseCount func(ctx context.Context) ([]models.BuyerWithPurchaseCount, error)
	FuncSetLogger               func(l logger.Logger)
}

func (m *PurchaseOrderRepositoryMock) Create(ctx context.Context, po models.PurchaseOrder) (*models.PurchaseOrder, error) {
	if m.FuncCreate != nil {
		return m.FuncCreate(ctx, po)
	}
	return nil, nil
}

func (m *PurchaseOrderRepositoryMock) GetAll(ctx context.Context) ([]models.PurchaseOrder, error) {
	if m.FuncGetAll != nil {
		return m.FuncGetAll(ctx)
	}
	return nil, nil
}

func (m *PurchaseOrderRepositoryMock) GetByID(ctx context.Context, id int) (*models.PurchaseOrder, error) {
	if m.FuncGetByID != nil {
		return m.FuncGetByID(ctx, id)
	}
	return nil, nil
}

func (m *PurchaseOrderRepositoryMock) ExistsOrderNumber(ctx context.Context, orderNumber string) bool {
	if m.FuncExistsOrderNumber != nil {
		return m.FuncExistsOrderNumber(ctx, orderNumber)
	}
	return false
}

func (m *PurchaseOrderRepositoryMock) GetCountByBuyer(ctx context.Context, buyerID int) ([]models.BuyerWithPurchaseCount, error) {
	if m.FuncGetCountByBuyer != nil {
		return m.FuncGetCountByBuyer(ctx, buyerID)
	}
	return nil, nil
}

func (m *PurchaseOrderRepositoryMock) GetAllWithPurchaseCount(ctx context.Context) ([]models.BuyerWithPurchaseCount, error) {
	if m.FuncGetAllWithPurchaseCount != nil {
		return m.FuncGetAllWithPurchaseCount(ctx)
	}
	return nil, nil
}

func (m *PurchaseOrderRepositoryMock) SetLogger(l logger.Logger) {
	if m.FuncSetLogger != nil {
		m.FuncSetLogger(l)
	}
}
