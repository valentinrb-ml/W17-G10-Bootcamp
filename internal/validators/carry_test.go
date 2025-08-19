package validators_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

func TestValidateCarryCreateRequest(t *testing.T) {
    t.Run("request válido", func(t *testing.T) {
        req := carry.CarryRequest{
            Address:     "Calle 1",
            Telephone:   "+1234567890",
            Cid:         "C123",
            CompanyName: "Empresa S.A.",
            LocalityId:  "LOC-1",
        }
        err := validators.ValidateCarryCreateRequest(req)
        assert.NoError(t, err)
    })

    t.Run("campos obligatorios vacíos", func(t *testing.T) {
        req := carry.CarryRequest{}
        err := validators.ValidateCarryCreateRequest(req)
        assert.Error(t, err)
    })

    t.Run("Address vacío", func(t *testing.T) {
        req := carry.CarryRequest{
            Address:     "",
            Telephone:   "+1234567890",
            Cid:         "C123",
            CompanyName: "Empresa S.A.",
            LocalityId:  "LOC-1",
        }
        err := validators.ValidateCarryCreateRequest(req)
        assert.Error(t, err)
    })

    t.Run("Telephone vacío", func(t *testing.T) {
        req := carry.CarryRequest{
            Address:     "Calle 1",
            Telephone:   "",
            Cid:         "C123",
            CompanyName: "Empresa S.A.",
            LocalityId:  "LOC-1",
        }
        err := validators.ValidateCarryCreateRequest(req)
        assert.Error(t, err)
    })

    t.Run("Cid vacío", func(t *testing.T) {
        req := carry.CarryRequest{
            Address:     "Calle 1",
            Telephone:   "+1234567890",
            Cid:         "",
            CompanyName: "Empresa S.A.",
            LocalityId:  "LOC-1",
        }
        err := validators.ValidateCarryCreateRequest(req)
        assert.Error(t, err)
    })

    t.Run("CompanyName vacío", func(t *testing.T) {
        req := carry.CarryRequest{
            Address:     "Calle 1",
            Telephone:   "+1234567890",
            Cid:         "C123",
            CompanyName: "",
            LocalityId:  "LOC-1",
        }
        err := validators.ValidateCarryCreateRequest(req)
        assert.Error(t, err)
    })

    t.Run("LocalityId vacío", func(t *testing.T) {
        req := carry.CarryRequest{
            Address:     "Calle 1",
            Telephone:   "+1234567890",
            Cid:         "C123",
            CompanyName: "Empresa S.A.",
            LocalityId:  "",
        }
        err := validators.ValidateCarryCreateRequest(req)
        assert.Error(t, err)
    })

    t.Run("teléfono inválido", func(t *testing.T) {
        req := carry.CarryRequest{
            Address:     "Calle 1",
            Telephone:   "abc123",
            Cid:         "C123",
            CompanyName: "Empresa S.A.",
            LocalityId:  "LOC-1",
        }
        err := validators.ValidateCarryCreateRequest(req)
        assert.Error(t, err)
    })

    t.Run("teléfono válido", func(t *testing.T) {
        req := carry.CarryRequest{
            Address:     "Calle 1",
            Telephone:   "+1234567890",
            Cid:         "C123",
            CompanyName: "Empresa S.A.",
            LocalityId:  "LOC-1",
        }
        err := validators.ValidateCarryCreateRequest(req)
        assert.NoError(t, err)
    })
}