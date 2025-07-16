package service

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

type ProductBatchesService interface {
	CreateProductBatches(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error)
	GetReportProductById(ctx context.Context, sectionNumber int) (*models.ReportProduct, error)
	GetReportProduct(ctx context.Context) ([]models.ReportProduct, error)
}

type productBatchesService struct {
	r repository.ProductBatchesRepository
}

func NewProductBatchesService(repo repository.ProductBatchesRepository) ProductBatchesService {
	return &productBatchesService{
		repo,
	}
}
