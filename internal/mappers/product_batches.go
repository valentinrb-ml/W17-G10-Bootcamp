package mappers

import (
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

func RequestToProductBatch(req models.PostProductBatches) models.ProductBatches {
	return models.ProductBatches{
		BatchNumber:        req.BatchNumber,
		CurrentQuantity:    *req.CurrentQuantity,
		CurrentTemperature: *req.CurrentTemperature,
		DueDate:            req.DueDate,
		InitialQuantity:    *req.InitialQuantity,
		ManufacturingDate:  req.ManufacturingDate,
		ManufacturingHour:  req.ManufacturingHour,
		MinimumTemperature: *req.MinimumTemperature,
		ProductId:          req.ProductId,
		SectionId:          req.SectionId,
	}
}

func ProductBatchesToResponse(proBa models.ProductBatches) models.ProductBatchesResponse {
	return models.ProductBatchesResponse{
		Id:                 proBa.Id,
		BatchNumber:        proBa.BatchNumber,
		CurrentQuantity:    proBa.CurrentQuantity,
		CurrentTemperature: proBa.CurrentTemperature,
		DueDate:            proBa.DueDate,
		InitialQuantity:    proBa.InitialQuantity,
		ManufacturingDate:  proBa.ManufacturingDate,
		ManufacturingHour:  proBa.ManufacturingHour,
		MinimumTemperature: proBa.MinimumTemperature,
		ProductId:          proBa.ProductId,
		SectionId:          proBa.SectionId,
	}
}
