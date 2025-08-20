package product_record

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"

	"github.com/stretchr/testify/mock"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

type MockProductRecordRepository struct {
	mock.Mock
	log logger.Logger
}

func (m *MockProductRecordRepository) Create(ctx context.Context, r models.ProductRecord) (models.ProductRecord, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(models.ProductRecord), args.Error(1)
}

func (m *MockProductRecordRepository) GetRecordsReport(ctx context.Context, id int) ([]models.ProductRecordReport, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]models.ProductRecordReport), args.Error(1)
}

func (m *MockProductRecordRepository) SetLogger(l logger.Logger) {
	m.log = l
}
