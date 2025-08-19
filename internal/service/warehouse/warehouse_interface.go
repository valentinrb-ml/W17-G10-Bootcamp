package service

import (
	"context"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

type WarehouseService interface {
	Create(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error)
	FindAll(ctx context.Context) ([]warehouse.Warehouse, error)
	FindById(ctx context.Context, id int) (*warehouse.Warehouse, error)
	Update(ctx context.Context, id int, patch warehouse.WarehousePatchDTO) (*warehouse.Warehouse, error)
	Delete(ctx context.Context, id int) error
}

type WarehouseDefault struct {
	rp     repository.WarehouseRepository
	logger logger.Logger
}

func NewWarehouseService(rp repository.WarehouseRepository) *WarehouseDefault {
	return &WarehouseDefault{rp: rp}
}

// SetLogger allows you to inject the logger after creation
func (s *WarehouseDefault) SetLogger(l logger.Logger) {
	s.logger = l
}
