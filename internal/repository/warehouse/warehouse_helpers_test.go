package repository_test

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

// createTestWarehouse creates a warehouse for testing purposes
func createTestWarehouse() warehouse.Warehouse {
	return warehouse.Warehouse{
		WarehouseCode:      "WH001",
		Address:            "123 Main St",
		MinimumTemperature: 10.5,
		MinimumCapacity:    1000,
		Telephone:          "555-1234",
		LocalityId:         "LOC001",
	}
}

// createExpectedWarehouse creates a warehouse with ID for expected results
func createExpectedWarehouse(id int) *warehouse.Warehouse {
	w := createTestWarehouse()
	w.Id = id
	return &w
}

func createExpectedWarehouses() []warehouse.Warehouse {
	return []warehouse.Warehouse{
		{
			Id:                 1,
			WarehouseCode:      "WH001",
			Address:            "123 Main St",
			MinimumTemperature: 10.5,
			MinimumCapacity:    1000,
			Telephone:          "555-1234",
			LocalityId:         "LOC001",
		},
		{
			Id:                 2,
			WarehouseCode:      "WH002",
			Address:            "456 Elm St",
			MinimumTemperature: 15.5,
			MinimumCapacity:    2000,
			Telephone:          "555-5678",
			LocalityId:         "LOC002",
		},
	}
}

// createMockDB creates a mock database for testing
func createMockDB() (sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return mock, db
}
