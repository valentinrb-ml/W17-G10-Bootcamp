package mocks

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

// SellerRepositoryMock is a manual mock of SellerRepository (happy path only).
type SellerRepositoryMock struct {
	CreateFn    func(ctx context.Context, s models.Seller) (*models.Seller, error)
	UpdateFn    func(ctx context.Context, id int, s models.Seller) error
	DeleteFn    func(ctx context.Context, id int) error
	FindAllFn   func(ctx context.Context) ([]models.Seller, error)
	FindByIdFn  func(ctx context.Context, id int) (*models.Seller, error)
	SetLoggerFn func(l logger.Logger)
}

func (m *SellerRepositoryMock) Create(ctx context.Context, s models.Seller) (*models.Seller, error) {
	return m.CreateFn(ctx, s)
}

func (m *SellerRepositoryMock) Update(ctx context.Context, id int, s models.Seller) error {
	return m.UpdateFn(ctx, id, s)
}

func (m *SellerRepositoryMock) Delete(ctx context.Context, id int) error {
	return m.DeleteFn(ctx, id)
}

func (m *SellerRepositoryMock) FindAll(ctx context.Context) ([]models.Seller, error) {
	return m.FindAllFn(ctx)
}

func (m *SellerRepositoryMock) FindById(ctx context.Context, id int) (*models.Seller, error) {
	return m.FindByIdFn(ctx, id)
}

func (m *SellerRepositoryMock) SetLogger(l logger.Logger) {
	if m.SetLoggerFn != nil {
		m.SetLoggerFn(l)
	}
}
