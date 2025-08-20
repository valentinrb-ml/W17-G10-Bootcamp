package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

const (
	// Query para insertar un empleado
	queryEmployeeInsert = `INSERT INTO employees (id_card_number, first_name, last_name, wareHouse_id) VALUES (?, ?, ?, ?)`
	// Query para traer empleado por card_number_id (para comprobar duplicados)
	queryEmployeeSelectByCardNumberID = `SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id_card_number=?`
	// Query para traer todos los empleados
	queryEmployeeSelectAll = `SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees`
	// Query para traer empleado por id
	queryEmployeeSelectByID = `SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id=?`
	// Query para actualizar empleado
	queryEmployeeUpdate = `UPDATE employees SET id_card_number=?, first_name=?, last_name=?, wareHouse_id=? WHERE id=?`
	// Query para borrar empleado
	queryEmployeeDelete = `DELETE FROM employees WHERE id=?`
)

// Implementación MySQL del repositorio de empleados
type EmployeeMySQLRepository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewEmployeeRepository(db *sql.DB) *EmployeeMySQLRepository {
	return &EmployeeMySQLRepository{db: db}
}
func (r *EmployeeMySQLRepository) SetLogger(l logger.Logger) {
	r.logger = l
}

// Crea un nuevo empleado
func (r *EmployeeMySQLRepository) Create(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Creating new employee", map[string]interface{}{
			"card_number_id": e.CardNumberID,
			"first_name":     e.FirstName,
			"last_name":      e.LastName,
			"warehouse_id":   e.WarehouseID,
		})
	}
	result, err := r.db.ExecContext(ctx, queryEmployeeInsert, e.CardNumberID, e.FirstName, e.LastName, e.WarehouseID)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "employee-repository", "Failed to create employee", err, map[string]interface{}{
				"card_number_id": e.CardNumberID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "database insert failed")
	}
	id, err := result.LastInsertId()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "employee-repository", "Failed to get last insert ID", err, map[string]interface{}{
				"card_number_id": e.CardNumberID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "could not get inserted ID")
	}
	e.ID = int(id)
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Employee created successfully", map[string]interface{}{
			"employee_id":    e.ID,
			"card_number_id": e.CardNumberID,
		})
	}
	return e, nil
}

// Busca un empleado por card_number_id (para unicidad)
func (r *EmployeeMySQLRepository) FindByCardNumberID(ctx context.Context, cardNumberID string) (*models.Employee, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Searching employee by card_number_id", map[string]interface{}{
			"card_number_id": cardNumberID,
		})
	}
	row := r.db.QueryRowContext(ctx, queryEmployeeSelectByCardNumberID, cardNumberID)
	e := &models.Employee{}
	err := row.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID)
	if errors.Is(err, sql.ErrNoRows) {
		if r.logger != nil {
			r.logger.Warning(ctx, "employee-repository", "Employee not found with card_number_id", map[string]interface{}{
				"card_number_id": cardNumberID,
			})
		}
		return nil, nil
	}
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "employee-repository", "Failed to scan employee by card_number_id", err, map[string]interface{}{
				"card_number_id": cardNumberID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "database scan failed")
	}
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Employee found by card_number_id", map[string]interface{}{
			"employee_id":    e.ID,
			"card_number_id": cardNumberID,
		})
	}
	return e, nil
}

// Devuelve todos los empleados
func (r *EmployeeMySQLRepository) FindAll(ctx context.Context) ([]*models.Employee, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Fetching all employees")
	}
	rows, err := r.db.QueryContext(ctx, queryEmployeeSelectAll)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "employee-repository", "Failed to query all employees", err)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "database query failed")
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		e := &models.Employee{}
		if err := rows.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID); err != nil {
			if r.logger != nil {
				r.logger.Error(ctx, "employee-repository", "Failed to scan employee row", err)
			}
			return nil, apperrors.NewAppError(apperrors.CodeInternal, "database scan failed")
		}
		employees = append(employees, e)
	}
	if err := rows.Err(); err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "employee-repository", "Error during row iteration", err)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "row iteration error")
	}
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Fetched all employees", map[string]interface{}{
			"count": len(employees),
		})
	}
	return employees, nil
}

// Busca un empleado por id
func (r *EmployeeMySQLRepository) FindByID(ctx context.Context, id int) (*models.Employee, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Searching employee by id", map[string]interface{}{
			"employee_id": id,
		})
	}
	row := r.db.QueryRowContext(ctx, queryEmployeeSelectByID, id)
	e := &models.Employee{}
	err := row.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID)
	if errors.Is(err, sql.ErrNoRows) {
		if r.logger != nil {
			r.logger.Warning(ctx, "employee-repository", "Employee not found by id", map[string]interface{}{
				"employee_id": id,
			})
		}
		return nil, nil // PARA QUE EL SERVICE PUEDA DEVOLVER 409 O 404
	}
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "employee-repository", "Failed to scan employee by id", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "database scan failed")
	}
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Employee found by id", map[string]interface{}{
			"employee_id": e.ID,
		})
	}
	return e, nil
}

// Actualiza un empleado existente
func (r *EmployeeMySQLRepository) Update(ctx context.Context, id int, e *models.Employee) error {
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Updating employee", map[string]interface{}{
			"employee_id":    id,
			"card_number_id": e.CardNumberID,
		})
	}
	_, err := r.db.ExecContext(ctx, queryEmployeeUpdate, e.CardNumberID, e.FirstName, e.LastName, e.WarehouseID, id)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "employee-repository", "Failed to update employee", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeInternal, "database update failed")
	}
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Employee updated successfully", map[string]interface{}{
			"employee_id": id,
		})
	}
	return nil
}

// Borra un empleado por id
func (r *EmployeeMySQLRepository) Delete(ctx context.Context, id int) error {
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Deleting employee", map[string]interface{}{
			"employee_id": id,
		})
	}
	result, err := r.db.ExecContext(ctx, queryEmployeeDelete, id)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "employee-repository", "Failed to delete employee", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeInternal, "database delete failed")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "employee-repository", "Failed to get rows affected", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeInternal, "rows affected failed")
	}
	if rows == 0 {
		if r.logger != nil {
			r.logger.Warning(ctx, "employee-repository", "Employee not found for deletion", map[string]interface{}{
				"employee_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
	}
	if r.logger != nil {
		r.logger.Info(ctx, "employee-repository", "Employee deleted successfully", map[string]interface{}{
			"employee_id": id,
		})
	}
	return nil
}
