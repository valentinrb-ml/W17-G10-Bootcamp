package mocks

import (
	"context"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

type ProductBatchServiceMock struct {
	FuncCreate        func(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error)
	FuncGetReportById func(ctx context.Context, id int) (*models.ReportProduct, error)
	FuncGetReport     func(ctx context.Context) ([]models.ReportProduct, error)
}

func (m *ProductBatchServiceMock) CreateProductBatches(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error) {
	return m.FuncCreate(ctx, proBa)
}
func (m *ProductBatchServiceMock) GetReportProductById(ctx context.Context, id int) (*models.ReportProduct, error) {

	return m.FuncGetReportById(ctx, id)
}
func (m *ProductBatchServiceMock) GetReportProduct(ctx context.Context) ([]models.ReportProduct, error) {
	return m.FuncGetReport(ctx)
}
