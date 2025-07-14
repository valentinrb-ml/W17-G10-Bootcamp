package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)



const (
	queryWarehouseExist    = `SELECT COUNT(*) FROM warehouse WHERE warehouse_code = ?`
	queryWarehouseCreate   = `INSERT INTO warehouse (warehouse_code, address, minimum_temperature, minimum_capacity, telephone, locality_id) VALUES (?, ?, ?, ?, ?, ?)`
	queryWarehouseFindAll  = `SELECT id, warehouse_code, address, minimum_temperature, minimum_capacity, telephone FROM warehouse`
	queryWarehouseFindById = `SELECT id, warehouse_code, address, minimum_temperature, minimum_capacity, telephone FROM warehouse WHERE id = ?`
	queryWarehouseUpdate   = `UPDATE warehouse SET warehouse_code = ?, address = ?, minimum_temperature = ?, minimum_capacity = ?, telephone = ? WHERE id = ?`
	queryWarehouseDelete   = `DELETE FROM warehouse WHERE id = ?`
)

func (r *WarehouseMySQL) Exist(ctx context.Context, wc string) (bool, *api.ServiceError) {
	var count int
	err := r.db.QueryRowContext(ctx, queryWarehouseExist, wc).Scan(&count)
	if err != nil {
		errVal := api.ServiceErrors[api.ErrInternalServer]
		errVal.InternalError = err
		return false, &errVal
	}

	// Si count > 0, existe; si count = 0, no existe (sin error)
	return count > 0, nil
}

func (r *WarehouseMySQL) Create(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError) {
	res, err := r.db.ExecContext(ctx, queryWarehouseCreate, w.WarehouseCode, w.Address, w.MinimumTemperature, w.MinimumCapacity, w.Telephone, w.LocalityId)
	if err != nil {
		errVal := api.ServiceErrors[api.ErrInternalServer]
		errVal.InternalError = err
		return nil, &errVal
	}
	id, err := res.LastInsertId()
	if err != nil {
		errVal := api.ServiceErrors[api.ErrInternalServer]
		errVal.InternalError = err
		return nil, &errVal
	}
	w.Id = int(id)
	return &w, nil
}

func (r *WarehouseMySQL) FindAll(ctx context.Context) ([]warehouse.Warehouse, *api.ServiceError) {
	rows, err := r.db.QueryContext(ctx, queryWarehouseFindAll)
	if err != nil {
		errVal := api.ServiceErrors[api.ErrInternalServer]
		errVal.InternalError = err
		return nil, &errVal
	}
	defer rows.Close()

	var whs []warehouse.Warehouse
	for rows.Next() {
		var wh warehouse.Warehouse
		err := rows.Scan(&wh.Id, &wh.WarehouseCode, &wh.Address, &wh.MinimumTemperature, &wh.MinimumCapacity, &wh.Telephone)
		if err != nil {
			// Log the error but continue processing other rows
			continue
		}
		whs = append(whs, wh)
	}

	if err := rows.Err(); err != nil {
		errVal := api.ServiceErrors[api.ErrInternalServer]
		errVal.InternalError = err
		return nil, &errVal
	}
	return whs, nil
}

func (r *WarehouseMySQL) FindById(ctx context.Context, id int) (*warehouse.Warehouse, *api.ServiceError) {
	var w warehouse.Warehouse
	err := r.db.QueryRowContext(ctx, queryWarehouseFindById, id).Scan(
		&w.Id, &w.WarehouseCode, &w.Address, &w.MinimumTemperature, &w.MinimumCapacity, &w.Telephone,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errVal := api.ServiceErrors[api.ErrNotFound]
			return nil, &errVal
		}
		errVal := api.ServiceErrors[api.ErrInternalServer]
		errVal.InternalError = err
		return nil, &errVal
	}
	return &w, nil
}

func (r *WarehouseMySQL) Update(ctx context.Context, id int, w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError) {
	res, err := r.db.ExecContext(ctx, queryWarehouseUpdate, w.WarehouseCode, w.Address, w.MinimumTemperature, w.MinimumCapacity, w.Telephone, id)
	if err != nil {
		errVal := api.ServiceErrors[api.ErrInternalServer]
		errVal.InternalError = err
		return nil, &errVal
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		errVal := api.ServiceErrors[api.ErrInternalServer]
		errVal.InternalError = err
		return nil, &errVal
	}
	if rowsAffected == 0 {
		errVal := api.ServiceErrors[api.ErrNotFound]
		return nil, &errVal
	}
	w.Id = id
	return &w, nil
}

func (r *WarehouseMySQL) Delete(ctx context.Context, id int) *api.ServiceError {
	res, err := r.db.ExecContext(ctx, queryWarehouseDelete, id)
	if err != nil {
		errVal := api.ServiceErrors[api.ErrInternalServer]
		errVal.InternalError = err
		return &errVal
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		errVal := api.ServiceErrors[api.ErrInternalServer]
		errVal.InternalError = err
		return &errVal
	}
	if rowsAffected == 0 {
		errVal := api.ServiceErrors[api.ErrNotFound]
		return &errVal
	}
	return nil
}
