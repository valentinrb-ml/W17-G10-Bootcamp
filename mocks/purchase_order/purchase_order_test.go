package mocks_test

import (
	"context"
	"testing"

	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/purchase_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

func TestPurchaseOrderServiceMock_DummyCoverage(t *testing.T) {
	m := &mocks.PurchaseOrderServiceMock{
		CreateFn: func(ctx context.Context, req models.RequestPurchaseOrder) (*models.ResponsePurchaseOrder, error) {
			return nil, nil
		},
		GetAllFn:           func(ctx context.Context) ([]models.ResponsePurchaseOrder, error) { return nil, nil },
		GetByIDFn:          func(ctx context.Context, id int) (*models.ResponsePurchaseOrder, error) { return nil, nil },
		GetReportByBuyerFn: func(ctx context.Context, buyerID *int) ([]models.BuyerWithPurchaseCount, error) { return nil, nil },
		SetLoggerFn:        func(l logger.Logger) {},
	}

	m.Create(context.TODO(), models.RequestPurchaseOrder{})
	m.GetAll(context.TODO())
	m.GetByID(context.TODO(), 0)
	m.GetReportByBuyer(context.TODO(), nil)
	m.SetLogger(nil)
}

func TestPurchaseOrderRepositoryMock_DummyCoverage(t *testing.T) {
	m := &mocks.PurchaseOrderRepositoryMock{
		FuncCreate:                  func(ctx context.Context, po models.PurchaseOrder) (*models.PurchaseOrder, error) { return nil, nil },
		FuncGetAll:                  func(ctx context.Context) ([]models.PurchaseOrder, error) { return nil, nil },
		FuncGetByID:                 func(ctx context.Context, id int) (*models.PurchaseOrder, error) { return nil, nil },
		FuncExistsOrderNumber:       func(ctx context.Context, orderNumber string) bool { return false },
		FuncGetCountByBuyer:         func(ctx context.Context, buyerID int) ([]models.BuyerWithPurchaseCount, error) { return nil, nil },
		FuncGetAllWithPurchaseCount: func(ctx context.Context) ([]models.BuyerWithPurchaseCount, error) { return nil, nil },
		FuncSetLogger:               func(l logger.Logger) {},
	}

	m.Create(context.TODO(), models.PurchaseOrder{})
	m.GetAll(context.TODO())
	m.GetByID(context.TODO(), 0)
	m.ExistsOrderNumber(context.TODO(), "")
	m.GetCountByBuyer(context.TODO(), 0)
	m.GetAllWithPurchaseCount(context.TODO())
	m.SetLogger(nil)
}
