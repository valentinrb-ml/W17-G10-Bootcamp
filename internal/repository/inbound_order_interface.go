package repository

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

// declaración de interface del repositorio de inbound orders
type InboundOrderRepository interface {
	Create(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error)
	ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error)
	ReportAll(ctx context.Context) ([]models.InboundOrderReport, error)
	ReportByID(ctx context.Context, employeeID int) (*models.InboundOrderReport, error)
}
