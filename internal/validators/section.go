package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

func ValidateSectionRequest(secReq section.RequestSection) *api.ServiceError {
	if secReq.SectionNumber == nil || secReq.WarehouseId == nil ||
		secReq.MaximumCapacity == nil || secReq.MinimumCapacity == nil ||
		secReq.CurrentCapacity == nil || secReq.MinimumTemperature == nil ||
		secReq.CurrentTemperature == nil || secReq.ProductId == 0 {
		orig := api.ServiceErrors[api.ErrUnprocessableEntity]
		err := api.ServiceError{
			Code:         orig.Code,
			ResponseCode: orig.ResponseCode,
			Message:      "All fields are required. They cannot be empty.",
		}
		return &err
	}
	return nil
}

func ValidateSectionPatch(secReq section.RequestSection) *api.ServiceError {
	if secReq.SectionNumber == nil && secReq.WarehouseId == nil &&
		secReq.MaximumCapacity == nil && secReq.MinimumCapacity == nil &&
		secReq.CurrentCapacity == nil && secReq.MinimumTemperature == nil &&
		secReq.CurrentTemperature == nil && secReq.ProductId == 0 {
		orig := api.ServiceErrors[api.ErrUnprocessableEntity]
		err := api.ServiceError{
			Code:         orig.Code,
			ResponseCode: orig.ResponseCode,
			Message:      "At least one field must be provided to update the section.",
		}
		return &err
	}
	return nil

}
