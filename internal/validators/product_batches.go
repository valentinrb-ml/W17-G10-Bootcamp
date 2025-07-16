package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

func ValidateProductBatchPost(p models.PostProductBatches) error {
	if p.SectionId == 0 || p.ProductId == 0 || p.MinimumTemperature == nil ||
		p.ManufacturingHour == 0 || p.ManufacturingDate.IsZero() || p.InitialQuantity == nil ||
		p.DueDate.IsZero() || p.CurrentTemperature == nil || p.CurrentQuantity == nil || p.BatchNumber == 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "All fields are required. They cannot be empty.")
	}
	if *p.CurrentQuantity < 0 || *p.InitialQuantity < 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "Quantity values cannot be negative.")
	}
	return nil
}
