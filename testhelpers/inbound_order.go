package testhelpers

import (
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

// Devuelve un inbound_order típico
func CreateTestInboundOrder() models.InboundOrder {
	return models.InboundOrder{
		OrderNumber:    "INV001",
		OrderDate:      "2024-06-01",
		EmployeeID:     1,
		ProductBatchID: 10,
		WarehouseID:    1,
	}
}

func CreateExpectedInboundOrder(id int) *models.InboundOrder {
	o := CreateTestInboundOrder()
	o.ID = id
	return &o
}

func CreateInboundOrderReport(id int) models.InboundOrderReport {
	return models.InboundOrderReport{
		ID:                 id,
		CardNumberID:       "CARDID",
		FirstName:          "Juan",
		LastName:           "Tester",
		WarehouseID:        1,
		InboundOrdersCount: 5,
	}
}

func CreateInboundOrderReports() []models.InboundOrderReport {
	return []models.InboundOrderReport{
		CreateInboundOrderReport(1),
		CreateInboundOrderReport(2),
	}
}
