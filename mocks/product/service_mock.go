package product

import (
	"context"

	"github.com/stretchr/testify/mock"
	model "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type MockService struct{ mock.Mock }

func (m *MockService) GetAll(ctx context.Context) ([]model.ProductResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.ProductResponse), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, p model.Product) (model.ProductResponse, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(model.ProductResponse), args.Error(1)
}

func (m *MockService) GetByID(ctx context.Context, id int) (model.ProductResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.ProductResponse), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockService) Patch(ctx context.Context, id int, req model.ProductPatchRequest) (model.ProductResponse, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(model.ProductResponse), args.Error(1)
}
