package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func ValidateSellerPost(sr models.RequestSeller) error {
	err := api.ServiceErrors[api.ErrUnprocessableEntity]

	if sr.Cid == nil || *sr.Cid <= 0 {
		err.Message = "Cid is required and must be greater than 0."
		return &err
	}
	if sr.CompanyName == nil || *sr.CompanyName == "" {
		err.Message = "CompanyName is required and cannot be empty."
		return &err
	}
	if sr.Address == nil || *sr.Address == "" {
		err.Message = "Address is required and cannot be empty."
		return &err
	}
	if sr.Telephone == nil || *sr.Telephone == "" {
		err.Message = "Telephone is required and cannot be empty."
		return &err
	}
	if sr.LocalityId == nil || *sr.LocalityId <= 0 {
		err.Message = "Locality is required and must be greater than 0."
		return &err
	}

	return nil
}

func ValidateSellerPatch(sr models.RequestSeller) error {
	err := api.ServiceErrors[api.ErrUnprocessableEntity]

	if sr.Cid != nil && *sr.Cid <= 0 {
		err.Message = "Cid must be greater than 0."
		return &err
	}
	if sr.CompanyName != nil && *sr.CompanyName == "" {
		err.Message = "CompanyName cannot be empty."
		return &err
	}
	if sr.Address != nil && *sr.Address == "" {
		err.Message = "Address cannot be empty."
		return &err
	}
	if sr.Telephone != nil && *sr.Telephone == "" {
		err.Message = "Telephone cannot be empty."
		return &err
	}
	if sr.LocalityId != nil && *sr.LocalityId <= 0 {
		err.Message = "Locality cannot be empty."
		return &err
	}

	return nil
}

func ValidateSellerPatchNotEmpty(sr models.RequestSeller) error {
	errDef := api.ServiceErrors[api.ErrUnprocessableEntity]

	if sr.Cid == nil &&
		sr.CompanyName == nil &&
		sr.Address == nil &&
		sr.Telephone == nil &&
		sr.LocalityId == nil {
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "at least one field is required to be updated.",
		}
	}

	return nil
}
