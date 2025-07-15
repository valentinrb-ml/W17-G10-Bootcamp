package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)



const (
	queryWarehouseCreate   = `INSERT INTO warehouse (warehouse_code, address, minimum_temperature, minimum_capacity, telephone, locality_id) VALUES (?, ?, ?, ?, ?, ?)`
	queryWarehouseFindAll  = `SELECT id, warehouse_code, address, minimum_temperature, minimum_capacity, telephone, locality_id FROM warehouse`
	queryWarehouseFindById = `SELECT id, warehouse_code, address, minimum_temperature, minimum_capacity, telephone, locality_id FROM warehouse WHERE id = ?`
	queryWarehouseUpdate   = `UPDATE warehouse SET warehouse_code = ?, address = ?, minimum_temperature = ?, minimum_capacity = ?, telephone = ? WHERE id = ?`
	queryWarehouseDelete   = `DELETE FROM warehouse WHERE id = ?`
)

func (r *WarehouseMySQL) Create(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
	res, err := r.db.ExecContext(ctx, queryWarehouseCreate, w.WarehouseCode, w.Address, w.MinimumTemperature, w.MinimumCapacity, w.Telephone, w.LocalityId)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
        if mysqlErr.Number == 1062 {
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "warehouse_code already exists")
        }
    }
		
		return nil, apperrors.Wrap(err, "error creating warehouse")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, apperrors.Wrap(err, "error creating warehouse")
	}
	w.Id = int(id)
	return &w, nil
}

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

func (r *WarehouseMySQL) Update(ctx context.Context, id int, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
	res, err := r.db.ExecContext(ctx, queryWarehouseUpdate, w.WarehouseCode, w.Address, w.MinimumTemperature, w.MinimumCapacity, w.Telephone, id)
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

func (r *WarehouseMySQL) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, queryWarehouseDelete, id)
	if err != nil {
		return apperrors.Wrap(err, "error deleting warehouse")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return apperrors.Wrap(err, "error deleting warehouse")
	}
	if rowsAffected == 0 {
		return apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found")
	}
	return nil
}
