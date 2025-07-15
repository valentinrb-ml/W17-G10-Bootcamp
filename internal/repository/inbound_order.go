package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

const (
	queryInboundOrderInsert = `
    	INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id)
    		VALUES (?, ?, ?, ?, ?)`
	queryInboundOrderExistsByOrderNumber = `
		SELECT COUNT(1) FROM inbound_orders WHERE order_number = ?`
	queryInboundOrdersReportAll = `
		SELECT e.id, e.id_card_number, e.first_name, e.last_name, e.warehouse_id, 
		  COUNT(io.id) as inbound_orders_count
		FROM employees e
		LEFT JOIN inbound_orders io ON e.id = io.employee_id
		GROUP BY e.id`
	queryInboundOrdersReportByEmployee = `
		SELECT e.id, e.id_card_number, e.first_name, e.last_name, e.warehouse_id, 
		  COUNT(io.id) as inbound_orders_count
		FROM employees e
		LEFT JOIN inbound_orders io ON e.id = io.employee_id
		WHERE e.id = ?
		GROUP BY e.id`
)

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

func (r *InboundOrderMySQLRepository) ReportAll(ctx context.Context) ([]models.InboundOrderReport, error) {
	rows, err := r.db.QueryContext(ctx, queryInboundOrdersReportAll)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "db query error")
	}
	defer rows.Close()

	var res []models.InboundOrderReport
	for rows.Next() {
		var rep models.InboundOrderReport
		err := rows.Scan(&rep.ID, &rep.CardNumberID, &rep.FirstName, &rep.LastName, &rep.WarehouseID, &rep.InboundOrdersCount)
		if err != nil {
			continue
		}
		res = append(res, rep)
	}
	return res, nil
}

func (r *InboundOrderMySQLRepository) ReportByID(ctx context.Context, employeeID int) (*models.InboundOrderReport, error) {
	row := r.db.QueryRowContext(ctx, queryInboundOrdersReportByEmployee, employeeID)
	rep := &models.InboundOrderReport{}
	err := row.Scan(&rep.ID, &rep.CardNumberID, &rep.FirstName, &rep.LastName, &rep.WarehouseID, &rep.InboundOrdersCount)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
	}
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "db scan error")
	}
	return rep, nil
}
