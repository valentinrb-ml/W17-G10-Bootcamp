package service

import (
	"context"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

// CreateProductBatches creates a new product batch using the repository.
// Delegates creation to the repository layer and returns the created batch.
func (s *productBatchesService) CreateProductBatches(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error) {
	newProBa, err := s.r.CreateProductBatches(ctx, proBa)
	if err != nil {
		return nil, err
	}
	return newProBa, nil
}

// GetReportProductById retrieves a report for products in a section by its number.
// Calls the repository to fetch the report for a specific section.
func (s *productBatchesService) GetReportProductById(ctx context.Context, sectionNumber int) (*models.ReportProduct, error) {
	reportProduct, err := s.r.GetReportProductById(ctx, sectionNumber)
	if err != nil {
		return nil, err
	}
	return reportProduct, nil
}

// GetReportProduct gets product report data for all sections.
// Uses the repository to retrieve aggregated product batch information.
func (s *productBatchesService) GetReportProduct(ctx context.Context) ([]models.ReportProduct, error) {
	reportsProduct, err := s.r.GetReportProduct(ctx)
	if err != nil {
		return nil, err
	}
	return reportsProduct, nil
}
