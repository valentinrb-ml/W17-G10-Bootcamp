package testhelpers

import (
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

// Devuelve un empleado genérico para testeo
func CreateTestEmployee() models.Employee {
	return models.Employee{
		CardNumberID: "EMP001",
		FirstName:    "John",
		LastName:     "Doe",
		WarehouseID:  1,
	}
}

// Devuelve un empleado esperado (con ID seteado)
func CreateExpectedEmployee(id int) *models.Employee {
	e := CreateTestEmployee()
	e.ID = id
	return &e
}

// Devuelve varios empleados de prueba
func CreateTestEmployees() []models.Employee {
	return []models.Employee{
		{ID: 1, CardNumberID: "EMP001", FirstName: "John", LastName: "Doe", WarehouseID: 1},
		{ID: 2, CardNumberID: "EMP002", FirstName: "Jane", LastName: "Smith", WarehouseID: 2},
	}
}
