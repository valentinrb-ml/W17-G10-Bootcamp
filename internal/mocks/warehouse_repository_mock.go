package mocks

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

type WarehouseRepositoryMock struct {
    FuncCreate   func(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error)
    FuncFindAll  func(ctx context.Context) ([]warehouse.Warehouse, error)
    FuncFindById func(ctx context.Context, id int) (*warehouse.Warehouse, error)
    FuncUpdate   func(ctx context.Context, id int, w warehouse.Warehouse) (*warehouse.Warehouse, error)
    FuncDelete   func(ctx context.Context, id int) error
}

func NewWarehouseRepositoryMock() *WarehouseRepositoryMock {
    return &WarehouseRepositoryMock{}
}

func (m *WarehouseRepositoryMock) Create(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
    return m.FuncCreate(ctx, w)
}

func (m *WarehouseRepositoryMock) FindAll(ctx context.Context) ([]warehouse.Warehouse, error) {
    return m.FuncFindAll(ctx)
}

func (m *WarehouseRepositoryMock) FindById(ctx context.Context, id int) (*warehouse.Warehouse, error) {
    return m.FuncFindById(ctx, id)
}

func (m *WarehouseRepositoryMock) Update(ctx context.Context, id int, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
    return m.FuncUpdate(ctx, id, w)
}

func (m *WarehouseRepositoryMock) Delete(ctx context.Context, id int) error {
    return m.FuncDelete(ctx, id)
}
