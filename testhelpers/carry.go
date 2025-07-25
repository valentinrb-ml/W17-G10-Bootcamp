package testhelpers

import (
	"fmt"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

// CreateTestCarry creates a test carry for use in tests
func CreateTestCarry(id int) *carry.Carry {
    return &carry.Carry{
        Id:          id,
        Cid:         fmt.Sprintf("CAR%03d", id),
        CompanyName: fmt.Sprintf("Test Company %d", id),
        Address:     fmt.Sprintf("Test Address %d", id),
        Telephone:   "5551234567",
        LocalityId:  "1",
    }
}

// CreateExpectedCarry creates expected carry for assertions
func CreateExpectedCarry(id int) *carry.Carry {
    return CreateTestCarry(id)
}

// CreateTestCarriesReport creates a test carries report
func CreateTestCarriesReport(localityID, localityName string, count int) *carry.CarriesReport {
    return &carry.CarriesReport{
        LocalityID:   localityID,
        LocalityName: localityName,
        CarriesCount: count,
    }
}