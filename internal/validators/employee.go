package validators

import (
	"strings"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

func ValidateEmployee(e *models.Employee) error {
	if e == nil {
		return apperrors.NewAppError(apperrors.CodeValidationError, "employee cannot be nil")
	}
	if strings.TrimSpace(e.CardNumberID) == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "card_number_id cannot be empty")
	}
	if strings.TrimSpace(e.FirstName) == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "first_name cannot be empty")
	}
	if strings.TrimSpace(e.LastName) == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "last_name cannot be empty")
	}
	if e.WarehouseID == 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "warehouse_id is required")
	}
	if e.WarehouseID < 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "warehouse_id must be greater than 0")
	}
	return nil
}

func ValidateEmployeeID(id int) error {
	if id <= 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "id must be positive")
	}
	return nil
}

func ValidateEmployeePatch(e *models.EmployeePatch) error {
	if e.CardNumberID != nil && strings.TrimSpace(*e.CardNumberID) == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "card_number_id cannot be empty")
	}
	if e.FirstName != nil && strings.TrimSpace(*e.FirstName) == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "first_name cannot be empty")
	}
	if e.LastName != nil && strings.TrimSpace(*e.LastName) == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "last_name cannot be empty")
	}
	if e.WarehouseID != nil && *e.WarehouseID <= 0 {
		if *e.WarehouseID == 0 {
			return apperrors.NewAppError(apperrors.CodeValidationError, "warehouse_id is required")
		}
		return apperrors.NewAppError(apperrors.CodeValidationError, "warehouse_id must be positive")
	}
	return nil
}
