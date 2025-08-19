package mappers_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

func TestRequestToCarry(t *testing.T) {
    t.Run("convierte correctamente CarryRequest a Carry", func(t *testing.T) {
        req := carry.CarryRequest{
            Cid:         "C123",
            CompanyName: "Transporte S.A.",
            Address:     "Calle Falsa 123",
            Telephone:   "123456789",
            LocalityId:  "LOC-1",
        }
        expected := carry.Carry{
            Id:          0,
            Cid:         "C123",
            CompanyName: "Transporte S.A.",
            Address:     "Calle Falsa 123",
            Telephone:   "123456789",
            LocalityId:  "LOC-1",
        }
        result := mappers.RequestToCarry(req)
        assert.Equal(t, expected, result)
    })

    t.Run("valores vacíos", func(t *testing.T) {
        req := carry.CarryRequest{}
        expected := carry.Carry{
            Id: 0,
        }
        result := mappers.RequestToCarry(req)
        assert.Equal(t, expected, result)
    })
}

func TestCarryToDoc(t *testing.T) {
    t.Run("convierte correctamente Carry a CarryDoc", func(t *testing.T) {
        c := &carry.Carry{
            Id:          1,
            Cid:         "C123",
            CompanyName: "Transporte S.A.",
            Address:     "Calle Falsa 123",
            Telephone:   "123456789",
            LocalityId:  "LOC-1",
        }
        expected := carry.CarryDoc{
            ID:          1,
            Cid:         "C123",
            CompanyName: "Transporte S.A.",
            Address:     "Calle Falsa 123",
            Telephone:   "123456789",
            LocalityId:  "LOC-1",
        }
        result := mappers.CarryToDoc(c)
        assert.Equal(t, expected, result)
    })

    t.Run("valores vacíos", func(t *testing.T) {
        c := &carry.Carry{}
        expected := carry.CarryDoc{}
        result := mappers.CarryToDoc(c)
        assert.Equal(t, expected, result)
    })

    t.Run("panic si el puntero es nil", func(t *testing.T) {
        assert.Panics(t, func() {
            mappers.CarryToDoc(nil)
        })
    })
}

func TestCarryToDocSlice(t *testing.T) {
    t.Run("convierte slice correctamente", func(t *testing.T) {
        carries := []carry.Carry{
            {
                Id:          1,
                Cid:         "C1",
                CompanyName: "Empresa 1",
                Address:     "Dirección 1",
                Telephone:   "111",
                LocalityId:  "LOC-1",
            },
            {
                Id:          2,
                Cid:         "C2",
                CompanyName: "Empresa 2",
                Address:     "Dirección 2",
                Telephone:   "222",
                LocalityId:  "LOC-2",
            },
        }
        result := mappers.CarryToDocSlice(carries)
        assert.Len(t, result, 2)
        assert.Equal(t, carries[0].Id, result[0].ID)
        assert.Equal(t, carries[1].Cid, result[1].Cid)
    })

    t.Run("slice vacío", func(t *testing.T) {
        carries := []carry.Carry{}
        result := mappers.CarryToDocSlice(carries)
        assert.Empty(t, result)
    })

    t.Run("slice nil", func(t *testing.T) {
        var carries []carry.Carry
        result := mappers.CarryToDocSlice(carries)
        assert.Empty(t, result)
    })

    t.Run("slice con un solo elemento", func(t *testing.T) {
        carries := []carry.Carry{
            {
                Id:          99,
                Cid:         "C99",
                CompanyName: "Empresa 99",
                Address:     "Dirección 99",
                Telephone:   "999",
                LocalityId:  "LOC-99",
            },
        }
        result := mappers.CarryToDocSlice(carries)
        assert.Len(t, result, 1)
        assert.Equal(t, carries[0].Id, result[0].ID)
        assert.Equal(t, carries[0].LocalityId, result[0].LocalityId)
    })
}