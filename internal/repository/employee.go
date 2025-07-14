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

type EmployeeRepository interface {
	Create(ctx context.Context, e *models.Employee) (*models.Employee, error)
	FindByCardNumberID(ctx context.Context, cardNumberID string) (*models.Employee, error)
	FindAll(ctx context.Context) ([]*models.Employee, error)
	FindByID(ctx context.Context, id int) (*models.Employee, error)
	Update(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error)
	Delete(ctx context.Context, id int) error
}

type EmployeeMySQLRepository struct {
	db *sql.DB
}

func NewEmployeeMySQLRepository(db *sql.DB) *EmployeeMySQLRepository {
	return &EmployeeMySQLRepository{db: db}
}

func (r *EmployeeMySQLRepository) Create(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	var count int
	err := r.db.QueryRowContext(ctx, queryEmployeeCountByCardNumberID, e.CardNumberID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		se := api.ServiceErrors[api.ErrConflict]
		se.Message = "card_number_id already exists"
		return nil, &se
	}

	result, err := r.db.ExecContext(ctx, queryEmployeeInsert, e.CardNumberID, e.FirstName, e.LastName, e.WarehouseID)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return e, nil
}

func (r *EmployeeMySQLRepository) FindAll(ctx context.Context) ([]*models.Employee, error) {
	rows, err := r.db.QueryContext(ctx, queryEmployeeSelectAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		e := &models.Employee{}
		if err := rows.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID); err != nil {
			return nil, err
		}
		employees = append(employees, e)
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
		return nil, err
	}
	return e, nil
}

func (r *EmployeeMySQLRepository) Update(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error) {
	emp, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		se := api.ServiceErrors[api.ErrNotFound]
		se.Message = "employee not found"
		return nil, &se
	}

	if patch.CardNumberID != nil {
		exist, err := r.FindByCardNumberID(ctx, *patch.CardNumberID)
		if err != nil {
			return nil, err
		}
		if exist != nil && exist.ID != id {
			se := api.ServiceErrors[api.ErrConflict]
			se.Message = "card_number_id already exists"
			return nil, &se
		}
		emp.CardNumberID = *patch.CardNumberID
	}
	if patch.FirstName != nil {
		emp.FirstName = *patch.FirstName
	}
	if patch.LastName != nil {
		emp.LastName = *patch.LastName
	}
	if patch.WarehouseID != nil && *patch.WarehouseID != 0 {
		emp.WarehouseID = *patch.WarehouseID
	}
	_, err = r.db.ExecContext(ctx, queryEmployeeUpdate,
		emp.CardNumberID, emp.FirstName, emp.LastName, emp.WarehouseID, emp.ID,
	)
	if err != nil {
		return nil, err
	}
	return emp, nil
}

func (r *EmployeeMySQLRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, queryEmployeeDelete, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		se := api.ServiceErrors[api.ErrNotFound]
		se.Message = "employee not found"
		return &se
	}
	return nil
}
