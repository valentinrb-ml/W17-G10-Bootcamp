package repository

import (
	"context"
	"database/sql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

const (
	queryInboundOrderInsert = `
    	INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id)
    		VALUES (?, ?, ?, ?, ?)`
	queryInboundOrderExistsByOrderNumber = `
		SELECT COUNT(1) FROM inbound_orders WHERE order_number = ?`
)

type InboundOrderRepository interface {
	Create(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error)
	ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error)
}

type InboundOrderMySQLRepository struct {
	db *sql.DB
}

func NewInboundOrderRepository(db *sql.DB) *InboundOrderMySQLRepository {
	return &InboundOrderMySQLRepository{db: db}
}

func (r *InboundOrderMySQLRepository) Create(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error) {
	res, err := r.db.ExecContext(ctx, queryInboundOrderInsert, o.OrderDate, o.OrderNumber, o.EmployeeID, o.ProductBatchID, o.WarehouseID)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "inbound order insert failed")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "could not get inserted ID")
	}
	o.ID = int(id)
	return o, nil
}
func (r *InboundOrderMySQLRepository) ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, queryInboundOrderExistsByOrderNumber, orderNumber).Scan(&count)
	if err != nil {
		return false, apperrors.NewAppError(apperrors.CodeInternal, "db error")
	}
	return count > 0, nil
}
