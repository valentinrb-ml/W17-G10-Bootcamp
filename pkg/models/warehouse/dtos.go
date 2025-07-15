package warehouse

type WarehouseDoc struct {
	ID                 int     `json:"id,omitempty"`
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	WarehouseCode      string  `json:"warehouse_code"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	LocalityId         int     `json:"locality_id"`
}

type WarehouseRequest struct {
	Id                 int      `json:"id,omitempty"`
	Address            string   `json:"address,omitempty"`
	Telephone          string   `json:"telephone,omitempty"`
	WarehouseCode      string   `json:"warehouse_code,omitempty"`
	MinimumCapacity    int      `json:"minimum_capacity,omitempty"`
	MinimumTemperature *float64 `json:"minimum_temperature,omitempty"`
	LocalityId         int      `json:"locality_id"`
}

type WarehousePatchDTO struct {
	Address            *string  `json:"address,omitempty"`
	Telephone          *string  `json:"telephone,omitempty"`
	WarehouseCode      *string  `json:"warehouse_code,omitempty"`
	MinimumCapacity    *int     `json:"minimum_capacity,omitempty"`
	MinimumTemperature *float64 `json:"minimum_temperature,omitempty"`
}
