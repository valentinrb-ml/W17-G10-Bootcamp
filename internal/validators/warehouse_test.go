package validators_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
)

func TestValidateWarehouseCreateRequest(t *testing.T) {
    validTemp := 10.0

    t.Run("request válido", func(t *testing.T) {
        req := warehouse.WarehouseRequest{
            Address:            "Calle 1",
            Telephone:          "+1234567890",
            WarehouseCode:      "WH-001",
            MinimumCapacity:    1,
            MinimumTemperature: &validTemp,
            LocalityId:         "LOC-1",
        }
        err := validators.ValidateWarehouseCreateRequest(req)
        assert.NoError(t, err)
    })

    t.Run("campos obligatorios vacíos", func(t *testing.T) {
        req := warehouse.WarehouseRequest{}
        err := validators.ValidateWarehouseCreateRequest(req)
        assert.Error(t, err)
    })

    t.Run("minimumCapacity <= 0", func(t *testing.T) {
        req := warehouse.WarehouseRequest{
            Address:            "Calle 1",
            Telephone:          "+1234567890",
            WarehouseCode:      "WH-001",
            MinimumCapacity:    0,
            MinimumTemperature: &validTemp,
            LocalityId:         "LOC-1",
        }
        err := validators.ValidateWarehouseCreateRequest(req)
        assert.Error(t, err)
    })

    t.Run("minimumTemperature nil", func(t *testing.T) {
        req := warehouse.WarehouseRequest{
            Address:            "Calle 1",
            Telephone:          "+1234567890",
            WarehouseCode:      "WH-001",
            MinimumCapacity:    10,
            MinimumTemperature: nil,
            LocalityId:         "LOC-1",
        }
        err := validators.ValidateWarehouseCreateRequest(req)
        assert.Error(t, err)
    })

    t.Run("LocalityId vacío", func(t *testing.T) {
        req := warehouse.WarehouseRequest{
            Address:            "Calle 1",
            Telephone:          "+1234567890",
            WarehouseCode:      "WH-001",
            MinimumCapacity:    10,
            MinimumTemperature: &validTemp,
            LocalityId:         "",
        }
        err := validators.ValidateWarehouseCreateRequest(req)
        assert.Error(t, err)
    })

    t.Run("teléfono inválido", func(t *testing.T) {
        req := warehouse.WarehouseRequest{
            Address:            "Calle 1",
            Telephone:          "abc123",
            WarehouseCode:      "WH-001",
            MinimumCapacity:    10,
            MinimumTemperature: &validTemp,
            LocalityId:         "LOC-1",
        }
        err := validators.ValidateWarehouseCreateRequest(req)
        assert.Error(t, err)
    })
}

func TestIsValidPhone(t *testing.T) {
    validPhones := []string{"+1234567890", "1234567890", "+123456789012345"}
    invalidPhones := []string{"12345", "phone", "+12abc345", "+1234567890123456"}

    for _, phone := range validPhones {
        assert.True(t, validators.IsValidPhone(phone), "debe ser válido: %s", phone)
    }
    for _, phone := range invalidPhones {
        assert.False(t, validators.IsValidPhone(phone), "debe ser inválido: %s", phone)
    }
}

func TestValidateMinimumCapacity(t *testing.T) {
    t.Run("capacidad válida", func(t *testing.T) {
        err := validators.ValidateMinimumCapacity(1)
        assert.NoError(t, err)
    })

    t.Run("capacidad cero", func(t *testing.T) {
        err := validators.ValidateMinimumCapacity(0)
        assert.Error(t, err)
    })

    t.Run("capacidad negativa", func(t *testing.T) {
        err := validators.ValidateMinimumCapacity(-10)
        assert.Error(t, err)
    })
}
