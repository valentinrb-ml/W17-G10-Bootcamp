package models

import "time"

type ProductBatches struct {
	Id                 int
	BatchNumber        int
	CurrentQuantity    int
	CurrentTemperature float64
	DueDate            time.Time
	InitialQuantity    int
	ManufacturingDate  time.Time
	ManufacturingHour  int
	MinimumTemperature float64
	ProductId          int
	SectionId          int
}

type ProductBatchesResponse struct {
	Id                 int       `json:"id"`
	BatchNumber        int       `json:"batch_number"`
	CurrentQuantity    int       `json:"current_quantity"`
	CurrentTemperature float64   `json:"current_temperature"`
	DueDate            time.Time `json:"due_date"`
	InitialQuantity    int       `json:"initial_quantity"`
	ManufacturingDate  time.Time `json:"manufacturing_date"`
	ManufacturingHour  int       `json:"manufacturing_hour"`
	MinimumTemperature float64   `json:"minimum_temperature"`
	ProductId          int       `json:"product_id"`
	SectionId          int       `json:"section_id"`
}

type PostProductBatches struct {
	BatchNumber        int       `json:"batch_number"`
	CurrentQuantity    *int      `json:"current_quantity"`
	CurrentTemperature *float64  `json:"current_temperature"`
	DueDate            time.Time `json:"due_date"`
	InitialQuantity    *int      `json:"initial_quantity"`
	ManufacturingDate  time.Time `json:"manufacturing_date"`
	ManufacturingHour  int       `json:"manufacturing_hour"`
	MinimumTemperature *float64  `json:"minimum_temperature"`
	ProductId          int       `json:"product_id"`
	SectionId          int       `json:"section_id"`
}

type ReportProduct struct {
	SectionId     int `json:"section_id"`
	SectionNumber int `json:"section_number"`
	ProductsCount int `json:"products_count"`
}
