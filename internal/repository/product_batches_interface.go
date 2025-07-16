package repository

import (
	"context"
	"database/sql"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

type ProductBatchesRepository interface {
	CreateProductBatches(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error)
	GetReportProductById(ctx context.Context, id int) (*models.ReportProduct, error)
	GetReportProduct(ctx context.Context) ([]models.ReportProduct, error)
}

type productBatchesRepository struct {
	mysql *sql.DB
}

func NewProductBatchesRepository(mysql *sql.DB) ProductBatchesRepository {
	return &productBatchesRepository{
		mysql,
	}
}
