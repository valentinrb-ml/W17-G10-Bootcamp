package models

import "time"

// ProductRecordCore contains the base fields of the product record
type ProductRecordCore struct {
	LastUpdateDate time.Time `json:"last_update_date" db:"last_update_date"`
	PurchasePrice  float64   `json:"purchase_price" db:"purchase_price"`
	SalePrice      float64   `json:"sale_price" db:"sale_price"`
	ProductID      int       `json:"product_id" db:"product_id"`
}

// ProductRecord represents a product record in the domain
type ProductRecord struct {
	ID int `json:"id" db:"id"`
	ProductRecordCore
}

// ProductRecordRequest represents the request from the POST endpoint
type ProductRecordRequest struct {
	Data ProductRecordCore `json:"data"`
}

// ProductRecordResponse represents the response (same structure as the domain)
type ProductRecordResponse = ProductRecord

// ProductRecordReport represents the report of records by product
type ProductRecordReport struct {
	ProductID    int    `json:"product_id" db:"product_id"`
	Description  string `json:"description" db:"description"`
	RecordsCount int    `json:"records_count" db:"records_count"`
}

// ProductRecordsReportResponse represents the report response
type ProductRecordsReportResponse struct {
	Data []ProductRecordReport `json:"data"`
}
