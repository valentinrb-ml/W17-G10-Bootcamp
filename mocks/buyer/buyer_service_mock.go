package mocks

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

// BuyerServiceMock implements BuyerService for testing
type BuyerServiceMock struct {
	CreateFn   func(ctx context.Context, req models.RequestBuyer) (*models.ResponseBuyer, error)
	UpdateFn   func(ctx context.Context, id int, req models.RequestBuyer) (*models.ResponseBuyer, error)
	DeleteFn   func(ctx context.Context, id int) error
	FindAllFn  func(ctx context.Context) ([]models.ResponseBuyer, error)
	FindByIdFn func(ctx context.Context, id int) (*models.ResponseBuyer, error)
}

func (m *BuyerServiceMock) Create(ctx context.Context, req models.RequestBuyer) (*models.ResponseBuyer, error) {
	return m.CreateFn(ctx, req)
}

func (m *BuyerServiceMock) Update(ctx context.Context, id int, req models.RequestBuyer) (*models.ResponseBuyer, error) {
	return m.UpdateFn(ctx, id, req)
}

func (m *BuyerServiceMock) Delete(ctx context.Context, id int) error {
	return m.DeleteFn(ctx, id)
}

func (m *BuyerServiceMock) FindAll(ctx context.Context) ([]models.ResponseBuyer, error) {
	return m.FindAllFn(ctx)
}

func (m *BuyerServiceMock) FindById(ctx context.Context, id int) (*models.ResponseBuyer, error) {
	return m.FindByIdFn(ctx, id)
}
