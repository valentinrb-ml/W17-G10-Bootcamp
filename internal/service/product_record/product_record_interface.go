package product_record

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

type ProductRecordService interface {
	Create(ctx context.Context, record models.ProductRecord) (models.ProductRecordResponse, error)
	GetRecordsReport(ctx context.Context, productID int) ([]models.ProductRecordReport, error)
}
