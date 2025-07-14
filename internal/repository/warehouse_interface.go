package repository

import (
	"context"
	"database/sql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

type WarehouseRepository interface {
	Create(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError)
	Exist(ctx context.Context, wc string) (bool, *api.ServiceError)
	FindAll(ctx context.Context) ([]warehouse.Warehouse, *api.ServiceError)
	FindById(ctx context.Context, id int) (*warehouse.Warehouse, *api.ServiceError)
	Update(ctx context.Context, id int, w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError)
	Delete(ctx context.Context, id int) *api.ServiceError
}

type WarehouseMySQL struct {
	db *sql.DB
}

func NewWarehouseRepository(db *sql.DB) *WarehouseMySQL {
	return &WarehouseMySQL{db}
}