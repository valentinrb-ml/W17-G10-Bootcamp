package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func ValidateRequestSeller(sr models.RequestSeller) *api.ServiceError {
	if sr.Cid == nil || *sr.Cid <= 0 ||
		sr.CompanyName == nil || *sr.CompanyName == "" ||
		sr.Address == nil || *sr.Address == "" ||
		sr.Telephone == nil || *sr.Telephone == "" {
		err := api.ServiceErrors[api.ErrUnprocessableEntity]
		err.Message = "All fields are required. They cannot be empty."
		return &err
	}

	return nil
}
