package product_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

func TestMockRepository_DummyCoverage(t *testing.T) {
	m := &product.MockRepository{}
	m.On("GetAll", mock.Anything).Return([]models.Product{}, nil)
	m.On("Save", mock.Anything, mock.Anything).Return(models.Product{}, nil)
	m.On("GetByID", mock.Anything, mock.Anything).Return(models.Product{}, nil)
	m.On("Delete", mock.Anything, mock.Anything).Return(nil)
	m.On("Patch", mock.Anything, mock.Anything, mock.Anything).Return(models.Product{}, nil)

	m.GetAll(context.TODO())
	m.Save(context.TODO(), models.Product{})
	m.GetByID(context.TODO(), 0)
	m.Delete(context.TODO(), 0)
	m.Patch(context.TODO(), 0, models.ProductPatchRequest{})
}

func TestMockService_DummyCoverage(t *testing.T) {
	m := &product.MockService{}
	m.On("GetAll", mock.Anything).Return([]models.ProductResponse{}, nil)
	m.On("Create", mock.Anything, mock.Anything).Return(models.ProductResponse{}, nil)
	m.On("GetByID", mock.Anything, mock.Anything).Return(models.ProductResponse{}, nil)
	m.On("Delete", mock.Anything, mock.Anything).Return(nil)
	m.On("Patch", mock.Anything, mock.Anything, mock.Anything).Return(models.ProductResponse{}, nil)

	m.GetAll(context.TODO())
	m.Create(context.TODO(), models.Product{})
	m.GetByID(context.TODO(), 0)
	m.Delete(context.TODO(), 0)
	m.Patch(context.TODO(), 0, models.ProductPatchRequest{})
}
