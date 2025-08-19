package mappers_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

func TestRequestPurchaseOrderToPurchaseOrder_OK(t *testing.T) {
	req := models.RequestPurchaseOrder{
		OrderNumber:     "ORD-42",
		OrderDate:       "2023-07-15",
		TrackingCode:    "TRK-123",
		BuyerID:         1,
		ProductRecordID: 2,
	}
	po, err := mappers.RequestPurchaseOrderToPurchaseOrder(req)
	assert.NoError(t, err)
	assert.Equal(t, "ORD-42", po.OrderNumber)
	assert.True(t, po.OrderDate.Equal(time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)), "time parse ok")
	assert.Equal(t, "TRK-123", po.TrackingCode)
	assert.Equal(t, 1, po.BuyerID)
	assert.Equal(t, 2, po.ProductRecordID)
}

func TestRequestPurchaseOrderToPurchaseOrder_ParseError(t *testing.T) {
	req := models.RequestPurchaseOrder{
		OrderNumber:     "ORD-42",
		OrderDate:       "not-a-date",
		TrackingCode:    "TRK-123",
		BuyerID:         1,
		ProductRecordID: 2,
	}
	po, err := mappers.RequestPurchaseOrderToPurchaseOrder(req)
	assert.Error(t, err)
	assert.Equal(t, models.PurchaseOrder{}, po)
}

func TestPurchaseOrderToResponse(t *testing.T) {
	orderDate := time.Date(2022, 11, 25, 0, 0, 0, 0, time.UTC)
	po := models.PurchaseOrder{
		ID:              99,
		OrderNumber:     "ONO",
		OrderDate:       orderDate,
		TrackingCode:    "TRK",
		BuyerID:         8,
		ProductRecordID: 20,
	}
	resp := mappers.PurchaseOrderToResponse(po)
	assert.Equal(t, 99, resp.ID)
	assert.Equal(t, "ONO", resp.OrderNumber)
	assert.Equal(t, "2022-11-25", resp.OrderDate)
	assert.Equal(t, "TRK", resp.TrackingCode)
	assert.Equal(t, 8, resp.BuyerID)
	assert.Equal(t, 20, resp.ProductRecordID)
}

func TestToResponsePurchaseOrderList(t *testing.T) {
	orders := []models.PurchaseOrder{
		{
			ID: 1, OrderNumber: "ORD1", OrderDate: time.Date(2021, 3, 4, 0, 0, 0, 0, time.UTC),
			TrackingCode: "T1", BuyerID: 10, ProductRecordID: 99,
		},
		{
			ID: 2, OrderNumber: "ORD2", OrderDate: time.Date(2022, 4, 5, 0, 0, 0, 0, time.UTC),
			TrackingCode: "T2", BuyerID: 20, ProductRecordID: 100,
		},
	}
	resp := mappers.ToResponsePurchaseOrderList(orders)
	assert.Len(t, resp, 2)
	assert.Equal(t, "ORD1", resp[0].OrderNumber)
	assert.Equal(t, "2021-03-04", resp[0].OrderDate)
	assert.Equal(t, 10, resp[0].BuyerID)
	assert.Equal(t, 99, resp[0].ProductRecordID)
	assert.Equal(t, "ORD2", resp[1].OrderNumber)
	assert.Equal(t, "2022-04-05", resp[1].OrderDate)
	assert.Equal(t, 20, resp[1].BuyerID)
	assert.Equal(t, 100, resp[1].ProductRecordID)
}
