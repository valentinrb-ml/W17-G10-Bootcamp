// internal/repository/purchase_order_interface.go
package repository

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type PurchaseOrderRepository interface {
	Create(ctx context.Context, po models.PurchaseOrder) (*models.PurchaseOrder, error)
	GetAll(ctx context.Context) ([]models.PurchaseOrder, error)
	GetByID(ctx context.Context, id int) (*models.PurchaseOrder, error)
	ExistsOrderNumber(ctx context.Context, orderNumber string) bool
	GetCountByBuyer(ctx context.Context, buyerID int) ([]models.BuyerWithPurchaseCount, error)
	GetAllWithPurchaseCount(ctx context.Context) ([]models.BuyerWithPurchaseCount, error)
}
