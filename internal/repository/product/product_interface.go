package repository

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (models.Product, error)
	Save(ctx context.Context, p models.Product) (models.Product, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, req models.ProductPatchRequest) (models.Product, error)
	SetLogger(l logger.Logger)
}
