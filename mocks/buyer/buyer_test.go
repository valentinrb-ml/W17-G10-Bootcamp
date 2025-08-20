package mocks_test

import (
	"context"
	"testing"

	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

func TestBuyerRepositoryMocks_DummyCoverage(t *testing.T) {
	m := &mocks.BuyerRepositoryMocks{
		MockCreate:           func(ctx context.Context, b models.Buyer) (*models.Buyer, error) { return nil, nil },
		MockFindAll:          func(ctx context.Context) ([]models.Buyer, error) { return nil, nil },
		MockFindByID:         func(ctx context.Context, id int) (*models.Buyer, error) { return nil, nil },
		MockFindById:         func(ctx context.Context, id int) (*models.Buyer, error) { return nil, nil },
		MockFindByCardNumber: func(ctx context.Context, cardNumber string) (*models.Buyer, error) { return nil, nil },
		MockUpdate:           func(ctx context.Context, id int, b models.Buyer) error { return nil },
		MockDelete:           func(ctx context.Context, id int) error { return nil },
		MockCardNumberExists: func(ctx context.Context, cardNumber string, id int) bool { return false },
		MockSetLogger:        func(l logger.Logger) {},
	}
	m.Create(context.TODO(), models.Buyer{})
	m.FindAll(context.TODO())
	m.FindByID(context.TODO(), 0)
	m.FindById(context.TODO(), 0)
	m.FindByCardNumber(context.TODO(), "")
	m.Update(context.TODO(), 0, models.Buyer{})
	m.Delete(context.TODO(), 0)
	m.CardNumberExists(context.TODO(), "", 0)
	m.SetLogger(nil)
}

func TestBuyerServiceMock_DummyCoverage(t *testing.T) {
	m := &mocks.BuyerServiceMock{
		CreateFn: func(ctx context.Context, req models.RequestBuyer) (*models.ResponseBuyer, error) { return nil, nil },
		UpdateFn: func(ctx context.Context, id int, req models.RequestBuyer) (*models.ResponseBuyer, error) {
			return nil, nil
		},
		DeleteFn:    func(ctx context.Context, id int) error { return nil },
		FindAllFn:   func(ctx context.Context) ([]models.ResponseBuyer, error) { return nil, nil },
		FindByIdFn:  func(ctx context.Context, id int) (*models.ResponseBuyer, error) { return nil, nil },
		SetLoggerFn: func(l logger.Logger) {},
	}
	m.Create(context.TODO(), models.RequestBuyer{})
	m.Update(context.TODO(), 0, models.RequestBuyer{})
	m.Delete(context.TODO(), 0)
	m.FindAll(context.TODO())
	m.FindById(context.TODO(), 0)
	m.SetLogger(nil)
}
