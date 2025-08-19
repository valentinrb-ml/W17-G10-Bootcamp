package mappers_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func TestWarehouseToDoc(t *testing.T) {
    t.Run("convert warehouse to doc successfully", func(t *testing.T) {
        // Arrange
        inputWarehouse := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "WH-001",
            Address:            "123 Main St",
            Telephone:          "+1234567890",
            MinimumCapacity:    100,
            MinimumTemperature: -18.5,
            LocalityId:         "LOC-42",
        }

        expectedDoc := warehouse.WarehouseDoc{
            ID:                 1,
            WarehouseCode:      "WH-001",
            Address:            "123 Main St",
            Telephone:          "+1234567890",
            MinimumCapacity:    100,
            MinimumTemperature: -18.5,
            LocalityId:         "LOC-42",
        }

        // Act
        result := mappers.WarehouseToDoc(inputWarehouse)

        // Assert
        assert.Equal(t, expectedDoc.ID, result.ID)
        assert.Equal(t, expectedDoc.WarehouseCode, result.WarehouseCode)
        assert.Equal(t, expectedDoc.Address, result.Address)
        assert.Equal(t, expectedDoc.Telephone, result.Telephone)
        assert.Equal(t, expectedDoc.MinimumCapacity, result.MinimumCapacity)
        assert.Equal(t, expectedDoc.MinimumTemperature, result.MinimumTemperature)
        assert.Equal(t, expectedDoc.LocalityId, result.LocalityId)
        assert.Equal(t, expectedDoc, result)
    })

    t.Run("convert warehouse with zero values", func(t *testing.T) {
        // Arrange
        inputWarehouse := &warehouse.Warehouse{
            Id:                 0,
            WarehouseCode:      "",
            Address:            "",
            Telephone:          "",
            MinimumCapacity:    0,
            MinimumTemperature: 0.0,
            LocalityId:         "",
        }

        expectedDoc := warehouse.WarehouseDoc{
            ID:                 0,
            WarehouseCode:      "",
            Address:            "",
            Telephone:          "",
            MinimumCapacity:    0,
            MinimumTemperature: 0.0,
            LocalityId:         "",
        }

        // Act
        result := mappers.WarehouseToDoc(inputWarehouse)

        // Assert
        assert.Equal(t, expectedDoc, result)
    })

    t.Run("convert warehouse with negative values", func(t *testing.T) {
        // Arrange
        inputWarehouse := &warehouse.Warehouse{
            Id:                 -1,
            WarehouseCode:      "NEG-WH",
            Address:            "Negative Street",
            Telephone:          "-123456789",
            MinimumCapacity:    -50,
            MinimumTemperature: -40.5,
            LocalityId:         "NEG-LOC",
        }

        expectedDoc := warehouse.WarehouseDoc{
            ID:                 -1,
            WarehouseCode:      "NEG-WH",
            Address:            "Negative Street",
            Telephone:          "-123456789",
            MinimumCapacity:    -50,
            MinimumTemperature: -40.5,
            LocalityId:         "NEG-LOC",
        }

        // Act
        result := mappers.WarehouseToDoc(inputWarehouse)

        // Assert
        assert.Equal(t, expectedDoc, result)
    })

    t.Run("convert warehouse with maximum values", func(t *testing.T) {
        // Arrange
        inputWarehouse := &warehouse.Warehouse{
            Id:                 2147483647, // max int32
            WarehouseCode:      "MAX-WAREHOUSE-CODE-WITH-VERY-LONG-NAME",
            Address:            "Very Long Address That Could Be Stored In Database Field",
            Telephone:          "+999999999999999999",
            MinimumCapacity:    2147483647, // max int32
            MinimumTemperature: 999999.999,
            LocalityId:         "MAX-LOCALITY-ID-STRING",
        }

        expectedDoc := warehouse.WarehouseDoc{
            ID:                 2147483647,
            WarehouseCode:      "MAX-WAREHOUSE-CODE-WITH-VERY-LONG-NAME",
            Address:            "Very Long Address That Could Be Stored In Database Field",
            Telephone:          "+999999999999999999",
            MinimumCapacity:    2147483647,
            MinimumTemperature: 999999.999,
            LocalityId:         "MAX-LOCALITY-ID-STRING",
        }

        // Act
        result := mappers.WarehouseToDoc(inputWarehouse)

        // Assert
        assert.Equal(t, expectedDoc, result)
    })

    t.Run("convert warehouse with special characters", func(t *testing.T) {
        // Arrange
        inputWarehouse := &warehouse.Warehouse{
            Id:                 123,
            WarehouseCode:      "WH@#$%",
            Address:            "123 Spëcîál Çhárãctérs St. & Co.",
            Telephone:          "+1-234-567-8900 ext.123",
            MinimumCapacity:    50,
            MinimumTemperature: -18.75,
            LocalityId:         "LOC@#$%",
        }

        expectedDoc := warehouse.WarehouseDoc{
            ID:                 123,
            WarehouseCode:      "WH@#$%",
            Address:            "123 Spëcîál Çhárãctérs St. & Co.",
            Telephone:          "+1-234-567-8900 ext.123",
            MinimumCapacity:    50,
            MinimumTemperature: -18.75,
            LocalityId:         "LOC@#$%",
        }

        // Act
        result := mappers.WarehouseToDoc(inputWarehouse)

        // Assert
        assert.Equal(t, expectedDoc, result)
    })

    t.Run("verify all fields are mapped correctly", func(t *testing.T) {
        // Arrange
        inputWarehouse := &warehouse.Warehouse{
            Id:                 777,
            WarehouseCode:      "FIELD-TEST",
            Address:            "Field Mapping Test Address",
            Telephone:          "+1111111111",
            MinimumCapacity:    200,
            MinimumTemperature: 25.5,
            LocalityId:         "LOC-888",
        }

        // Act
        result := mappers.WarehouseToDoc(inputWarehouse)

        // Assert - Verify each field individually
        assert.Equal(t, inputWarehouse.Id, result.ID, "ID field should be mapped correctly")
        assert.Equal(t, inputWarehouse.WarehouseCode, result.WarehouseCode, "WarehouseCode field should be mapped correctly")
        assert.Equal(t, inputWarehouse.Address, result.Address, "Address field should be mapped correctly")
        assert.Equal(t, inputWarehouse.Telephone, result.Telephone, "Telephone field should be mapped correctly")
        assert.Equal(t, inputWarehouse.MinimumCapacity, result.MinimumCapacity, "MinimumCapacity field should be mapped correctly")
        assert.Equal(t, inputWarehouse.MinimumTemperature, result.MinimumTemperature, "MinimumTemperature field should be mapped correctly")
        assert.Equal(t, inputWarehouse.LocalityId, result.LocalityId, "LocalityId field should be mapped correctly")
    })

    t.Run("verify struct types are different", func(t *testing.T) {
        // Arrange
        inputWarehouse := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "TYPE-TEST",
            Address:            "Type Test Address",
            Telephone:          "+9999999999",
            MinimumCapacity:    300,
            MinimumTemperature: 10.0,
            LocalityId:         "LOC-999",
        }

        // Act
        result := mappers.WarehouseToDoc(inputWarehouse)

        // Assert - Verify types
        assert.IsType(t, warehouse.WarehouseDoc{}, result, "Result should be of type WarehouseDoc")
        assert.IsType(t, &warehouse.Warehouse{}, inputWarehouse, "Input should be of type *Warehouse")
    })

    t.Run("handle nil pointer input", func(t *testing.T) {
        // This test verifies the function behavior with nil input
        // Note: This will panic, which is expected behavior for this function
        assert.Panics(t, func() {
            mappers.WarehouseToDoc(nil)
        }, "Function should panic when given nil pointer")
    })

    t.Run("test with floating point precision", func(t *testing.T) {
        // Arrange
        inputWarehouse := &warehouse.Warehouse{
            Id:                 456,
            WarehouseCode:      "FLOAT-TEST",
            Address:            "Float Precision Test",
            Telephone:          "+7777777777",
            MinimumCapacity:    150,
            MinimumTemperature: -18.123456789,
            LocalityId:         "LOC-321",
        }

        // Act
        result := mappers.WarehouseToDoc(inputWarehouse)

        // Assert
        assert.Equal(t, -18.123456789, result.MinimumTemperature, "Should preserve floating point precision")
        assert.InDelta(t, -18.123456789, result.MinimumTemperature, 0.000000001, "Should be within delta tolerance")
    })

    t.Run("test immutability - input not modified", func(t *testing.T) {
        // Arrange
        originalWarehouse := &warehouse.Warehouse{
            Id:                 555,
            WarehouseCode:      "IMMUTABLE-TEST",
            Address:            "Immutability Test Address",
            Telephone:          "+5555555555",
            MinimumCapacity:    250,
            MinimumTemperature: 5.5,
            LocalityId:         "LOC-666",
        }

        // Create copies of original values for comparison
        originalId := originalWarehouse.Id
        originalCode := originalWarehouse.WarehouseCode
        originalAddress := originalWarehouse.Address
        originalTelephone := originalWarehouse.Telephone
        originalCapacity := originalWarehouse.MinimumCapacity
        originalTemperature := originalWarehouse.MinimumTemperature
        originalLocalityId := originalWarehouse.LocalityId

        // Act
        _ = mappers.WarehouseToDoc(originalWarehouse)

        // Assert - Original warehouse should not be modified
        assert.Equal(t, originalId, originalWarehouse.Id, "Original Id should not be modified")
        assert.Equal(t, originalCode, originalWarehouse.WarehouseCode, "Original WarehouseCode should not be modified")
        assert.Equal(t, originalAddress, originalWarehouse.Address, "Original Address should not be modified")
        assert.Equal(t, originalTelephone, originalWarehouse.Telephone, "Original Telephone should not be modified")
        assert.Equal(t, originalCapacity, originalWarehouse.MinimumCapacity, "Original MinimumCapacity should not be modified")
        assert.Equal(t, originalTemperature, originalWarehouse.MinimumTemperature, "Original MinimumTemperature should not be modified")
        assert.Equal(t, originalLocalityId, originalWarehouse.LocalityId, "Original LocalityId should not be modified")
    })
}

