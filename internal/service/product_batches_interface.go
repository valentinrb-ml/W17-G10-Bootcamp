package service

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

// ProductBatchesService defines business logic for product batches.
type ProductBatchesService interface {
	CreateProductBatches(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error)
	GetReportProductById(ctx context.Context, sectionNumber int) (*models.ReportProduct, error)
	GetReportProduct(ctx context.Context) ([]models.ReportProduct, error)
}

// productBatchesService implements ProductBatchesService using a repository.
type productBatchesService struct {
	r repository.ProductBatchesRepository
}

// NewProductBatchesService creates a new ProductBatchesService using the provided repository.
func NewProductBatchesService(repo repository.ProductBatchesRepository) ProductBatchesService {
	return &productBatchesService{
		repo,
	}
}
