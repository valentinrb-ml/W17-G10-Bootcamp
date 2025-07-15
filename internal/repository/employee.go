package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

const (
	queryEmployeeCountByCardNumberID  = `SELECT COUNT(*) FROM employees WHERE id_card_number=?`
	queryEmployeeInsert               = `INSERT INTO employees (id_card_number, first_name, last_name, wareHouse_id) VALUES (?, ?, ?, ?)`
	queryEmployeeSelectByCardNumberID = `SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id_card_number=?`
	queryEmployeeSelectAll            = `SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees`
	queryEmployeeSelectByID           = `SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id=?`
	queryEmployeeUpdate               = `UPDATE employees SET id_card_number=?, first_name=?, last_name=?, wareHouse_id=? WHERE id=?`
	queryEmployeeDelete               = `DELETE FROM employees WHERE id=?`
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
		return false, apperrors.NewAppError(apperrors.CodeInternal, "database query failed")
	}
	return count > 0, nil
}

func (r *EmployeeMySQLRepository) Create(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	result, err := r.db.ExecContext(ctx, queryEmployeeInsert, e.CardNumberID, e.FirstName, e.LastName, e.WarehouseID)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "database insert failed")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "could not get inserted ID")
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
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "database scan failed")
	}
	return e, nil
}

func (r *EmployeeMySQLRepository) FindAll(ctx context.Context) ([]*models.Employee, error) {
	rows, err := r.db.QueryContext(ctx, queryEmployeeSelectAll)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "database query failed")
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		e := &models.Employee{}
		if err := rows.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID); err != nil {
			return nil, apperrors.NewAppError(apperrors.CodeInternal, "database scan failed")
		}
		employees = append(employees, e)
	}
	if err := rows.Err(); err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "row iteration error")
	}
	return employees, nil
}

func (r *EmployeeMySQLRepository) FindByID(ctx context.Context, id int) (*models.Employee, error) {
	row := r.db.QueryRowContext(ctx, queryEmployeeSelectByID, id)
	e := &models.Employee{}
	err := row.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // PARA QUE EL SERVICE PUEDA DEVOLVER 409 O 404
	}
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "database scan failed")
	}
	return e, nil
}

func (r *EmployeeMySQLRepository) Update(ctx context.Context, id int, e *models.Employee) error {
	_, err := r.db.ExecContext(ctx, queryEmployeeUpdate, e.CardNumberID, e.FirstName, e.LastName, e.WarehouseID, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternal, "database update failed")
	}
	return nil
}

func (r *EmployeeMySQLRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, queryEmployeeDelete, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternal, "database delete failed")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternal, "rows affected failed")
	}
	if rows == 0 {
		return apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
	}
	return nil
}
