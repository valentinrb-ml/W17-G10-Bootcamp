package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

func ValidatePurchaseOrderPost(po models.RequestPurchaseOrder) error {
	if po.OrderNumber == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "order_number is required")
	}

	if po.OrderDate == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "order_date is required")
	}

	if po.TrackingCode == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "tracking_code is required")
	}

	if po.BuyerID <= 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "buyer_id must be greater than 0")
	}

	if po.ProductRecordID <= 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "product_record_id must be greater than 0")
	}

	return nil
}
