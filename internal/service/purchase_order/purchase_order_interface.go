package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

// PurchaseOrderService define los métodos disponibles para interactuar con Purchase Orders
type PurchaseOrderService interface {
	// Create registra una nueva Purchase Order
	Create(ctx context.Context, req models.RequestPurchaseOrder) (*models.ResponsePurchaseOrder, error)

	// GetAll obtiene todas las Purchase Orders
	GetAll(ctx context.Context) ([]models.ResponsePurchaseOrder, error)

	// GetByID obtiene una Purchase Order por su ID
	GetByID(ctx context.Context, id int) (*models.ResponsePurchaseOrder, error)

	// GetReportByBuyer genera el reporte de Purchase Orders por Buyer
	// Si buyerID es nil, devuelve el reporte para todos los buyers
	GetReportByBuyer(ctx context.Context, buyerID *int) ([]models.BuyerWithPurchaseCount, error)

	// SetLogger allows injecting the logger after creation
	SetLogger(l logger.Logger)
}
