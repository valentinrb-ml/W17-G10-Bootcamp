package mappers

import "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"

func WarehouseToDoc(w *warehouse.Warehouse) warehouse.WarehouseDoc{
	return warehouse.WarehouseDoc{
		ID: w.Id,
		WarehouseCode: w.WarehouseCode,
		Address: w.Address,
		Telephone: w.Telephone,
		MinimumCapacity: w.MinimumCapacity,
		MinimumTemperature: w.MinimumTemperature,
	}
}

func WarehouseToDocSlice(w []warehouse.Warehouse) []warehouse.WarehouseDoc{
	newWarehouses := make([]warehouse.WarehouseDoc, 0, len(w))
	for _, wh := range w {
		wDoc := WarehouseToDoc(&wh)
		newWarehouses = append(newWarehouses, wDoc)
	}
	return newWarehouses
}