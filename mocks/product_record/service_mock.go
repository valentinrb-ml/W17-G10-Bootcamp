package product_record

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"

	"github.com/stretchr/testify/mock"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

type MockProductRecordService struct {
	mock.Mock
	log logger.Logger
}

func (m *MockProductRecordService) Create(ctx context.Context, r models.ProductRecord) (models.ProductRecord, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(models.ProductRecord), args.Error(1)
}

/*
GetRecordsReport mimics the SELECT-report query.

If the test returns nil instead of an empty slice we guard against the panic that would occur when doing nil.(slice).
*/
func (m *MockProductRecordService) GetRecordsReport(ctx context.Context, id int) ([]models.ProductRecordReport, error) {
	args := m.Called(ctx, id)

	// nil-safe conversion of the first return value
	var rep []models.ProductRecordReport
	if arg0 := args.Get(0); arg0 != nil {
		rep = arg0.([]models.ProductRecordReport)
	}
	return rep, args.Error(1)
}

func (m *MockProductRecordService) SetLogger(l logger.Logger) {
	m.log = l
}
