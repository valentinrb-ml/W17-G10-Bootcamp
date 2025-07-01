package mappers

import models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"

func MapEmployeeDocToEmployee(doc models.EmployeeDoc) *models.Employee {
	return &models.Employee{
		ID:           doc.ID,
		CardNumberID: doc.CardNumberID,
		FirstName:    doc.FirstName,
		LastName:     doc.LastName,
		WarehouseID:  doc.WarehouseID,
	}
}

func MapEmployeeToEmployeeDoc(e *models.Employee) models.EmployeeDoc {
	return models.EmployeeDoc{
		ID:           e.ID,
		CardNumberID: e.CardNumberID,
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		WarehouseID:  e.WarehouseID,
	}
}
func MapEmployeePatchToEmployee(existing *models.Employee, patch *models.EmployeePatch) *models.Employee {
	emp := *existing
	if patch.CardNumberID != nil {
		emp.CardNumberID = *patch.CardNumberID
	}
	if patch.FirstName != nil {
		emp.FirstName = *patch.FirstName
	}
	if patch.LastName != nil {
		emp.LastName = *patch.LastName
	}
	if patch.WarehouseID != nil && *patch.WarehouseID != 0 {
		emp.WarehouseID = *patch.WarehouseID
	}
	return &emp
}
