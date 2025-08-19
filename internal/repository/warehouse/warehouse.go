package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

// SQL queries for warehouse operations
const (
	queryWarehouseCreate   = `INSERT INTO warehouse (warehouse_code, address, minimum_temperature, minimum_capacity, telephone, locality_id) VALUES (?, ?, ?, ?, ?, ?)`
	queryWarehouseFindAll  = `SELECT id, warehouse_code, address, minimum_temperature, minimum_capacity, telephone, locality_id FROM warehouse ORDER BY id ASC`
	queryWarehouseFindById = `SELECT id, warehouse_code, address, minimum_temperature, minimum_capacity, telephone, locality_id FROM warehouse WHERE id = ?`
	queryWarehouseUpdate   = `UPDATE warehouse SET warehouse_code = ?, address = ?, minimum_temperature = ?, minimum_capacity = ?, telephone = ?, locality_id = ? WHERE id = ?`
	queryWarehouseDelete   = `DELETE FROM warehouse WHERE id = ?`
)

// Create inserts a new warehouse into the database
// Returns the created warehouse with its generated ID or an error if the operation fails
func (r *WarehouseMySQL) Create(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "warehouse-repository", "Creating new warehouse", map[string]interface{}{
			"warehouse_code": w.WarehouseCode,
			"locality_id":    w.LocalityId,
		})
	}

	res, err := r.db.ExecContext(ctx, queryWarehouseCreate, w.WarehouseCode, w.Address, w.MinimumTemperature, w.MinimumCapacity, w.Telephone, w.LocalityId)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				if r.logger != nil {
					r.logger.Warning(ctx, "warehouse-repository", "Duplicate warehouse code", map[string]interface{}{
						"warehouse_code": w.WarehouseCode,
						"mysql_error":    mysqlErr.Number,
					})
				}
				return nil, apperrors.NewAppError(apperrors.CodeConflict, "warehouse_code already exists")
			}
		}

		if r.logger != nil {
			r.logger.Error(ctx, "warehouse-repository", "Failed to create warehouse", err, map[string]interface{}{
				"warehouse_code": w.WarehouseCode,
			})
		}
		return nil, apperrors.Wrap(err, "error creating warehouse")
	}

	id, err := res.LastInsertId()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "warehouse-repository", "Failed to get last insert ID", err)
		}
		return nil, apperrors.Wrap(err, "error creating warehouse")
	}

	w.Id = int(id)

	if r.logger != nil {
		r.logger.Info(ctx, "warehouse-repository", "Warehouse created successfully", map[string]interface{}{
			"warehouse_id":   w.Id,
			"warehouse_code": w.WarehouseCode,
		})
	}

	return &w, nil
}

// FindAll retrieves all warehouses from the database
// Returns a slice of warehouses or an error if the operation fails
func (r *WarehouseMySQL) FindAll(ctx context.Context) ([]warehouse.Warehouse, error) {
	rows, err := r.db.QueryContext(ctx, queryWarehouseFindAll)
	if err != nil {
		return nil, apperrors.Wrap(err, "error getting warehouses")
	}
	defer rows.Close()

	var whs []warehouse.Warehouse
	for rows.Next() {
		var wh warehouse.Warehouse
		err := rows.Scan(&wh.Id, &wh.WarehouseCode, &wh.Address, &wh.MinimumTemperature, &wh.MinimumCapacity, &wh.Telephone, &wh.LocalityId)
		if err != nil {
			continue
		}
		whs = append(whs, wh)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, "error getting warehouses")
	}
	return whs, nil
}

// FindById retrieves a specific warehouse by its ID from the database
// Returns the warehouse if found or an error if not found or operation fails
func (r *WarehouseMySQL) FindById(ctx context.Context, id int) (*warehouse.Warehouse, error) {
	var w warehouse.Warehouse
	err := r.db.QueryRowContext(ctx, queryWarehouseFindById, id).Scan(
		&w.Id, &w.WarehouseCode, &w.Address, &w.MinimumTemperature, &w.MinimumCapacity, &w.Telephone, &w.LocalityId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found")
		}

		return nil, apperrors.Wrap(err, "error getting warehouse")
	}
	return &w, nil
}

// Update modifies an existing warehouse in the database
// Returns the updated warehouse or an error if the warehouse doesn't exist or operation fails
func (r *WarehouseMySQL) Update(ctx context.Context, id int, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
	res, err := r.db.ExecContext(ctx, queryWarehouseUpdate, w.WarehouseCode, w.Address, w.MinimumTemperature, w.MinimumCapacity, w.Telephone, w.LocalityId, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return nil, apperrors.NewAppError(apperrors.CodeConflict, "warehouse_code already exists")
			}
		}

		return nil, apperrors.Wrap(err, "error updating warehouse")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, apperrors.Wrap(err, "error updating warehouse")
	}
	if rowsAffected == 0 {
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found")
	}
	w.Id = id
	return &w, nil
}

// Delete removes a warehouse from the database by its ID
// Returns an error if the warehouse doesn't exist or operation fails
func (r *WarehouseMySQL) Delete(ctx context.Context, id int) error {
	if r.logger != nil {
		r.logger.Info(ctx, "warehouse-repository", "Deleting warehouse", map[string]interface{}{
			"warehouse_id": id,
		})
	}

	res, err := r.db.ExecContext(ctx, queryWarehouseDelete, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1451 {
				if r.logger != nil {
					r.logger.Warning(ctx, "warehouse-repository", "Cannot delete warehouse: foreign key constraint", map[string]interface{}{
						"warehouse_id": id,
						"mysql_error":  mysqlErr.Number,
					})
				}
				return apperrors.NewAppError(apperrors.CodeConflict, "cannot delete warehouse: it is being referenced by other records")
			}
		}

		if r.logger != nil {
			r.logger.Error(ctx, "warehouse-repository", "Failed to delete warehouse", err, map[string]interface{}{
				"warehouse_id": id,
			})
		}
		return apperrors.Wrap(err, "error deleting warehouse")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "warehouse-repository", "Failed to get rows affected", err, map[string]interface{}{
				"warehouse_id": id,
			})
		}
		return apperrors.Wrap(err, "error deleting warehouse")
	}

	if rowsAffected == 0 {
		if r.logger != nil {
			r.logger.Warning(ctx, "warehouse-repository", "Warehouse not found for deletion", map[string]interface{}{
				"warehouse_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found")
	}

	if r.logger != nil {
		r.logger.Info(ctx, "warehouse-repository", "Warehouse deleted successfully", map[string]interface{}{
			"warehouse_id": id,
		})
	}

	return nil
}
