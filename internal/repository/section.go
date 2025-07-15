package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

type SectionMap struct {
	db *sql.DB
}

// NewSectionMap is a function that returns a new instance of SectionMap
func NewSectionMap(db *sql.DB) *SectionMap {
	return &SectionMap{db}
}

const (
	querySectionGetAll = `SELECT id, section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id FROM sections `
	querySectionGetOne = `SELECT id, section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id FROM sections WHERE id = ?`
	querySectionDelete = `DELETE FROM sections WHERE id =?`
	querySectionUpdate = `UPDATE sections SET section_number = ?, current_capacity = ?, current_temperature = ? , maximum_capacity = ?, minimum_capacity = ?, minimum_temperature = ?, product_type_id = ?, warehouse_id = ? WHERE id = ?`
	querySectionCreate = `INSERT INTO sections (section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id) VALUES (?,?,?,?,?,?,?,?)`
)

// TODO Mejorar errores
func (r *SectionMap) FindAllSections(ctx context.Context) ([]section.Section, error) {
	rows, err := r.db.QueryContext(ctx, querySectionGetAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []section.Section

	for rows.Next() {
		var s section.Section
		err := rows.Scan(&s.Id, &s.SectionNumber, &s.CurrentCapacity, &s.CurrentTemperature, &s.MaximumCapacity, &s.MinimumCapacity, &s.MinimumTemperature, &s.ProductTypeId, &s.WarehouseId)
		if err != nil {
			return nil, err
		}
		sections = append(sections, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}

func (r *SectionMap) FindById(ctx context.Context, id int) (*section.Section, error) {
	var s section.Section
	err := r.db.QueryRowContext(ctx, querySectionGetOne, id).Scan(&s.Id, &s.SectionNumber, &s.CurrentCapacity, &s.CurrentTemperature, &s.MaximumCapacity, &s.MinimumCapacity, &s.MinimumTemperature, &s.ProductTypeId, &s.WarehouseId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := api.ServiceErrors[api.ErrNotFound]
			return nil, &api.ServiceError{
				Code:         err.Code,
				ResponseCode: err.ResponseCode,
				Message:      "The section you are looking for does not exist.",
			}
		}
		err := api.ServiceErrors[api.ErrInternalServer]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "An internal server error occurred while retrieving the section.",
		}
	}

	return &s, nil
}

func (r *SectionMap) DeleteSection(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, querySectionDelete, id)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1451 {
			err := api.ServiceErrors[api.ErrConflict]
			return &api.ServiceError{
				Code:         err.Code,
				ResponseCode: err.ResponseCode,
				Message:      "Cannot delete section: there are products batches associated with this section.",
			}
		}

		errVal := api.ServiceErrors[api.ErrInternalServer]
		return &api.ServiceError{
			Code:         errVal.Code,
			ResponseCode: errVal.ResponseCode,
			Message:      err.Error(),
		}
	}
	rows, err := result.RowsAffected()
	if err != nil {
		errVal := api.ServiceErrors[api.ErrInternalServer]
		return &api.ServiceError{
			Code:         errVal.Code,
			ResponseCode: errVal.ResponseCode,
			Message:      "An internal server error occurred while deleting the section.",
		}
	}
	if rows == 0 {
		err := api.ServiceErrors[api.ErrNotFound]
		return &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The section you are trying to delete does not exist.",
		}
	}

	return nil
}

func (r *SectionMap) CreateSection(ctx context.Context, sec section.Section) (*section.Section, error) {
	result, err := r.db.ExecContext(ctx, querySectionCreate,
		sec.SectionNumber, sec.CurrentCapacity, sec.CurrentTemperature, sec.MaximumCapacity, sec.MinimumCapacity, sec.MinimumTemperature, sec.ProductTypeId, sec.WarehouseId)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			err := api.ServiceErrors[api.ErrConflict]
			return nil, &api.ServiceError{
				Code:         err.Code,
				ResponseCode: err.ResponseCode,
				Message:      "Section number already exists.",
			}
		}
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1452 {
			err := api.ServiceErrors[api.ErrBadRequest]
			return nil, &api.ServiceError{
				Code:         err.Code,
				ResponseCode: err.ResponseCode,
				Message:      "Warehouse id or product type id does not exist.",
			}
		}
		err := api.ServiceErrors[api.ErrInternalServer]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "An internal error occurred while creating the section.",
		}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	sec.Id = int(id)

	return &sec, nil
}

func (r *SectionMap) UpdateSection(ctx context.Context, id int, sec *section.Section) (*section.Section, error) {
	result, err := r.db.ExecContext(ctx, querySectionUpdate,
		sec.SectionNumber, sec.CurrentCapacity, sec.CurrentTemperature,
		sec.MaximumCapacity, sec.MinimumCapacity, sec.MinimumTemperature,
		sec.ProductTypeId, sec.WarehouseId, id)
	if err != nil {
		//constraint UNIQUE
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			// Duplicate entry (section_number)
			err := api.ServiceErrors[api.ErrConflict]
			return nil, &api.ServiceError{
				Code:         err.Code,
				ResponseCode: err.ResponseCode,
				Message:      "Section number already exists.",
			}
		}
		// constraint FOREIGN KEY
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1452 {
			err := api.ServiceErrors[api.ErrBadRequest]
			return nil, &api.ServiceError{
				Code:         err.Code,
				ResponseCode: err.ResponseCode,
				Message:      "Warehouse id or product type id does not exist.",
			}
		}
		err := api.ServiceErrors[api.ErrInternalServer]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "An internal server error occurred while updating the section.",
		}

	}

	rows, err := result.RowsAffected()
	if err != nil {
		err := api.ServiceErrors[api.ErrInternalServer]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "An internal server error occurred while updating the section.",
		}
	}
	if rows == 0 {
		err := api.ServiceErrors[api.ErrNotFound]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The section you are trying to update does not exist.",
		}

	}

	return sec, nil
}
