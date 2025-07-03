package validators

import (
	"regexp"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func ValidateWarehouseCreateRequest(req warehouse.WarehouseRequest) *api.ServiceError {
	if req.Address == "" || req.Telephone == "" || req.WarehouseCode == "" || req.MinimumCapacity <= 0 || req.MinimumTemperature == nil {
		err := api.ServiceErrors[api.ErrUnprocessableEntity]
		err.Message = "invalid request body"
		return &err
	}

	if !isValidPhone(req.Telephone) {
		err := api.ServiceErrors[api.ErrUnprocessableEntity]
		err.Message = "invalid phone number"
		return &err
	}
	return nil
}

func isValidPhone(phone string) bool {
    re := regexp.MustCompile(`^\+?\d{8,15}$`)
    return re.MatchString(phone)
}

