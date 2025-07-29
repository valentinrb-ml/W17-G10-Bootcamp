// internal/models/purchase_order.go
package models

import "time"

type PurchaseOrder struct {
	ID              int       `json:"id"`
	OrderNumber     string    `json:"order_number"`
	OrderDate       time.Time `json:"order_date"`
	TrackingCode    string    `json:"tracking_code"`
	BuyerID         int       `json:"buyer_id"`
	ProductRecordID int       `json:"product_record_id"`
}

type RequestPurchaseOrder struct {
	OrderNumber     string `json:"order_number" validate:"required"`
	OrderDate       string `json:"order_date" validate:"required"`
	TrackingCode    string `json:"tracking_code" validate:"required"`
	BuyerID         int    `json:"buyer_id" validate:"required"`
	ProductRecordID int    `json:"product_record_id" validate:"required"`
}

type ResponsePurchaseOrder struct {
	ID              int    `json:"id"`
	OrderNumber     string `json:"order_number"`
	OrderDate       string `json:"order_date"`
	TrackingCode    string `json:"tracking_code"`
	BuyerID         int    `json:"buyer_id"`
	ProductRecordID int    `json:"product_record_id"`
}

type BuyerWithPurchaseCount struct {
	ID                  int    `json:"id"`
	CardNumberID        string `json:"id_card_number"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	PurchaseOrdersCount int    `json:"purchase_orders_count"`
}

// En tu archivo models/purchase_order.go
type PurchaseOrderRequestWrapper struct {
	Data RequestPurchaseOrder `json:"data"`
}
