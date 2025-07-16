package service

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]product.ProductResponse, error)
	Create(ctx context.Context, prod product.Product) (product.ProductResponse, error)
	GetByID(ctx context.Context, id int) (product.ProductResponse, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, req product.ProductPatchRequest) (product.ProductResponse, error)
}
