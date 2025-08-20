package mocks_test

import (
	"context"
	"testing"

	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product_batch"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

func TestProductBatchServiceMock_DummyCoverage(t *testing.T) {
	m := &mocks.ProductBatchServiceMock{
		FuncCreate: func(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error) {
			return nil, nil
		},
		FuncGetReportById: func(ctx context.Context, id int) (*models.ReportProduct, error) { return nil, nil },
		FuncGetReport:     func(ctx context.Context) ([]models.ReportProduct, error) { return nil, nil },
		FuncSetLogger:     func(l logger.Logger) {},
	}

	m.CreateProductBatches(context.TODO(), models.ProductBatches{})
	m.GetReportProductById(context.TODO(), 0)
	m.GetReportProduct(context.TODO())
	m.SetLogger(nil)
}

func TestProductBatchRepositoryMock_DummyCoverage(t *testing.T) {
	m := &mocks.ProductBatchRepositoryMock{
		FuncCreate: func(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error) {
			return nil, nil
		},
		FuncGetReportById: func(ctx context.Context, id int) (*models.ReportProduct, error) { return nil, nil },
		FuncGetReport:     func(ctx context.Context) ([]models.ReportProduct, error) { return nil, nil },
		FuncSetLogger:     func(l logger.Logger) {},
	}

	m.CreateProductBatches(context.TODO(), models.ProductBatches{})
	m.GetReportProductById(context.TODO(), 0)
	m.GetReportProduct(context.TODO())
	m.SetLogger(nil)
}
