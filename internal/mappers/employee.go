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
