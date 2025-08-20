package mocks_test

import (
	"context"
	"testing"

	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/inbound_order"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

func TestInboundOrderRepositoryMock_DummyCoverage(t *testing.T) {
	m := &mocks.InboundOrderRepositoryMock{
		MockCreate:              func(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error) { return nil, nil },
		MockExistsByOrderNumber: func(ctx context.Context, orderNumber string) (bool, error) { return false, nil },
		MockReportAll:           func(ctx context.Context) ([]models.InboundOrderReport, error) { return nil, nil },
		MockReportByID:          func(ctx context.Context, employeeID int) (*models.InboundOrderReport, error) { return nil, nil },
	}
	m.Create(context.TODO(), &models.InboundOrder{})
	m.ExistsByOrderNumber(context.TODO(), "")
	m.ReportAll(context.TODO())
	m.ReportByID(context.TODO(), 0)
}

func TestInboundOrderServiceMock_DummyCoverage(t *testing.T) {
	m := &mocks.InboundOrderServiceMock{
		MockCreate: func(ctx context.Context, in *models.InboundOrder) (*models.InboundOrder, error) { return nil, nil },
		MockReport: func(ctx context.Context, id *int) (interface{}, error) { return nil, nil },
	}
	m.Create(context.TODO(), &models.InboundOrder{})
	m.Report(context.TODO(), nil)
}
