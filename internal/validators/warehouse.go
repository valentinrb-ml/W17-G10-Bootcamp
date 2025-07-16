package validators

import (
	"regexp"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func ValidateWarehouseCreateRequest(req warehouse.WarehouseRequest) error {
	if req.Address == "" || req.Telephone == "" || req.WarehouseCode == "" || req.MinimumCapacity <= 0 || req.MinimumTemperature == nil || req.LocalityId == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "invalid request body")
	}

	if !isValidPhone(req.Telephone) {
		return apperrors.NewAppError(apperrors.CodeValidationError, "invalid phone number")
	}
	return nil
}

func isValidPhone(phone string) bool {
    re := regexp.MustCompile(`^\+?\d{8,15}$`)
    return re.MatchString(phone)
}

func ValidateMinimumCapacity(minimumCapacity int) error {
	if minimumCapacity <= 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "minimum capacity must be greater than 0")
	}
	return nil
}

