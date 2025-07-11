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

func RequestToWarehouse(req warehouse.WarehouseRequest) warehouse.Warehouse{
	return warehouse.Warehouse{
		Id:                 0,
		Address:            req.Address,
		WarehouseCode:      req.WarehouseCode,
		Telephone:          req.Telephone,
		MinimumCapacity:    req.MinimumCapacity,
		MinimumTemperature: *req.MinimumTemperature,
		LocalityId: 	   req.LocalityId,
	}
}

func ApplyWarehousePatch(existing *warehouse.Warehouse, patch warehouse.WarehousePatchDTO) {
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
}
