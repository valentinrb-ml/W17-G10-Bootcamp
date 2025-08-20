package product

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"

	"github.com/stretchr/testify/mock"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type MockRepository struct {
	mock.Mock
	log logger.Logger
}

// interface ProductRepository
func (m *MockRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Product), args.Error(1)
}
func (m *MockRepository) Save(ctx context.Context, p models.Product) (models.Product, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(models.Product), args.Error(1)
}
func (m *MockRepository) GetByID(ctx context.Context, id int) (models.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Product), args.Error(1)
}
func (m *MockRepository) Delete(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockRepository) Patch(ctx context.Context, id int, req models.ProductPatchRequest) (models.Product, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockRepository) SetLogger(l logger.Logger) {
	m.log = l
}
