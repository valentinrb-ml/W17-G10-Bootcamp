package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

func ValidateSectionRequest(secReq section.PostSection) *api.ServiceError {
	if secReq.SectionNumber == 0 || secReq.WarehouseId == 0 ||
		secReq.MaximumCapacity == 0 || secReq.MinimumCapacity == 0 ||
		secReq.CurrentCapacity == 0 || secReq.MinimumTemperature == nil ||
		secReq.CurrentTemperature == nil || secReq.ProductTypeId == 0 {
		orig := api.ServiceErrors[api.ErrUnprocessableEntity]
		err := api.ServiceError{
			Code:         orig.Code,
			ResponseCode: orig.ResponseCode,
			Message:      "All fields are required. They cannot be empty.",
		}
		return &err
	}
	if secReq.MaximumCapacity < 0 || secReq.MinimumCapacity < 0 || secReq.CurrentCapacity < 0 {
		orig := api.ServiceErrors[api.ErrUnprocessableEntity]
		return &api.ServiceError{
			Code:         orig.Code,
			ResponseCode: orig.ResponseCode,
			Message:      "Capacity values cannot be negative.",
		}
	}
	if secReq.MaximumCapacity < secReq.MinimumCapacity {
		errDef := api.ServiceErrors[api.ErrUnprocessableEntity]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "Maximum capacity cannot be less than minimum capacity.",
		}
	}
	
	if secReq.CurrentCapacity > secReq.MaximumCapacity {
		errDef := api.ServiceErrors[api.ErrUnprocessableEntity]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "Current capacity cannot exceed maximum capacity.",
		}
	}
	return nil
}

func ValidateSectionPatch(secReq section.PatchSection) *api.ServiceError {
	if secReq.SectionNumber == nil && secReq.WarehouseId == nil &&
		secReq.MaximumCapacity == nil && secReq.MinimumCapacity == nil &&
		secReq.CurrentCapacity == nil && secReq.MinimumTemperature == nil &&
		secReq.CurrentTemperature == nil && secReq.ProductTypeId == nil {
		orig := api.ServiceErrors[api.ErrUnprocessableEntity]
		err := api.ServiceError{
			Code:         orig.Code,
			ResponseCode: orig.ResponseCode,
			Message:      "At least one field must be provided to update the section.",
		}
		return &err
	}

	if secReq.MaximumCapacity != nil && *secReq.MaximumCapacity < 0 {
		errDef := api.ServiceErrors[api.ErrUnprocessableEntity]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "Maximum capacity cannot be negative.",
		}
	}
	if secReq.MinimumCapacity != nil && *secReq.MinimumCapacity < 0 {
		errDef := api.ServiceErrors[api.ErrUnprocessableEntity]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "Minimum capacity cannot be negative.",
		}
	}
	if secReq.CurrentCapacity != nil && *secReq.CurrentCapacity < 0 {
		errDef := api.ServiceErrors[api.ErrUnprocessableEntity]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "Current capacity cannot be negative.",
		}
	}

	if secReq.MaximumCapacity != nil && secReq.MinimumCapacity != nil &&
		*secReq.MaximumCapacity < *secReq.MinimumCapacity {
		errDef := api.ServiceErrors[api.ErrUnprocessableEntity]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "Maximum capacity cannot be less than minimum capacity.",
		}
	}

	if secReq.CurrentCapacity != nil && secReq.MaximumCapacity != nil &&
		*secReq.CurrentCapacity > *secReq.MaximumCapacity {
		errDef := api.ServiceErrors[api.ErrUnprocessableEntity]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "Current capacity cannot exceed maximum capacity.",
		}
	}

	return nil

}
