package service

import (
	"sort"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

type WarehouseService interface {
	Create(w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError)
	FindAll() ([]warehouse.Warehouse, *api.ServiceError)
	FindById(id int) (*warehouse.Warehouse, *api.ServiceError)
	Update(id int, patch warehouse.WarehousePatchDTO) (*warehouse.Warehouse, *api.ServiceError)
	Delete(id int) *api.ServiceError
}

type WarehouseDefault struct {
	rp repository.WarehouseRepository
}

func NewWarehouseDefault(rp repository.WarehouseRepository) *WarehouseDefault {
	return &WarehouseDefault{rp: rp}
}

func (s *WarehouseDefault) Create(w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError) {
	ok, err := s.rp.Exist(w.WarehouseCode)
	if err != nil {
		if err.Code != api.ErrNotFound {
			serviceErr := api.ServiceErrors[api.ErrInternalServer]
			serviceErr.Message = err.Message
			return nil, &serviceErr
		}
	}
	if ok {
		err := api.ServiceErrors[api.ErrConflict]
		return nil, &err
	}

	wh, err := s.rp.Create(w)
	if err != nil {
		err := api.ServiceErrors[api.ErrBadRequest]
		return nil, &err
	}
	return wh, nil
}

func (s *WarehouseDefault) FindAll() ([]warehouse.Warehouse, *api.ServiceError) {
	whs, err := s.rp.FindAll()
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

func (s *WarehouseDefault) FindById(id int) (*warehouse.Warehouse, *api.ServiceError) {
	w, err := s.rp.FindById(id)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WarehouseDefault) Update(id int, patch warehouse.WarehousePatchDTO) (*warehouse.Warehouse, *api.ServiceError) {
	existing, err := s.rp.FindById(id)
	if err != nil {
		return nil, err
	}

	if patch.Address != nil {
		existing.Address = *patch.Address
	}
	if patch.Telephone != nil {
		existing.Telephone = *patch.Telephone
	}
	if patch.WarehouseCode != nil {
		existing.WarehouseCode = *patch.WarehouseCode
	}
	if patch.MinimumCapacity != nil {
		existing.MinimumCapacity = *patch.MinimumCapacity
	}
	if patch.MinimumTemperature != nil {
		existing.MinimumTemperature = float64(*patch.MinimumTemperature)
	}
	updated, errRepo := s.rp.Update(id, *existing)
	if errRepo != nil {
		return nil, errRepo
	}
	return updated, nil
}

func (s *WarehouseDefault) Delete(id int) *api.ServiceError {
	_, err := s.rp.FindById(id)
	if err != nil {
		return err
	}
	return s.rp.Delete(id)
}
