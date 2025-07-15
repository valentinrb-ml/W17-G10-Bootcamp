package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

func ValidateCarryCreateRequest(req carry.CarryRequest) error {
	if req.Address == "" || req.Telephone == "" || req.Cid == "" || req.CompanyName == "" || req.LocalityId <= 0 {
		return apperrors.NewAppError(apperrors.CodeValidationError, "invalid request body")
	}

	if !isValidPhone(req.Telephone) {
		return apperrors.NewAppError(apperrors.CodeValidationError, "invalid phone number")
	}
	return nil
}



