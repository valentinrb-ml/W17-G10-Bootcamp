package service

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/productrecord"
)

type ProductRecordService interface {
	Create(ctx context.Context, record productrecord.ProductRecord) (productrecord.ProductRecordResponse, error)
	GetRecordsReport(ctx context.Context, productID int) (productrecord.ProductRecordsReportResponse, error)
}
