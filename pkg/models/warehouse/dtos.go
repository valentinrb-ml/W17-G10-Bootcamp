package warehouse

type WarehouseDoc struct {
	ID                 int     `json:"id,omitempty"`
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	WarehouseCode      string  `json:"warehouse_code"`
	MinimumCapacity    int     `json:"minimun_capacity"`
	MinimumTemperature float64 `json:"minimun_temperature"`
}

type WarehouseRequest struct {
	Id                 int      `json:"id,omitempty"`
	Address            string   `json:"address,omitempty"`
	Telephone          string   `json:"telephone,omitempty"`
	WarehouseCode      string   `json:"warehouse_code,omitempty"`
	MinimumCapacity    int      `json:"minimum_capacity,omitempty"`
	MinimumTemperature *float64 `json:"minimum_temperature,omitempty"`
}

type WarehousePatchDTO struct {
	Address            *string 		`json:"address,omitempty"`
	Telephone          *string 		`json:"telephone,omitempty"`
	WarehouseCode      *string 		`json:"warehouse_code,omitempty"`
	MinimumCapacity    *int    		`json:"minimum_capacity,omitempty"`
	MinimumTemperature *float64    	`json:"minimum_temperature,omitempty"`
}
