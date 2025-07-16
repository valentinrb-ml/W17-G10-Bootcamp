package validators

import "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"

func ValidateID(value int, fieldName string) error {
	if value <= 0 {
		return apperrors.NewAppError(apperrors.CodeBadRequest, fieldName+" must be a positive integer")
	}
	return nil
}
