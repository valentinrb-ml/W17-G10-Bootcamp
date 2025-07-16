package validators

import (
	"time"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/productrecord"
)

func ValidateProductRecordCreateRequest(req productrecord.ProductRecordRequest) error {
	if req.Data.LastUpdateDate.IsZero() {
		return apperrors.NewAppError(apperrors.CodeValidationError, "last_update_date is required")
	}

	if req.Data.ProductID <= 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "product_id must be greater than 0")
	}

	return nil
}

func ValidateProductRecordBusinessRules(record productrecord.ProductRecord) error {
	if record.PurchasePrice < 0 {
		return apperrors.NewAppError(apperrors.CodeBadRequest, "purchase_price must be greater than or equal to 0")
	}

	if record.SalePrice < 0 {
		return apperrors.NewAppError(apperrors.CodeBadRequest, "sale_price must be greater than or equal to 0")
	}

	if record.LastUpdateDate.After(time.Now()) {
		return apperrors.NewAppError(apperrors.CodeBadRequest, "last_update_date cannot be in the future")
	}

	return nil
}
