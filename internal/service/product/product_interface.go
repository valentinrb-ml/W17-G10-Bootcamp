package product

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]models.ProductResponse, error)
	Create(ctx context.Context, prod models.Product) (models.ProductResponse, error)
	GetByID(ctx context.Context, id int) (models.ProductResponse, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, req models.ProductPatchRequest) (models.ProductResponse, error)
}
