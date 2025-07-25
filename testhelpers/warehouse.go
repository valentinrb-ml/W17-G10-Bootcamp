package testhelpers

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

// CreateTestWarehouse creates a warehouse for testing purposes
func CreateTestWarehouse() warehouse.Warehouse {
	return warehouse.Warehouse{
		WarehouseCode:      "WH001",
		Address:            "123 Main St",
		MinimumTemperature: 10.5,
		MinimumCapacity:    1000,
		Telephone:          "5551234567", // Formato válido sin guiones
		LocalityId:         "LOC001",
	}
}

// CreateTestWarehouseRequest creates a WarehouseRequest for testing handler endpoints
func CreateTestWarehouseRequest() warehouse.WarehouseRequest {
	return warehouse.WarehouseRequest{
		WarehouseCode:      "WH001",
		Address:            "123 Main St",
		MinimumTemperature: Float64Ptr(10.5),
		MinimumCapacity:    1000,
		Telephone:          "5551234567", // Formato válido sin guiones
		LocalityId:         "LOC001",
	}
}

// CreateExpectedWarehouse creates a warehouse with ID for expected results
func CreateExpectedWarehouse(id int) *warehouse.Warehouse {
	w := CreateTestWarehouse()
	w.Id = id
	return &w
}

// CreateTestWarehouses creates multiple warehouses for testing
func CreateTestWarehouses() []warehouse.Warehouse {
	return []warehouse.Warehouse{
		{
			Id:                 1,
			WarehouseCode:      "WH001",
			Address:            "123 Main St",
			MinimumTemperature: 10.5,
			MinimumCapacity:    1000,
			Telephone:          "5551234567", // Formato válido sin guiones
			LocalityId:         "LOC001",
		},
		{
			Id:                 2,
			WarehouseCode:      "WH002",
			Address:            "456 Elm St",
			MinimumTemperature: 15.5,
			MinimumCapacity:    2000,
			Telephone:          "5555678901", // Formato válido sin guiones
			LocalityId:         "LOC002",
		},
	}
}

// CreateMockDB creates a mock database for testing (for repository tests)
func CreateMockDB() (sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return mock, db
}

// Helper functions for pointer values
func IntPtr(i int) *int {
	return &i
}

func StringPtr(s string) *string {
	return &s
}

func Float64Ptr(f float64) *float64 {
	return &f
}
