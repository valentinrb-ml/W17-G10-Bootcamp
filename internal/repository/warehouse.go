package repository

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

type WarehouseRepository interface {
	Create(w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError)
	Exist(wc string) (bool, *api.ServiceError)
	FindAll() ([]warehouse.Warehouse, *api.ServiceError)
	FindById(id int) (*warehouse.Warehouse, *api.ServiceError)
	Update(id int, w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError)
	Delete(id int) *api.ServiceError
}

type WarehouseMap struct {
	db map[int]warehouse.Warehouse
}

func NewWarehouseMap(db map[int]warehouse.Warehouse) *WarehouseMap {
	defaultDb := make(map[int]warehouse.Warehouse)
	if db != nil {
		defaultDb = db
	}
	return &WarehouseMap{db: defaultDb}
}

func (r *WarehouseMap) Exist(wc string) (bool, *api.ServiceError) {
	for _, warehouse := range r.db {
		if warehouse.WarehouseCode == wc {
			return true, nil
		}
	}
	err := api.ServiceErrors[api.ErrNotFound]
	return false, &err
}

func (r *WarehouseMap) Create(w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError) {
	maxId := 0
	for _, wh := range r.db {
		if wh.Id > maxId {
			maxId = wh.Id
		}
	}

	w.Id = maxId + 1
	r.db[w.Id] = w
	return &w, nil
}

func (r *WarehouseMap) FindAll() ([]warehouse.Warehouse, *api.ServiceError) {
	w := make([]warehouse.Warehouse, 0, len(r.db))

	for _, wh := range r.db {
		w = append(w, wh)
	}
	return w, nil
}

func (r *WarehouseMap) FindById(id int) (*warehouse.Warehouse, *api.ServiceError) {
	wh, ok := r.db[id]
	if !ok {
		err := api.ServiceErrors[api.ErrNotFound]
		return nil, &err
	}
	return &wh, nil
}

func (r *WarehouseMap) Update(id int, w warehouse.Warehouse) (*warehouse.Warehouse, *api.ServiceError) {
	_, ok := r.db[id]
	if !ok {
		err := api.ServiceErrors[api.ErrNotFound]
		return nil, &err
	}
	r.db[id] = w
	return &w, nil
}

func (r *WarehouseMap) Delete(id int) *api.ServiceError {
	_, ok := r.db[id]
	if !ok {
		err := api.ServiceErrors[api.ErrNotFound]
		return &err
	}
	delete(r.db, id)
	return nil
}
