package repository

import (
	"context"
	"database/sql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

// ProductBatchesRepository defines the interface for product batches data operations.
type ProductBatchesRepository interface {
	CreateProductBatches(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error)
	GetReportProductById(ctx context.Context, id int) (*models.ReportProduct, error)
	GetReportProduct(ctx context.Context) ([]models.ReportProduct, error)
}

// productBatchesRepository is the implementation of ProductBatchesRepository using MySQL.
// This struct holds the DB connection.
type productBatchesRepository struct {
	mysql  *sql.DB
	logger logger.Logger
}

// SetLogger allows you to inject the logger after creation
func (r *productBatchesRepository) SetLogger(l logger.Logger) {
	r.logger = l
}

// NewProductBatchesRepository returns a new ProductBatchesRepository using the given MySQL connection.
func NewProductBatchesRepository(mysql *sql.DB) ProductBatchesRepository {
	return &productBatchesRepository{
		mysql: mysql,
	}
}
