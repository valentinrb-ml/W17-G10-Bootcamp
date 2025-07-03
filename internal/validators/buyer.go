package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

func ValidateRequestBuyer(br models.RequestBuyer) *api.ServiceError {
	if br.CardNumberId == nil || *br.CardNumberId == "" ||
		br.FirstName == nil || *br.FirstName == "" ||
		br.LastName == nil || *br.LastName == "" {
		err := api.ServiceErrors[api.ErrUnprocessableEntity]
		err.Message = "All fields are required. They cannot be empty."
		return &err
	}
	return nil
}

func ValidateUpdateBuyer(br models.RequestBuyer) *api.ServiceError {
	if br.CardNumberId != nil && *br.CardNumberId == "" {
		return createValidationError("card_number_id cannot be empty")
	}
	if br.FirstName != nil && *br.FirstName == "" {
		return createValidationError("first_name cannot be empty")
	}
	if br.LastName != nil && *br.LastName == "" {
		return createValidationError("last_name cannot be empty")
	}
	return nil
}

func createValidationError(message string) *api.ServiceError {
	err := api.ServiceErrors[api.ErrUnprocessableEntity]
	err.Message = message
	return &err
}