// Benchmark test for WarehouseToDoc function
func BenchmarkWarehouseToDoc(b *testing.B) {
    warehouse := &warehouse.Warehouse{
        Id:                 1,
        WarehouseCode:      "BENCH-001",
        Address:            "Benchmark Test Address",
        Telephone:          "+1234567890",
        MinimumCapacity:    100,
        MinimumTemperature: -18.5,
        LocalityId:         "LOC-1",
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        mappers.WarehouseToDoc(warehouse)
    }
}

func TestWarehouseToDocSlice(t *testing.T) {
    t.Run("convert slice of warehouses successfully", func(t *testing.T) {
        // Arrange
        inputWarehouses := []warehouse.Warehouse{
            {
                Id:                 1,
                WarehouseCode:      "WH-001",
                Address:            "123 Main St",
                Telephone:          "+1234567890",
                MinimumCapacity:    100,
                MinimumTemperature: -18.5,
                LocalityId:         "LOC-42",
            },
            {
                Id:                 2,
                WarehouseCode:      "WH-002", 
                Address:            "456 Oak Ave",
                Telephone:          "+0987654321",
                MinimumCapacity:    200,
                MinimumTemperature: -20.0,
                LocalityId:         "LOC-43",
            },
        }

        // Act
        result := mappers.WarehouseToDocSlice(inputWarehouses)

        // Assert
        assert.Len(t, result, 2)
        assert.Equal(t, inputWarehouses[0].Id, result[0].ID)
        assert.Equal(t, inputWarehouses[0].WarehouseCode, result[0].WarehouseCode)
        assert.Equal(t, inputWarehouses[1].Id, result[1].ID)
        assert.Equal(t, inputWarehouses[1].WarehouseCode, result[1].WarehouseCode)
        assert.Equal(t, inputWarehouses[0].LocalityId, result[0].LocalityId)
        assert.Equal(t, inputWarehouses[1].LocalityId, result[1].LocalityId)
    })

    t.Run("convert empty slice", func(t *testing.T) {
        // Arrange
        inputWarehouses := []warehouse.Warehouse{}

        // Act
        result := mappers.WarehouseToDocSlice(inputWarehouses)

        // Assert
        assert.Empty(t, result)
        assert.NotNil(t, result)
        assert.Len(t, result, 0)
    })

    t.Run("convert nil slice", func(t *testing.T) {
        // Arrange
        var inputWarehouses []warehouse.Warehouse

        // Act
        result := mappers.WarehouseToDocSlice(inputWarehouses)

        // Assert
        assert.Empty(t, result)
        assert.NotNil(t, result)
        assert.Len(t, result, 0)
    })

    t.Run("convert single warehouse slice", func(t *testing.T) {
        // Arrange
        inputWarehouses := []warehouse.Warehouse{
            {
                Id:                 999,
                WarehouseCode:      "SINGLE-WH",
                Address:            "Single Address",
                Telephone:          "+1111111111",
                MinimumCapacity:    50,
                MinimumTemperature: 0.0,
                LocalityId:         "LOC-1",
            },
        }

        // Act
        result := mappers.WarehouseToDocSlice(inputWarehouses)

        // Assert
        assert.Len(t, result, 1)
        assert.Equal(t, inputWarehouses[0].Id, result[0].ID)
        assert.Equal(t, inputWarehouses[0].WarehouseCode, result[0].WarehouseCode)
        assert.Equal(t, inputWarehouses[0].LocalityId, result[0].LocalityId)
    })
}

