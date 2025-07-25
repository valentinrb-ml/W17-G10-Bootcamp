package mocks

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

type InboundOrderRepositoryMock struct {
	MockCreate              func(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error)
	MockExistsByOrderNumber func(ctx context.Context, orderNumber string) (bool, error)
	MockReportAll           func(ctx context.Context) ([]models.InboundOrderReport, error)
	MockReportByID          func(ctx context.Context, employeeID int) (*models.InboundOrderReport, error)
}

func (m *InboundOrderRepositoryMock) Create(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error) {
	return m.MockCreate(ctx, o)
}
func (m *InboundOrderRepositoryMock) ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	return m.MockExistsByOrderNumber(ctx, orderNumber)
}
func (m *InboundOrderRepositoryMock) ReportAll(ctx context.Context) ([]models.InboundOrderReport, error) {
	return m.MockReportAll(ctx)
}
func (m *InboundOrderRepositoryMock) ReportByID(ctx context.Context, employeeID int) (*models.InboundOrderReport, error) {
	return m.MockReportByID(ctx, employeeID)
}
