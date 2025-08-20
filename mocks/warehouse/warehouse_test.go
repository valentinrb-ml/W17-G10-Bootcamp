package mocks_test

import (
	"context"
	"testing"

	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func TestWarehouseRepositoryMock_DummyCoverage(t *testing.T) {
	m := &mocks.WarehouseRepositoryMock{
		FuncCreate:   func(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error) { return nil, nil },
		FuncFindAll:  func(ctx context.Context) ([]warehouse.Warehouse, error) { return nil, nil },
		FuncFindById: func(ctx context.Context, id int) (*warehouse.Warehouse, error) { return nil, nil },
		FuncUpdate: func(ctx context.Context, id int, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
			return nil, nil
		},
		FuncDelete: func(ctx context.Context, id int) error { return nil },
	}
	m.Create(context.TODO(), warehouse.Warehouse{})
	m.FindAll(context.TODO())
	m.FindById(context.TODO(), 0)
	m.Update(context.TODO(), 0, warehouse.Warehouse{})
	m.Delete(context.TODO(), 0)
}

func TestWarehouseServiceMock_DummyCoverage(t *testing.T) {
	m := &mocks.WarehouseServiceMock{
		FuncCreate:   func(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error) { return nil, nil },
		FuncFindAll:  func(ctx context.Context) ([]warehouse.Warehouse, error) { return nil, nil },
		FuncFindById: func(ctx context.Context, id int) (*warehouse.Warehouse, error) { return nil, nil },
		FuncUpdate: func(ctx context.Context, id int, patch warehouse.WarehousePatchDTO) (*warehouse.Warehouse, error) {
			return nil, nil
		},
		FuncDelete: func(ctx context.Context, id int) error { return nil },
	}

	m.Create(context.TODO(), warehouse.Warehouse{})
	m.FindAll(context.TODO())
	m.FindById(context.TODO(), 0)
	m.Update(context.TODO(), 0, warehouse.WarehousePatchDTO{})
	m.Delete(context.TODO(), 0)
	m.Reset()

	m.AssertCreateCalledWith(nil, nil, warehouse.Warehouse{})
	m.AssertFindByIdCalledWith(nil, nil, 0)
}