func TestRequestToWarehouse(t *testing.T) {
    t.Run("convert request to warehouse successfully", func(t *testing.T) {
        // Arrange
        minTemp := -18.5
        inputRequest := warehouse.WarehouseRequest{
            WarehouseCode:      "REQ-001",
            Address:            "Request Address",
            Telephone:          "+1234567890",
            MinimumCapacity:    100,
            MinimumTemperature: &minTemp,
            LocalityId:         "LOC-42",
        }

        expectedWarehouse := warehouse.Warehouse{
            Id:                 0, // Always set to 0
            WarehouseCode:      "REQ-001",
            Address:            "Request Address",
            Telephone:          "+1234567890",
            MinimumCapacity:    100,
            MinimumTemperature: -18.5,
            LocalityId:         "LOC-42",
        }

        // Act
        result := mappers.RequestToWarehouse(inputRequest)

        // Assert
        assert.Equal(t, expectedWarehouse.Id, result.Id)
        assert.Equal(t, expectedWarehouse.WarehouseCode, result.WarehouseCode)
        assert.Equal(t, expectedWarehouse.Address, result.Address)
        assert.Equal(t, expectedWarehouse.Telephone, result.Telephone)
        assert.Equal(t, expectedWarehouse.MinimumCapacity, result.MinimumCapacity)
        assert.Equal(t, expectedWarehouse.MinimumTemperature, result.MinimumTemperature)
        assert.Equal(t, expectedWarehouse.LocalityId, result.LocalityId)
    })

    t.Run("convert request with zero values", func(t *testing.T) {
        // Arrange
        minTemp := 0.0
        inputRequest := warehouse.WarehouseRequest{
            WarehouseCode:      "",
            Address:            "",
            Telephone:          "",
            MinimumCapacity:    0,
            MinimumTemperature: &minTemp,
            LocalityId:         "",
        }

        // Act
        result := mappers.RequestToWarehouse(inputRequest)

        // Assert
        assert.Equal(t, 0, result.Id)
        assert.Equal(t, "", result.WarehouseCode)
        assert.Equal(t, "", result.Address)
        assert.Equal(t, "", result.Telephone)
        assert.Equal(t, 0, result.MinimumCapacity)
        assert.Equal(t, 0.0, result.MinimumTemperature)
        assert.Equal(t, "", result.LocalityId)
    })

    t.Run("convert request with negative temperature", func(t *testing.T) {
        // Arrange
        minTemp := -40.5
        inputRequest := warehouse.WarehouseRequest{
            WarehouseCode:      "COLD-WH",
            Address:            "Antarctica",
            Telephone:          "+000000000",
            MinimumCapacity:    1000,
            MinimumTemperature: &minTemp,
            LocalityId:         "LOC-999",
        }

        // Act
        result := mappers.RequestToWarehouse(inputRequest)

        // Assert
        assert.Equal(t, -40.5, result.MinimumTemperature)
        assert.Equal(t, "COLD-WH", result.WarehouseCode)
        assert.Equal(t, "LOC-999", result.LocalityId)
    })
}

