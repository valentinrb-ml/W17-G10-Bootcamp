package service

import (
	"context"
	"sort"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
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

func NewWarehouseDefault(rp repository.WarehouseRepository) *WarehouseDefault {
	return &WarehouseDefault{rp: rp}
}

func (s *WarehouseDefault) Create(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError) {
	ok, err := s.rp.Exist(ctx, w.WarehouseCode)
	if err != nil {
		serviceErr := api.ServiceErrors[api.ErrInternalServer]
		serviceErr.Message = err.Message
		return nil, &serviceErr
	}
	if ok {
		err := api.ServiceErrors[api.ErrConflict]
		err.Message = "warehouse_code already exists"
		return nil, &err
	}

	wh, err := s.rp.Create(ctx, w)
	if err != nil {
		return nil, err
	}
	return wh, nil
}

func (s *WarehouseDefault) FindAll(ctx context.Context) ([]warehouse.Warehouse, *api.ServiceError) {
	whs, err := s.rp.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(whs) == 0 {
		return whs, nil
	}

	sort.Slice(whs, func(i, j int) bool {
		return whs[i].Id < whs[j].Id
	})

	return whs, err
}

func (s *WarehouseDefault) FindById(ctx context.Context, id int) (*warehouse.Warehouse, *api.ServiceError) {
	w, err := s.rp.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WarehouseDefault) Update(ctx context.Context, id int, patch warehouse.WarehousePatchDTO) (*warehouse.Warehouse, *api.ServiceError) {
	existing, err := s.rp.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if patch.MinimumCapacity != nil {
		if err := validators.ValidateMinimumCapacity(*patch.MinimumCapacity); err != nil {
			return nil, err
		}
	}

	if patch.WarehouseCode != nil {
		ok, err := s.rp.Exist(ctx, *patch.WarehouseCode)
		if err != nil {
			return nil, err
		}
		if ok {
			err := api.ServiceErrors[api.ErrConflict]
			err.Message = "warehouse_code already exists"
			return nil, &err
		}
	}

	mappers.ApplyWarehousePatch(existing, patch)

	updated, errRepo := s.rp.Update(ctx, id, *existing)
	if errRepo != nil {
		return nil, errRepo
	}
	return updated, nil
}

func (s *WarehouseDefault) Delete(ctx context.Context, id int) *api.ServiceError {
	_, err := s.rp.FindById(ctx, id)
	if err != nil {
		return err
	}
	return s.rp.Delete(ctx, id)
}
