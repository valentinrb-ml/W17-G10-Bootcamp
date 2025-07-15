package validators

import (
	"strings"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

func ValidateInboundOrder(o *models.InboundOrder) error {
	if o == nil {
		return apperrors.NewAppError(apperrors.CodeValidationError, "inbound order cannot be nil")
	}

	if strings.TrimSpace(o.OrderDate) == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "order_date is required")
	}
	if strings.TrimSpace(o.OrderNumber) == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "order_number is required")
	}
	if o.EmployeeID <= 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "employee_id is required and must be positive")
	}
	if o.ProductBatchID <= 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "product_batch_id is required and must be positive")
	}
	if o.WarehouseID <= 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "warehouse_id is required and must be positive")
	}
	return nil
}
