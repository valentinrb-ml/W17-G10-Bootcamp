package service

import (
	"context"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

func (s *productBatchesService) CreateProductBatches(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error) {
	newProBa, err := s.r.CreateProductBatches(ctx, proBa)
	if err != nil {
		return nil, err
	}
	return newProBa, nil
}

func (s *productBatchesService) GetReportProductById(ctx context.Context, sectionNumber int) (*models.ReportProduct, error) {
	reportProduct, err := s.r.GetReportProductById(ctx, sectionNumber)
	if err != nil {
		return nil, err
	}
	return reportProduct, nil
}
func (s *productBatchesService) GetReportProduct(ctx context.Context) ([]models.ReportProduct, error) {
	reportsProduct, err := s.r.GetReportProduct(ctx)
	if err != nil {
		return nil, err
	}
	return reportsProduct, nil
}
