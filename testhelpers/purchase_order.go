package testhelpers

import (
	"time"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

// PurchaseOrderDummyMap contiene datos dummy de órdenes de compra
var PurchaseOrderDummyMap = map[int]models.PurchaseOrder{
	1: {
		ID:              1,
		OrderNumber:     "PO-001",
		OrderDate:       time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		TrackingCode:    "TRACK001",
		BuyerID:         101,
		ProductRecordID: 201,
	},
	2: {
		ID:              2,
		OrderNumber:     "PO-002",
		OrderDate:       time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
		TrackingCode:    "TRACK002",
		BuyerID:         102,
		ProductRecordID: 202,
	},
}

// BuyerWithPurchaseCountDummyMap contiene datos dummy de compradores con conteo
var BuyerWithPurchaseCountDummyMap = map[int]models.BuyerWithPurchaseCount{
	101: {
		ID:                  101,
		CardNumberID:        "CARD101",
		FirstName:           "John",
		LastName:            "Doe",
		PurchaseOrdersCount: 3,
	},
	102: {
		ID:                  102,
		CardNumberID:        "CARD102",
		FirstName:           "Jane",
		LastName:            "Smith",
		PurchaseOrdersCount: 5,
	},
}

// CreateTestPurchaseOrder crea una orden de compra de prueba
// CreateTestPurchaseOrder crea una orden de compra de prueba
func CreateTestPurchaseOrder(id int) models.PurchaseOrder {
	return models.PurchaseOrder{
		ID:              id,
		OrderNumber:     "TEST-PO",
		OrderDate:       time.Now(),
		TrackingCode:    "TEST-TRACK",
		BuyerID:         100,
		ProductRecordID: 200,
	}
}

// DummyPurchaseOrderRequest crea una solicitud dummy de orden de compra
func DummyPurchaseOrderRequest() models.PurchaseOrder {
	return models.PurchaseOrder{
		OrderNumber:     "PO-999",
		OrderDate:       time.Now(),
		TrackingCode:    "TRACK999",
		BuyerID:         999,
		ProductRecordID: 999,
	}
}

// DummyBuyerWithPurchaseCount crea un comprador dummy con conteo de compras
func DummyBuyerWithPurchaseCount() models.BuyerWithPurchaseCount {
	return models.BuyerWithPurchaseCount{
		ID:                  100,
		CardNumberID:        "CARD100",
		FirstName:           "Test",
		LastName:            "Buyer",
		PurchaseOrdersCount: 10,
	}
}
