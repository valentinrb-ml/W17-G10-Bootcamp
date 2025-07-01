package validators

import (
	"strings"

	api "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

// Retorna un ServiceError 422 con mensaje específico
func ValidateEmployee(e *models.Employee) error {
	if e == nil {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.InternalError = nil
		se.Message = "employee cannot be nil"
		return se
	}
	if strings.TrimSpace(e.CardNumberID) == "" {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.InternalError = nil
		se.Message = "card_number_id cannot be empty"
		return se
	}
	if strings.TrimSpace(e.FirstName) == "" {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.InternalError = nil
		se.Message = "first_name cannot be empty"
		return se
	}
	if strings.TrimSpace(e.LastName) == "" {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.InternalError = nil
		se.Message = "last_name cannot be empty"
		return se
	}
	if e.WarehouseID <= 0 {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.InternalError = nil
		se.Message = "warehouse_id must be greater than 0"
		return se
	}
	return nil
}
func ValidateEmployeeID(id int) error {
	if id <= 0 {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.Message = "id must be positive"
		return se
	}
	return nil
}
func ValidateEmployeePatch(e *models.EmployeePatch) error {
	if e.CardNumberID != nil && strings.TrimSpace(*e.CardNumberID) == "" {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.Message = "card_number_id cannot be empty"
		return se
	}
	if e.FirstName != nil && strings.TrimSpace(*e.FirstName) == "" {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.Message = "first_name cannot be empty"
		return se
	}
	if e.LastName != nil && strings.TrimSpace(*e.LastName) == "" {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.Message = "last_name cannot be empty"
		return se
	}
	if e.WarehouseID != nil && *e.WarehouseID <= 0 {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.Message = "warehouse_id must be positive"
		return se
	}
	return nil
}
