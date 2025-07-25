package mocks

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

// Interfaz de service
type InboundOrderServiceMock struct {
	MockCreate func(ctx context.Context, in *models.InboundOrder) (*models.InboundOrder, error)
	MockReport func(ctx context.Context, id *int) (interface{}, error)
}

func (m *InboundOrderServiceMock) Create(ctx context.Context, in *models.InboundOrder) (*models.InboundOrder, error) {
	return m.MockCreate(ctx, in)
}
func (m *InboundOrderServiceMock) Report(ctx context.Context, id *int) (interface{}, error) {
	return m.MockReport(ctx, id)
}
