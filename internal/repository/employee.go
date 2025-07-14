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
	Update(ctx context.Context, id int, e *models.Employee) (*models.Employee, error)
	Delete(ctx context.Context, id int) error
	ExistsByCardNumberID(ctx context.Context, cardNumberID string) (bool, error)
}

type EmployeeMySQLRepository struct {
	db *sql.DB
}

func NewEmployeeMySQLRepository(db *sql.DB) *EmployeeMySQLRepository {
	return &EmployeeMySQLRepository{db: db}
}

func (r *EmployeeMySQLRepository) ExistsByCardNumberID(ctx context.Context, cardNumberID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, queryEmployeeCountByCardNumberID, cardNumberID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *EmployeeMySQLRepository) Create(ctx context.Context, e *models.Employee) (*models.Employee, error) {
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

func (r *EmployeeMySQLRepository) Update(ctx context.Context, id int, e *models.Employee) (*models.Employee, error) {
	_, err := r.db.ExecContext(
		ctx,
		queryEmployeeUpdate,
		e.CardNumberID, e.FirstName, e.LastName, e.WarehouseID, id,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
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
