package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

func ValidateSectionRequest(secReq section.PostSection) error {
	if secReq.SectionNumber == 0 || secReq.WarehouseId == 0 ||
		secReq.MaximumCapacity == 0 || secReq.MinimumCapacity == 0 ||
		secReq.CurrentCapacity == 0 || secReq.MinimumTemperature == nil ||
		secReq.CurrentTemperature == nil || secReq.ProductTypeId == 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "All fields are required. They cannot be empty.")
	}
	if secReq.MaximumCapacity < 0 || secReq.MinimumCapacity < 0 || secReq.CurrentCapacity < 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "Capacity values cannot be negative.")
	}
	if secReq.MaximumCapacity < secReq.MinimumCapacity {
		return apperrors.NewAppError(apperrors.CodeValidationError, "Maximum capacity cannot be less than minimum capacity.")
	}

	if secReq.CurrentCapacity > secReq.MaximumCapacity {
		return apperrors.NewAppError(apperrors.CodeValidationError, "Current capacity cannot exceed maximum capacity.")
	}
	return nil
}

func ValidateSectionPatch(secReq section.PatchSection) error {
	if secReq.SectionNumber == nil && secReq.WarehouseId == nil &&
		secReq.MaximumCapacity == nil && secReq.MinimumCapacity == nil &&
		secReq.CurrentCapacity == nil && secReq.MinimumTemperature == nil &&
		secReq.CurrentTemperature == nil && secReq.ProductTypeId == nil {
		return apperrors.NewAppError(apperrors.CodeValidationError, "At least one field must be provided to update the section.")
	}

	if secReq.MaximumCapacity != nil && *secReq.MaximumCapacity < 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "Maximum capacity cannot be negative.")
	}
	if secReq.MinimumCapacity != nil && *secReq.MinimumCapacity < 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "Minimum capacity cannot be negative.")
	}
	if secReq.CurrentCapacity != nil && *secReq.CurrentCapacity < 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "Current capacity cannot be negative.")
	}

	if secReq.MaximumCapacity != nil && secReq.MinimumCapacity != nil &&
		*secReq.MaximumCapacity < *secReq.MinimumCapacity {
		return apperrors.NewAppError(apperrors.CodeValidationError, "Maximum capacity cannot be less than minimum capacity.")
	}

	if secReq.CurrentCapacity != nil && secReq.MaximumCapacity != nil &&
		*secReq.CurrentCapacity > *secReq.MaximumCapacity {
		return apperrors.NewAppError(apperrors.CodeValidationError, "Current capacity cannot exceed maximum capacity.")
	}

	return nil

}
