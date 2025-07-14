package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

const (
	queryEmployeeCountByCardNumberID = `
		SELECT COUNT(*) FROM employees WHERE id_card_number=?`
	queryEmployeeInsert = `
		INSERT INTO employees (id_card_number, first_name, last_name, warehouse_id)
		VALUES (?, ?, ?, ?)`
	queryEmployeeSelectByCardNumberID = `
		SELECT id, id_card_number, first_name, last_name, warehouse_id
		FROM employees WHERE id_card_number=?`
	queryEmployeeSelectAll = `
		SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees`
	queryEmployeeSelectByID = `
		SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id=?`
	queryEmployeeUpdate = `
		UPDATE employees SET id_card_number=?, first_name=?, last_name=?, warehouse_id=? WHERE id=?`
	queryEmployeeDelete = `
		DELETE FROM employees WHERE id=?`
)

type EmployeeMySQLRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeMySQLRepository {
	return &EmployeeMySQLRepository{db: db}
}

func (r *EmployeeMySQLRepository) ExistsByCardNumberID(ctx context.Context, cardNumberID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, queryEmployeeCountByCardNumberID, cardNumberID).Scan(&count)
	if err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "database query failed"
		se.InternalError = err
		return false, &se
	}
	return count > 0, nil
}

func (r *EmployeeMySQLRepository) Create(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	result, err := r.db.ExecContext(ctx, queryEmployeeInsert, e.CardNumberID, e.FirstName, e.LastName, e.WarehouseID)
	if err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "database insert failed"
		se.InternalError = err
		return nil, &se
	}
	id, err := result.LastInsertId()
	if err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "could not get inserted ID"
		se.InternalError = err
		return nil, &se
	}
	e.ID = int(id)
	return e, nil
}

func (r *EmployeeMySQLRepository) FindByCardNumberID(ctx context.Context, cardNumberID string) (*models.Employee, error) {
	row := r.db.QueryRowContext(ctx, queryEmployeeSelectByCardNumberID, cardNumberID)
	e := &models.Employee{}
	err := row.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "database scan failed"
		se.InternalError = err
		return nil, &se
	}
	return e, nil
}

func (r *EmployeeMySQLRepository) FindAll(ctx context.Context) ([]*models.Employee, error) {
	rows, err := r.db.QueryContext(ctx, queryEmployeeSelectAll)
	if err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "database query failed"
		se.InternalError = err
		return nil, &se
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		e := &models.Employee{}
		if err := rows.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID); err != nil {
			se := api.ServiceErrors[api.ErrInternalServer]
			se.Message = "database scan failed"
			se.InternalError = err
			return nil, &se
		}
		employees = append(employees, e)
	}
	if err := rows.Err(); err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "row iteration error"
		se.InternalError = err
		return nil, &se
	}
	return employees, nil
}

func (r *EmployeeMySQLRepository) FindByID(ctx context.Context, id int) (*models.Employee, error) {
	row := r.db.QueryRowContext(ctx, queryEmployeeSelectByID, id)
	e := &models.Employee{}
	err := row.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID)
	if errors.Is(err, sql.ErrNoRows) {
		se := api.ServiceErrors[api.ErrNotFound]
		se.Message = "employee not found"
		return nil, &se
	}
	if err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "database scan failed"
		se.InternalError = err
		return nil, &se
	}
	return e, nil
}

func (r *EmployeeMySQLRepository) Update(ctx context.Context, id int, e *models.Employee) (*models.Employee, error) {
	_, err := r.db.ExecContext(
		ctx,
		queryEmployeeUpdate,
		e.CardNumberID, e.FirstName, e.LastName, e.WarehouseID, id,
	)
	if err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "database update failed"
		se.InternalError = err
		return nil, &se
	}
	return r.FindByID(ctx, id)
}

func (r *EmployeeMySQLRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, queryEmployeeDelete, id)
	if err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "database delete failed"
		se.InternalError = err
		return &se
	}
	rows, err := result.RowsAffected()
	if err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "rows affected failed"
		se.InternalError = err
		return &se
	}
	if rows == 0 {
		se := api.ServiceErrors[api.ErrNotFound]
		se.Message = "employee not found"
		return &se
	}
	return nil
}
