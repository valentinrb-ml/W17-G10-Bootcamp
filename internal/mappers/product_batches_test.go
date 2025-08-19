package mappers_test

import (
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
	"testing"
)

func TestRequestToProductBatch(t *testing.T) {
	t.Run("maps all fields from PostProductBatches to ProductBatches", func(t *testing.T) {
		req := testhelpers.DummyPostProductBatch(1)
		expected := models.ProductBatches{
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

		got := mappers.RequestToProductBatch(req)
		require.Equal(t, expected, got)
	})

	t.Run("panics if CurrentQuantity is nil", func(t *testing.T) {
		req := testhelpers.DummyPostProductBatch(1)
		req.CurrentQuantity = nil
		require.Panics(t, func() {
			_ = mappers.RequestToProductBatch(req)
		})
	})

	t.Run("panics if CurrentTemperature is nil", func(t *testing.T) {
		req := testhelpers.DummyPostProductBatch(1)
		req.CurrentTemperature = nil
		require.Panics(t, func() {
			_ = mappers.RequestToProductBatch(req)
		})
	})

	t.Run("panics if InitialQuantity is nil", func(t *testing.T) {
		req := testhelpers.DummyPostProductBatch(1)
		req.InitialQuantity = nil
		require.Panics(t, func() {
			_ = mappers.RequestToProductBatch(req)
		})
	})

	t.Run("panics if MinimumTemperature is nil", func(t *testing.T) {
		req := testhelpers.DummyPostProductBatch(1)
		req.MinimumTemperature = nil
		require.Panics(t, func() {
			_ = mappers.RequestToProductBatch(req)
		})
	})
}

func TestProductBatchesToResponse(t *testing.T) {
	t.Run("maps all fields from ProductBatches to ProductBatchesResponse", func(t *testing.T) {
		input := testhelpers.DummyProductBatch(1)
		expected := testhelpers.DummyResponseProductBatch(1) // Si tienes un helper para esto

		got := mappers.ProductBatchesToResponse(input)
		require.Equal(t, expected, got)
	})

	t.Run("maps zero values as is", func(t *testing.T) {
		input := models.ProductBatches{}
		expected := models.ProductBatchesResponse{}
		got := mappers.ProductBatchesToResponse(input)
		require.Equal(t, expected, got)
	})
}
