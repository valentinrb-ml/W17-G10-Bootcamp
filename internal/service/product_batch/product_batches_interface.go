package service

import (
	"context"

	productBatchRepository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/product_batch"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
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
	r      productBatchRepository.ProductBatchesRepository
	logger logger.Logger
}

// SetLogger allows you to inject the logger after creation
func (s *productBatchesService) SetLogger(l logger.Logger) {
	s.logger = l
}

// NewProductBatchesService creates a new ProductBatchesService using the provided repository.
func NewProductBatchesService(repo productBatchRepository.ProductBatchesRepository) ProductBatchesService {
	return &productBatchesService{
		r: repo,
	}
}
