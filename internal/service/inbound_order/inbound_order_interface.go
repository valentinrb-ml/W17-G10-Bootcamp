package service

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

type InboundOrderService interface {
	Create(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error)
	Report(ctx context.Context, employeeID *int) (interface{}, error)
}
