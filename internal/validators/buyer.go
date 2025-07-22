package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

func ValidateRequestBuyer(br models.RequestBuyer) *apperrors.AppError {
	if br.CardNumberId == nil || *br.CardNumberId == "" ||
		br.FirstName == nil || *br.FirstName == "" ||
		br.LastName == nil || *br.LastName == "" {
		return apperrors.NewAppError(apperrors.CodeValidationError, "All fields are required. They cannot be empty")
	}
	return nil
}

func ValidateUpdateBuyer(br models.RequestBuyer) *apperrors.AppError {
	if br.CardNumberId != nil && *br.CardNumberId == "" {
		return newValidationError("id_card_number cannot be empty")
	}
	if br.FirstName != nil && *br.FirstName == "" {
		return newValidationError("first_name cannot be empty")
	}
	if br.LastName != nil && *br.LastName == "" {
		return newValidationError("last_name cannot be empty")
	}
	return nil
}

func ValidateBuyerPatchNotEmpty(br models.RequestBuyer) *apperrors.AppError {
	if br.CardNumberId == nil && br.FirstName == nil && br.LastName == nil {
		return apperrors.NewAppError(apperrors.CodeValidationError, "At least one field is required to be updated")
	}
	return nil
}

func newValidationError(message string) *apperrors.AppError {
	return apperrors.NewAppError(apperrors.CodeValidationError, message)
}
