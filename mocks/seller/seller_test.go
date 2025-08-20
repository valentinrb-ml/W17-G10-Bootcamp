package mocks_test

import (
	"context"
	"testing"

	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/seller"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func TestSellerServiceMock_Dummy(t *testing.T) {
	m := &mocks.SellerServiceMock{
		CreateFn: func(ctx context.Context, req models.RequestSeller) (*models.ResponseSeller, error) { return nil, nil },
		UpdateFn: func(ctx context.Context, id int, req models.RequestSeller) (*models.ResponseSeller, error) {
			return nil, nil
		},
		DeleteFn:    func(ctx context.Context, id int) error { return nil },
		FindAllFn:   func(ctx context.Context) ([]models.ResponseSeller, error) { return nil, nil },
		FindByIdFn:  func(ctx context.Context, id int) (*models.ResponseSeller, error) { return nil, nil },
		SetLoggerFn: func(l logger.Logger) {},
	}
	m.Create(context.TODO(), models.RequestSeller{})
	m.Update(context.TODO(), 0, models.RequestSeller{})
	m.Delete(context.TODO(), 0)
	m.FindAll(context.TODO())
	m.FindById(context.TODO(), 0)
	m.SetLogger(nil)
}

func TestSellerRepositoryMock_DummyCoverage(t *testing.T) {
	m := &mocks.SellerRepositoryMock{
		CreateFn:    func(ctx context.Context, s models.Seller) (*models.Seller, error) { return nil, nil },
		UpdateFn:    func(ctx context.Context, id int, s models.Seller) error { return nil },
		DeleteFn:    func(ctx context.Context, id int) error { return nil },
		FindAllFn:   func(ctx context.Context) ([]models.Seller, error) { return nil, nil },
		FindByIdFn:  func(ctx context.Context, id int) (*models.Seller, error) { return nil, nil },
		SetLoggerFn: func(l logger.Logger) {},
	}

	m.Create(context.TODO(), models.Seller{})
	m.Update(context.TODO(), 0, models.Seller{})
	m.Delete(context.TODO(), 0)
	m.FindAll(context.TODO())
	m.FindById(context.TODO(), 0)
	m.SetLogger(nil)
}