func TestApplyWarehousePatch(t *testing.T) {
    t.Run("apply patch with all fields", func(t *testing.T) {
        // Arrange
        existing := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "OLD-CODE",
            Address:            "Old Address",
            Telephone:          "+1111111111",
            MinimumCapacity:    50,
            MinimumTemperature: -10.0,
            LocalityId:         "LOC-1",
        }

        newAddress := "New Address"
        newTelephone := "+2222222222"
        newCode := "NEW-CODE"
        newCapacity := 100
        newTemp := -20.5
        newLocalityId := "LOC-2"

        patch := warehouse.WarehousePatchDTO{
            Address:            &newAddress,
            Telephone:          &newTelephone,
            WarehouseCode:      &newCode,
            MinimumCapacity:    &newCapacity,
            MinimumTemperature: &newTemp,
            LocalityId:         &newLocalityId,
        }

        // Act
        mappers.ApplyWarehousePatch(existing, patch)

        // Assert
        assert.Equal(t, "NEW-CODE", existing.WarehouseCode)
        assert.Equal(t, "New Address", existing.Address)
        assert.Equal(t, "+2222222222", existing.Telephone)
        assert.Equal(t, 100, existing.MinimumCapacity)
        assert.Equal(t, -20.5, existing.MinimumTemperature)
        assert.Equal(t, "LOC-2", existing.LocalityId)
        assert.Equal(t, 1, existing.Id) // ID should remain unchanged
    })

    t.Run("apply patch with only address", func(t *testing.T) {
        // Arrange
        existing := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "UNCHANGED",
            Address:            "Old Address",
            Telephone:          "+1111111111",
            MinimumCapacity:    50,
            MinimumTemperature: -10.0,
            LocalityId:         "LOC-1",
        }

        newAddress := "Only Address Changed"
        patch := warehouse.WarehousePatchDTO{
            Address: &newAddress,
        }

        // Act
        mappers.ApplyWarehousePatch(existing, patch)

        // Assert
        assert.Equal(t, "Only Address Changed", existing.Address)
        assert.Equal(t, "UNCHANGED", existing.WarehouseCode) // Should remain unchanged
        assert.Equal(t, "+1111111111", existing.Telephone)   // Should remain unchanged
        assert.Equal(t, 50, existing.MinimumCapacity)         // Should remain unchanged
        assert.Equal(t, -10.0, existing.MinimumTemperature)   // Should remain unchanged
        assert.Equal(t, "LOC-1", existing.LocalityId)         // Should remain unchanged
    })

    t.Run("apply patch with only telephone", func(t *testing.T) {
        // Arrange
        existing := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "UNCHANGED",
            Address:            "Unchanged Address",
            Telephone:          "+1111111111",
            MinimumCapacity:    50,
            MinimumTemperature: -10.0,
            LocalityId:         "LOC-1",
        }

        newTelephone := "+9999999999"
        patch := warehouse.WarehousePatchDTO{
            Telephone: &newTelephone,
        }

        // Act
        mappers.ApplyWarehousePatch(existing, patch)

        // Assert
        assert.Equal(t, "+9999999999", existing.Telephone)
        assert.Equal(t, "UNCHANGED", existing.WarehouseCode)
        assert.Equal(t, "Unchanged Address", existing.Address)
        assert.Equal(t, "LOC-1", existing.LocalityId)
    })

    t.Run("apply patch with only warehouse code", func(t *testing.T) {
        // Arrange
        existing := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "OLD-CODE",
            Address:            "Unchanged Address",
            Telephone:          "+1111111111",
            MinimumCapacity:    50,
            MinimumTemperature: -10.0,
            LocalityId:         "LOC-1",
        }

        newCode := "UPDATED-CODE"
        patch := warehouse.WarehousePatchDTO{
            WarehouseCode: &newCode,
        }

        // Act
        mappers.ApplyWarehousePatch(existing, patch)

        // Assert
        assert.Equal(t, "UPDATED-CODE", existing.WarehouseCode)
        assert.Equal(t, "Unchanged Address", existing.Address)
        assert.Equal(t, "+1111111111", existing.Telephone)
        assert.Equal(t, "LOC-1", existing.LocalityId)
    })

    t.Run("apply patch with only minimum capacity", func(t *testing.T) {
        // Arrange
        existing := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "UNCHANGED",
            Address:            "Unchanged Address",
            Telephone:          "+1111111111",
            MinimumCapacity:    50,
            MinimumTemperature: -10.0,
            LocalityId:         "LOC-1",
        }

        newCapacity := 200
        patch := warehouse.WarehousePatchDTO{
            MinimumCapacity: &newCapacity,
        }

        // Act
        mappers.ApplyWarehousePatch(existing, patch)

        // Assert
        assert.Equal(t, 200, existing.MinimumCapacity)
        assert.Equal(t, "UNCHANGED", existing.WarehouseCode)
        assert.Equal(t, -10.0, existing.MinimumTemperature)
        assert.Equal(t, "LOC-1", existing.LocalityId)
    })

    t.Run("apply patch with only minimum temperature", func(t *testing.T) {
        // Arrange
        existing := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "UNCHANGED",
            Address:            "Unchanged Address",
            Telephone:          "+1111111111",
            MinimumCapacity:    50,
            MinimumTemperature: -10.0,
            LocalityId:         "LOC-1",
        }

        newTemp := -25.5
        patch := warehouse.WarehousePatchDTO{
            MinimumTemperature: &newTemp,
        }

        // Act
        mappers.ApplyWarehousePatch(existing, patch)

        // Assert
        assert.Equal(t, -25.5, existing.MinimumTemperature)
        assert.Equal(t, 50, existing.MinimumCapacity)
        assert.Equal(t, "UNCHANGED", existing.WarehouseCode)
        assert.Equal(t, "LOC-1", existing.LocalityId)
    })

    t.Run("apply patch with only locality id", func(t *testing.T) {
        // Arrange
        existing := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "UNCHANGED",
            Address:            "Unchanged Address",
            Telephone:          "+1111111111",
            MinimumCapacity:    50,
            MinimumTemperature: -10.0,
            LocalityId:         "LOC-1",
        }

        newLocalityId := "LOC-99"
        patch := warehouse.WarehousePatchDTO{
            LocalityId: &newLocalityId,
        }

        // Act
        mappers.ApplyWarehousePatch(existing, patch)

        // Assert
        assert.Equal(t, "LOC-99", existing.LocalityId)
        assert.Equal(t, "UNCHANGED", existing.WarehouseCode)
        assert.Equal(t, 50, existing.MinimumCapacity)
    })

    t.Run("apply empty patch - no changes", func(t *testing.T) {
        // Arrange
        existing := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "UNCHANGED",
            Address:            "Unchanged Address",
            Telephone:          "+1111111111",
            MinimumCapacity:    50,
            MinimumTemperature: -10.0,
            LocalityId:         "LOC-1",
        }

        originalValues := *existing // Make a copy
        patch := warehouse.WarehousePatchDTO{} // Empty patch

        // Act
        mappers.ApplyWarehousePatch(existing, patch)

        // Assert - Nothing should change
        assert.Equal(t, originalValues.Id, existing.Id)
        assert.Equal(t, originalValues.WarehouseCode, existing.WarehouseCode)
        assert.Equal(t, originalValues.Address, existing.Address)
        assert.Equal(t, originalValues.Telephone, existing.Telephone)
        assert.Equal(t, originalValues.MinimumCapacity, existing.MinimumCapacity)
        assert.Equal(t, originalValues.MinimumTemperature, existing.MinimumTemperature)
        assert.Equal(t, originalValues.LocalityId, existing.LocalityId)
    })

    t.Run("apply patch with zero values", func(t *testing.T) {
        // Arrange
        existing := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "OLD-CODE",
            Address:            "Old Address",
            Telephone:          "+1111111111",
            MinimumCapacity:    50,
            MinimumTemperature: -10.0,
            LocalityId:         "LOC-1",
        }

        emptyAddress := ""
        emptyTelephone := ""
        emptyCode := ""
        zeroCapacity := 0
        zeroTemp := 0.0
        emptyLocalityId := ""

        patch := warehouse.WarehousePatchDTO{
            Address:            &emptyAddress,
            Telephone:          &emptyTelephone,
            WarehouseCode:      &emptyCode,
            MinimumCapacity:    &zeroCapacity,
            MinimumTemperature: &zeroTemp,
            LocalityId:         &emptyLocalityId,
        }

        // Act
        mappers.ApplyWarehousePatch(existing, patch)

        // Assert
        assert.Equal(t, "", existing.WarehouseCode)
        assert.Equal(t, "", existing.Address)
        assert.Equal(t, "", existing.Telephone)
        assert.Equal(t, 0, existing.MinimumCapacity)
        assert.Equal(t, 0.0, existing.MinimumTemperature)
        assert.Equal(t, "", existing.LocalityId)
    })
}

