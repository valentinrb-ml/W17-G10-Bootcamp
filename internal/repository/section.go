package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

const (
	querySectionGetAll = `SELECT id, section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id FROM sections `
	querySectionGetOne = `SELECT id, section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id FROM sections WHERE id = ?`
	querySectionDelete = `DELETE FROM sections WHERE id =?`
	querySectionUpdate = `UPDATE sections SET section_number = ?, current_capacity = ?, current_temperature = ? , maximum_capacity = ?, minimum_capacity = ?, minimum_temperature = ?, product_type_id = ?, warehouse_id = ?, updated_at = NOW() WHERE id = ?`
	querySectionCreate = `INSERT INTO sections (section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id) VALUES (?,?,?,?,?,?,?,?)`
)

// FindAllSections retrieves all Section records from the database.
// Returns a slice of Section or an error if the query fails.
func (r *sectionRepository) FindAllSections(ctx context.Context) ([]models.Section, error) {
	rows, err := r.mysql.QueryContext(ctx, querySectionGetAll)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the sections.")
	}
	defer rows.Close()

	var sections []models.Section

	for rows.Next() {
		var s models.Section
		if err := rows.Scan(&s.Id, &s.SectionNumber, &s.CurrentCapacity, &s.CurrentTemperature, &s.MaximumCapacity, &s.MinimumCapacity, &s.MinimumTemperature, &s.ProductTypeId, &s.WarehouseId); err != nil {
			return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the section.")
		}
		sections = append(sections, s)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the section.")
	}
	return sections, nil
}

// FindById retrieves a Section by its id.
// Returns a pointer to Section or error if not found.
func (r sectionRepository) FindById(ctx context.Context, id int) (*models.Section, error) {
	var s models.Section
	err := r.mysql.QueryRowContext(ctx, querySectionGetOne, id).Scan(&s.Id, &s.SectionNumber, &s.CurrentCapacity, &s.CurrentTemperature, &s.MaximumCapacity, &s.MinimumCapacity, &s.MinimumTemperature, &s.ProductTypeId, &s.WarehouseId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The section you are looking for does not exist.")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the section.")
	}

	return &s, nil
}

// DeleteSection deletes a Section by its id.
// Returns error if the section does not exist or if there are foreign key constraints.
func (r sectionRepository) DeleteSection(ctx context.Context, id int) error {
	result, err := r.mysql.ExecContext(ctx, querySectionDelete, id)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1451 {
			return apperrors.NewAppError(apperrors.CodeConflict, "Cannot delete section: there are products batches associated with this section.")
		}
		return apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while deleting the section.")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while deleting the section.")
	}
	if rows == 0 {
		return apperrors.NewAppError(apperrors.CodeNotFound, "The section you are trying to delete does not exist.")
	}
	return nil
}

// CreateSection inserts a new Section into the database, setting its Id.
// Returns pointer to the created Section or error if validation fails.
func (r *sectionRepository) CreateSection(ctx context.Context, sec models.Section) (*models.Section, error) {
	result, err := r.mysql.ExecContext(ctx, querySectionCreate,
		sec.SectionNumber, sec.CurrentCapacity, sec.CurrentTemperature, sec.MaximumCapacity, sec.MinimumCapacity, sec.MinimumTemperature, sec.ProductTypeId, sec.WarehouseId)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "Section number already exists.")
		}
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1452 {
			return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "Warehouse id or product type id does not exist.")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while creating the section.")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	sec.Id = int(id)

	return &sec, nil
}

// UpdateSection updates an existing Section by id with new data in sec.
// Returns pointer to updated Section or error if the section does not exist or constraints fail.
func (r *sectionRepository) UpdateSection(ctx context.Context, id int, sec *models.Section) (*models.Section, error) {
	result, err := r.mysql.ExecContext(ctx, querySectionUpdate,
		sec.SectionNumber, sec.CurrentCapacity, sec.CurrentTemperature,
		sec.MaximumCapacity, sec.MinimumCapacity, sec.MinimumTemperature,
		sec.ProductTypeId, sec.WarehouseId, id)
	if err != nil {
		//constraint UNIQUE
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "Section number already exists.")
		}
		// constraint FOREIGN KEY
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1452 {
			return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "Warehouse id or product type id does not exist.")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while updating the section.")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while updating the section.")
	}
	if rows == 0 {
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The section you are trying to update does not exist.")
	}

	return sec, nil
}
