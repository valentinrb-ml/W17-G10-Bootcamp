package testhelpers

import (
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
	"time"
)

func DummyProductBatch(id int) models.ProductBatches {
	return models.ProductBatches{
		Id:                 id,
		BatchNumber:        100 + id,
		CurrentQuantity:    50,
		CurrentTemperature: 7,
		DueDate:            time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		InitialQuantity:    100,
		ManufacturingDate:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		ManufacturingHour:  10,
		MinimumTemperature: 1,
		ProductId:          22,
		SectionId:          33,
	}
}

func DummyReportProduct() models.ReportProduct {
	return models.ReportProduct{
		SectionId:     10,
		SectionNumber: 5,
		ProductsCount: 123,
	}
}

func DummyReportProductsList() []models.ReportProduct {
	return []models.ReportProduct{
		{
			SectionId:     10,
			SectionNumber: 5,
			ProductsCount: 123,
		},
		{
			SectionId:     20,
			SectionNumber: 7,
			ProductsCount: 42,
		},
	}
}
