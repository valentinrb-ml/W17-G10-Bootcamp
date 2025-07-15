package mappers

import (
	"time"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

/*
*

	func RequestPurchaseOrderToPurchaseOrder(req models.RequestPurchaseOrder) (models.PurchaseOrder, error) {
		orderDate, err := time.Parse("2006-01-02", req.OrderDate)
		if err != nil {
			return models.PurchaseOrder{}, err
		}

*
*/
func RequestPurchaseOrderToPurchaseOrder(req models.RequestPurchaseOrder) (models.PurchaseOrder, error) {
	orderDate, err := time.Parse("2006-01-02 15:04:05", req.OrderDate)
	if err != nil {
		return models.PurchaseOrder{}, err
	}

	return models.PurchaseOrder{
		OrderNumber:     req.OrderNumber,
		OrderDate:       orderDate,
		TrackingCode:    req.TrackingCode,
		BuyerID:         req.BuyerID,
		ProductRecordID: req.ProductRecordID,
	}, nil
}

func PurchaseOrderToResponse(po models.PurchaseOrder) models.ResponsePurchaseOrder {
	return models.ResponsePurchaseOrder{
		ID:              po.ID,
		OrderNumber:     po.OrderNumber,
		OrderDate:       po.OrderDate.Format("2006-01-02"),
		TrackingCode:    po.TrackingCode,
		BuyerID:         po.BuyerID,
		ProductRecordID: po.ProductRecordID,
	}
}

func ToResponsePurchaseOrderList(pos []models.PurchaseOrder) []models.ResponsePurchaseOrder {
	res := make([]models.ResponsePurchaseOrder, 0, len(pos))
	for _, po := range pos {
		res = append(res, PurchaseOrderToResponse(po))
	}
	return res
}
