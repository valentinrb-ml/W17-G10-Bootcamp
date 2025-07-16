package repository

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/productrecord"
)

type ProductRecordRepository interface {
	Create(ctx context.Context, record productrecord.ProductRecord) (productrecord.ProductRecord, error)
	GetRecordsReport(ctx context.Context, productID int) ([]productrecord.ProductRecordReport, error)
}
