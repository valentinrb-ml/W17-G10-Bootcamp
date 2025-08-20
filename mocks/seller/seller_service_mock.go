package mocks

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

// SellerServiceMock implements SellerService, but allows customizing behavior in tests.
type SellerServiceMock struct {
	CreateFn    func(ctx context.Context, req models.RequestSeller) (*models.ResponseSeller, error)
	UpdateFn    func(ctx context.Context, id int, req models.RequestSeller) (*models.ResponseSeller, error)
	DeleteFn    func(ctx context.Context, id int) error
	FindAllFn   func(ctx context.Context) ([]models.ResponseSeller, error)
	FindByIdFn  func(ctx context.Context, id int) (*models.ResponseSeller, error)
	SetLoggerFn func(l logger.Logger)
}

func (m *SellerServiceMock) Create(ctx context.Context, req models.RequestSeller) (*models.ResponseSeller, error) {
	return m.CreateFn(ctx, req)
}

func (m *SellerServiceMock) Update(ctx context.Context, id int, req models.RequestSeller) (*models.ResponseSeller, error) {
	return m.UpdateFn(ctx, id, req)
}

func (m *SellerServiceMock) Delete(ctx context.Context, id int) error {
	return m.DeleteFn(ctx, id)
}

func (m *SellerServiceMock) FindAll(ctx context.Context) ([]models.ResponseSeller, error) {
	return m.FindAllFn(ctx)
}

func (m *SellerServiceMock) FindById(ctx context.Context, id int) (*models.ResponseSeller, error) {
	return m.FindByIdFn(ctx, id)
}

func (m *SellerServiceMock) SetLogger(l logger.Logger) {
	if m.SetLoggerFn != nil {
		m.SetLoggerFn(l)
	}
}
