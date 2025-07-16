package models

type PostSection struct {
	SectionNumber      int      `json:"section_number"`
	CurrentTemperature *float64 `json:"current_temperature"`
	MinimumTemperature *float64 `json:"minimum_temperature"`
	CurrentCapacity    int      `json:"current_capacity"`
	MinimumCapacity    int      `json:"minimum_capacity"`
	MaximumCapacity    int      `json:"maximum_capacity"`
	WarehouseId        int      `json:"warehouse_id"`
	ProductTypeId      int      `json:"product_type_id"`
}
type PatchSection struct {
	SectionNumber      *int     `json:"section_number"`
	CurrentTemperature *float64 `json:"current_temperature"`
	MinimumTemperature *float64 `json:"minimum_temperature"`
	CurrentCapacity    *int     `json:"current_capacity"`
	MinimumCapacity    *int     `json:"minimum_capacity"`
	MaximumCapacity    *int     `json:"maximum_capacity"`
	WarehouseId        *int     `json:"warehouse_id"`
	ProductTypeId      *int     `json:"product_type_id"`
}

type ResponseSection struct {
	Id                 int     `json:"id"`
	SectionNumber      int     `json:"section_number"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	WarehouseId        int     `json:"warehouse_id"`
	ProductTypeId      int     `json:"product_type_id"`
}

type Section struct {
	Id                 int     `json:"id"`
	SectionNumber      int     `json:"section_number"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	WarehouseId        int     `json:"warehouse_id"`
	ProductTypeId      int     `json:"product_type_id"`
}
