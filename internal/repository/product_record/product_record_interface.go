package repository

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

type ProductRecordRepository interface {
	Create(ctx context.Context, record models.ProductRecord) (models.ProductRecord, error)
	GetRecordsReport(ctx context.Context, productID int) ([]models.ProductRecordReport, error)
}
