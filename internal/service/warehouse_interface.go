package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

type WarehouseService interface {
	Create(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError)
	FindAll(ctx context.Context) ([]warehouse.Warehouse, *api.ServiceError)
	FindById(ctx context.Context, id int) (*warehouse.Warehouse, *api.ServiceError)
	Update(ctx context.Context, id int, patch warehouse.WarehousePatchDTO) (*warehouse.Warehouse, *api.ServiceError)
	Delete(ctx context.Context, id int) *api.ServiceError
}

type WarehouseDefault struct {
	rp repository.WarehouseRepository
}

func NewWarehouseService(rp repository.WarehouseRepository) *WarehouseDefault {
	return &WarehouseDefault{rp: rp}
}