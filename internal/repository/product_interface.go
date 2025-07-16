package repository

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]product.Product, error)
	GetByID(ctx context.Context, id int) (product.Product, error)
	Save(ctx context.Context, p product.Product) (product.Product, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, req product.ProductPatchRequest) (product.Product, error)
	create(ctx context.Context, p product.Product) (product.Product, error)
	update(ctx context.Context, p product.Product) (product.Product, error)
	handleDBError(err error, msg string) error
}
