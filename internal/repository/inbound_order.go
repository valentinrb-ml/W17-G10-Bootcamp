package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

const (
	// Inserta un nuevo inbound order, validando FKs con la base
	queryInboundOrderInsert = `
    	INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id)
    		VALUES (?, ?, ?, ?, ?)`
	// Verifica si un order_number ya existe (debe ser único)
	queryInboundOrderExistsByOrderNumber = `
		SELECT COUNT(1) FROM inbound_orders WHERE order_number = ?`
	// Trae el reporte de inbound orders agrupado por employee
	queryInboundOrdersReportAll = `
		SELECT e.id, e.id_card_number, e.first_name, e.last_name, e.warehouse_id, 
		  COUNT(io.id) as inbound_orders_count
		FROM employees e
		LEFT JOIN inbound_orders io ON e.id = io.employee_id
		GROUP BY e.id`
	// Trae el reporte para un sólo employee especificado por id
	queryInboundOrdersReportByEmployee = `
		SELECT e.id, e.id_card_number, e.first_name, e.last_name, e.warehouse_id, 
		  COUNT(io.id) as inbound_orders_count
		FROM employees e
		LEFT JOIN inbound_orders io ON e.id = io.employee_id
		WHERE e.id = ?
		GROUP BY e.id`
)

// Repositorio MySQL para inbound orders
type InboundOrderMySQLRepository struct {
	db *sql.DB
}

// Inserta un inbound order, maneja errores de duplicidad y FK (1452)
func NewInboundOrderRepository(db *sql.DB) *InboundOrderMySQLRepository {
	return &InboundOrderMySQLRepository{db: db}
}
func (r *InboundOrderMySQLRepository) Create(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error) {
	res, err := r.db.ExecContext(ctx, queryInboundOrderInsert, o.OrderDate, o.OrderNumber, o.EmployeeID, o.ProductBatchID, o.WarehouseID)
	if err != nil {
		// Manejo específico para errores de MySQL
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062: // Unique constraint violation
				return nil, apperrors.NewAppError(apperrors.CodeConflict, "order_number already exists")
			case 1452: // FK constraint violation
				return nil, apperrors.NewAppError(apperrors.CodeUnprocessableEntity, "invalid foreign key: check employee_id, product_batch_id, warehouse_id")
			}
		}
		return nil, apperrors.Wrap(err, "inbound order insert failed")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "could not get inserted ID")
	}
	o.ID = int(id)
	return o, nil
}

// Verifica si un order_number ya existe
func (r *InboundOrderMySQLRepository) ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, queryInboundOrderExistsByOrderNumber, orderNumber).Scan(&count)
	if err != nil {
		return false, apperrors.NewAppError(apperrors.CodeInternal, "db error")
	}
	return count > 0, nil
}

// Genera el reporte de inbound orders para todos los empleados
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

// Genera el reporte de inbound orders para un empleado por id
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
