package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"

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
	db     *sql.DB
	logger logger.Logger
}

// Inserta un inbound order, maneja errores de duplicidad y FK (1452)
func NewInboundOrderRepository(db *sql.DB) *InboundOrderMySQLRepository {
	return &InboundOrderMySQLRepository{db: db}
}
func (r *InboundOrderMySQLRepository) SetLogger(l logger.Logger) {
	r.logger = l
}
func (r *InboundOrderMySQLRepository) Create(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "inboundorder-repository", "Creating new inbound order", map[string]interface{}{
			"order_number": o.OrderNumber,
			"employee_id":  o.EmployeeID,
			"warehouse_id": o.WarehouseID,
		})
	}
	res, err := r.db.ExecContext(ctx, queryInboundOrderInsert, o.OrderDate, o.OrderNumber, o.EmployeeID, o.ProductBatchID, o.WarehouseID)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "inboundorder-repository", "Failed to insert inbound order", err, map[string]interface{}{
				"order_number": o.OrderNumber,
			})
		}
		// Manejo específico para errores de MySQL
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062: // Unique constraint violation
				return nil, apperrors.NewAppError(apperrors.CodeConflict, "order_number already exists")
			case 1452: // FK constraint violation
				return nil, apperrors.NewAppError(apperrors.CodeNotFound, "related resource not found (check employee_id, product_batch_id, warehouse_id)")
			}
		}
		return nil, apperrors.Wrap(err, "inbound order insert failed")
	}
	id, err := res.LastInsertId()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "inboundorder-repository", "Failed to get inserted ID", err, map[string]interface{}{
				"order_number": o.OrderNumber,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "could not get inserted ID")
	}
	o.ID = int(id)
	if r.logger != nil {
		r.logger.Info(ctx, "inboundorder-repository", "Inbound order created successfully", map[string]interface{}{
			"inbound_order_id": o.ID,
			"order_number":     o.OrderNumber,
		})
	}
	return o, nil
}

// Verifica si un order_number ya existe
func (r *InboundOrderMySQLRepository) ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "inboundorder-repository", "Checking if order_number exists", map[string]interface{}{
			"order_number": orderNumber,
		})
	}
	var count int
	err := r.db.QueryRowContext(ctx, queryInboundOrderExistsByOrderNumber, orderNumber).Scan(&count)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "inboundorder-repository", "DB error while checking order_number existence", err, map[string]interface{}{
				"order_number": orderNumber,
			})
		}
		return false, apperrors.NewAppError(apperrors.CodeInternal, "db error")
	}
	exists := count > 0
	if r.logger != nil {
		r.logger.Info(ctx, "inboundorder-repository", "order_number exists check completed", map[string]interface{}{
			"order_number": orderNumber,
			"exists":       exists,
		})
	}
	return exists, nil
}

// Genera el reporte de inbound orders para todos los empleados
func (r *InboundOrderMySQLRepository) ReportAll(ctx context.Context) ([]models.InboundOrderReport, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "inboundorder-repository", "Generating report for all employees")
	}
	rows, err := r.db.QueryContext(ctx, queryInboundOrdersReportAll)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "inboundorder-repository", "DB query error in ReportAll", err)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "db query error")
	}
	defer rows.Close()

	var res []models.InboundOrderReport
	for rows.Next() {
		var rep models.InboundOrderReport
		err := rows.Scan(&rep.ID, &rep.CardNumberID, &rep.FirstName, &rep.LastName, &rep.WarehouseID, &rep.InboundOrdersCount)
		if err != nil {
			if r.logger != nil {
				r.logger.Error(ctx, "inboundorder-repository", "Error scanning row in ReportAll", err)
			}
			continue
		}
		res = append(res, rep)
	}
	if r.logger != nil {
		r.logger.Info(ctx, "inboundorder-repository", "Report for all employees generated successfully", map[string]interface{}{
			"employees_count": len(res),
		})
	}
	return res, nil
}

// Genera el reporte de inbound orders para un empleado por id
func (r *InboundOrderMySQLRepository) ReportByID(ctx context.Context, employeeID int) (*models.InboundOrderReport, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "inboundorder-repository", "Generating report for employee", map[string]interface{}{
			"employee_id": employeeID,
		})
	}
	row := r.db.QueryRowContext(ctx, queryInboundOrdersReportByEmployee, employeeID)
	rep := &models.InboundOrderReport{}
	err := row.Scan(&rep.ID, &rep.CardNumberID, &rep.FirstName, &rep.LastName, &rep.WarehouseID, &rep.InboundOrdersCount)
	if errors.Is(err, sql.ErrNoRows) {
		if r.logger != nil {
			r.logger.Warning(ctx, "inboundorder-repository", "Employee not found in report", map[string]interface{}{
				"employee_id": employeeID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
	}
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "inboundorder-repository", "DB scan error in ReportByID", err, map[string]interface{}{
				"employee_id": employeeID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "db scan error")
	}
	if r.logger != nil {
		r.logger.Info(ctx, "inboundorder-repository", "Report for employee generated successfully", map[string]interface{}{
			"employee_id":          rep.ID,
			"inbound_orders_count": rep.InboundOrdersCount,
		})
	}
	return rep, nil
}
