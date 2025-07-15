package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

func ValidateGeographyPost(rg models.RequestGeography) error {
	err := apperrors.NewAppError(apperrors.CodeValidationError, "")

	if rg.CountryName == nil || *rg.CountryName == "" {
		err.Message = "Country Name is required and cannot be empty."
		return err
	}
	if rg.ProvinceName == nil || *rg.ProvinceName == "" {
		err.Message = "Pronvice Name is required and cannot be empty."
		return err
	}
	if rg.LocalityName == nil || *rg.LocalityName == "" {
		err.Message = "Locality Name is required and cannot be empty."
		return err
	}

	return nil
}
