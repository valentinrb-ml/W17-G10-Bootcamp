package repository

import (
	"context"
	"database/sql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

type WarehouseRepository interface {
	Create(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error)
	FindAll(ctx context.Context) ([]warehouse.Warehouse, error)
	FindById(ctx context.Context, id int) (*warehouse.Warehouse, error)
	Update(ctx context.Context, id int, w warehouse.Warehouse) (*warehouse.Warehouse, error)
	Delete(ctx context.Context, id int) error
}

type WarehouseMySQL struct {
	db *sql.DB
}

func NewWarehouseRepository(db *sql.DB) *WarehouseMySQL {
	return &WarehouseMySQL{db}
}
