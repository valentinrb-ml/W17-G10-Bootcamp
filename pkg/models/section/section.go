package section

type RequestSection struct {
	SectionNumber      *int `json:"section_number"`
	CurrentTemperature *int `json:"current_temperature"`
	MinimumTemperature *int `json:"minimum_temperature"`
	CurrentCapacity    *int `json:"current_capacity"`
	MinimumCapacity    *int `json:"minimum_capacity"`
	MaximumCapacity    *int `json:"maximum_capacity"`
	WarehouseId        *int `json:"warehouse_id"`
	ProductId          int  `json:"product_id"`
}

type ResponseSection struct {
	Id                 int `json:"id"`
	SectionNumber      int `json:"section_number"`
	CurrentTemperature int `json:"current_temperature"`
	MinimumTemperature int `json:"minimum_temperature"`
	CurrentCapacity    int `json:"current_capacity"`
	MinimumCapacity    int `json:"minimum_capacity"`
	MaximumCapacity    int `json:"maximum_capacity"`
	WarehouseId        int `json:"warehouse_id"`
	ProductId          int `json:"product_id"`
}

type Section struct {
	Id                 int `json:"id"`
	SectionNumber      int `json:"section_number"`
	CurrentTemperature int `json:"current_temperature"`
	MinimumTemperature int `json:"minimum_temperature"`
	CurrentCapacity    int `json:"current_capacity"`
	MinimumCapacity    int `json:"minimum_capacity"`
	MaximumCapacity    int `json:"maximum_capacity"`
	WarehouseId        int `json:"warehouse_id"`
	ProductId          int `json:"product_id"`
}
