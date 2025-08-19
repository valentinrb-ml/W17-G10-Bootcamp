package mappers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

func TestMapEmployeeToEmployeeDoc(t *testing.T) {
	testCases := map[string]struct {
		input    *models.Employee
		expected models.EmployeeDoc
	}{
		"normal_with_warehouse_id": {
			input: &models.Employee{
				ID:           1,
				CardNumberID: "C1",
				FirstName:    "A",
				LastName:     "B",
				WarehouseID:  7,
			},
			expected: models.EmployeeDoc{
				ID:           1,
				CardNumberID: "C1",
				FirstName:    "A",
				LastName:     "B",
				WarehouseID:  func() *int { v := 7; return &v }(),
			},
		},
		"warehouse_id_zero_to_nil": {
			input: &models.Employee{
				ID:           10,
				CardNumberID: "C10",
				FirstName:    "John",
				LastName:     "Smith",
				WarehouseID:  0,
			},
			expected: models.EmployeeDoc{
				ID:           10,
				CardNumberID: "C10",
				FirstName:    "John",
				LastName:     "Smith",
				WarehouseID:  nil,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := mappers.MapEmployeeToEmployeeDoc(tc.input)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestMapEmployeeDocToEmployee(t *testing.T) {
	testCases := map[string]struct {
		input    models.EmployeeDoc
		expected *models.Employee
	}{
		"nil_warehouseID_to_zero": {
			input: models.EmployeeDoc{
				ID:           2,
				CardNumberID: "X",
				FirstName:    "Y",
				LastName:     "Z",
				WarehouseID:  nil,
			},
			expected: &models.Employee{
				ID:           2,
				CardNumberID: "X",
				FirstName:    "Y",
				LastName:     "Z",
				WarehouseID:  0,
			},
		},
		"warehouseID_pointer": {
			input: models.EmployeeDoc{
				ID:           3,
				CardNumberID: "CC",
				FirstName:    "FF",
				LastName:     "LL",
				WarehouseID:  func() *int { v := 77; return &v }(),
			},
			expected: &models.Employee{
				ID:           3,
				CardNumberID: "CC",
				FirstName:    "FF",
				LastName:     "LL",
				WarehouseID:  77,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := mappers.MapEmployeeDocToEmployee(tc.input)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestMapEmployeePatchToEmployee(t *testing.T) {
	card := "P12"
	fname := "AF"
	lname := "BL"
	wid := 3

	testCases := map[string]struct {
		orig     *models.Employee
		patch    *models.EmployeePatch
		expected *models.Employee
	}{
		"no_changes": {
			orig: &models.Employee{
				ID:           4,
				CardNumberID: "BASE",
				FirstName:    "PEP",
				LastName:     "QUI",
				WarehouseID:  2,
			},
			patch: &models.EmployeePatch{},
			expected: &models.Employee{
				ID:           4,
				CardNumberID: "BASE",
				FirstName:    "PEP",
				LastName:     "QUI",
				WarehouseID:  2,
			},
		},
		"update_all": {
			orig: &models.Employee{
				ID: 7, CardNumberID: "A", FirstName: "B", LastName: "C", WarehouseID: 1,
			},
			patch: &models.EmployeePatch{
				CardNumberID: &card,
				FirstName:    &fname,
				LastName:     &lname,
				WarehouseID:  &wid,
			},
			expected: &models.Employee{
				ID:           7,
				CardNumberID: "P12",
				FirstName:    "AF",
				LastName:     "BL",
				WarehouseID:  3,
			},
		},
		"update_partial_fields": {
			orig: &models.Employee{
				ID: 5, CardNumberID: "X", FirstName: "Y", LastName: "Z", WarehouseID: 9,
			},
			patch: &models.EmployeePatch{
				LastName: &lname,
			},
			expected: &models.Employee{
				ID:           5,
				CardNumberID: "X",
				FirstName:    "Y",
				LastName:     "BL",
				WarehouseID:  9,
			},
		},
		"patch_warehouse_id_zero_noop": {
			orig: &models.Employee{
				ID:           6,
				CardNumberID: "XW",
				FirstName:    "FW",
				LastName:     "LW",
				WarehouseID:  11,
			},
			patch: &models.EmployeePatch{
				WarehouseID: func() *int { v := 0; return &v }(),
			},
			expected: &models.Employee{
				ID:           6,
				CardNumberID: "XW",
				FirstName:    "FW",
				LastName:     "LW",
				WarehouseID:  11, // no cambio
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := mappers.MapEmployeePatchToEmployee(tc.orig, tc.patch)
			require.Equal(t, tc.expected, got)
		})
	}
}