// Benchmarks adicionales
func BenchmarkWarehouseToDocSlice(b *testing.B) {
    warehouses := make([]warehouse.Warehouse, 100)
    for i := 0; i < 100; i++ {
        warehouses[i] = warehouse.Warehouse{
            Id:                 i,
            WarehouseCode:      "BENCH-" + string(rune(i+65)), // A, B, C...
            Address:            "Benchmark Address",
            Telephone:          "+1234567890",
            MinimumCapacity:    100,
            MinimumTemperature: -18.5,
            LocalityId:         "LOC-1",
        }
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        mappers.WarehouseToDocSlice(warehouses)
    }
}

func BenchmarkRequestToWarehouse(b *testing.B) {
    minTemp := -18.5
    request := warehouse.WarehouseRequest{
        WarehouseCode:      "BENCH-001",
        Address:            "Benchmark Address",
        Telephone:          "+1234567890",
        MinimumCapacity:    100,
        MinimumTemperature: &minTemp,
        LocalityId:         "LOC-1",
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        mappers.RequestToWarehouse(request)
    }
}

func BenchmarkApplyWarehousePatch(b *testing.B) {
    newAddress := "New Address"
    newTelephone := "+2222222222"
    patch := warehouse.WarehousePatchDTO{
        Address:   &newAddress,
        Telephone: &newTelephone,
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        existing := &warehouse.Warehouse{
            Id:                 1,
            WarehouseCode:      "BENCH-CODE",
            Address:            "Old Address",
            Telephone:          "+1111111111",
            MinimumCapacity:    50,
            MinimumTemperature: -10.0,
            LocalityId:         "LOC-1",
        }
        mappers.ApplyWarehousePatch(existing, patch)
    }
}