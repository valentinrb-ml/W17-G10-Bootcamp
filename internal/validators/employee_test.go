package validators_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

func TestValidateEmployee(t *testing.T) {
	testCases := map[string]struct {
		input   *models.Employee
		wantErr string // mensaje de error esperado, "" si no quiero error
	}{
		"ok": {
			input: &models.Employee{
				CardNumberID: "123456",
				FirstName:    "John",
				LastName:     "Doe",
				WarehouseID:  1,
			},
			wantErr: "",
		},
		"nil_employee": {
			input:   nil,
			wantErr: "employee cannot be nil",
		},
		"empty_card_number_id": {
			input: &models.Employee{
				CardNumberID: "",
				FirstName:    "John",
				LastName:     "Doe",
				WarehouseID:  1,
			},
			wantErr: "card_number_id cannot be empty",
		},
		"empty_first_name": {
			input: &models.Employee{
				CardNumberID: "123",
				FirstName:    "",
				LastName:     "Doe",
				WarehouseID:  1,
			},
			wantErr: "first_name cannot be empty",
		},
		"empty_last_name": {
			input: &models.Employee{
				CardNumberID: "123",
				FirstName:    "F",
				LastName:     "",
				WarehouseID:  1,
			},
			wantErr: "last_name cannot be empty",
		},
		"zero_warehouse_id": {
			input: &models.Employee{
				CardNumberID: "123",
				FirstName:    "F",
				LastName:     "L",
				WarehouseID:  0,
			},
			wantErr: "warehouse_id is required",
		},
		"negative_warehouse_id": {
			input: &models.Employee{
				CardNumberID: "123",
				FirstName:    "F",
				LastName:     "L",
				WarehouseID:  -5,
			},
			wantErr: "warehouse_id must be greater than 0",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validators.ValidateEmployee(tc.input)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

func TestValidateEmployeeID(t *testing.T) {
	testCases := map[string]struct {
		id      int
		wantErr string
	}{
		"ok":   {id: 1, wantErr: ""},
		"zero": {id: 0, wantErr: "id must be positive"},
		"neg":  {id: -2, wantErr: "id must be positive"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validators.ValidateEmployeeID(tc.id)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

func TestValidateEmployeePatch(t *testing.T) {
	cn := "999"
	fn := "N"
	ln := "X"
	wid := 2
	zero := 0
	neg := -4
	testCases := map[string]struct {
		patch   *models.EmployeePatch
		wantErr string
	}{
		"ok_card_number": {
			patch: &models.EmployeePatch{
				CardNumberID: &cn,
			},
			wantErr: "",
		},
		"ok_first_name": {
			patch: &models.EmployeePatch{
				FirstName: &fn,
			},
			wantErr: "",
		},
		"ok_last_name": {
			patch: &models.EmployeePatch{
				LastName: &ln,
			},
			wantErr: "",
		},
		"ok_warehouse": {
			patch: &models.EmployeePatch{
				WarehouseID: &wid,
			},
			wantErr: "",
		},
		"empty_patch": {
			patch:   &models.EmployeePatch{},
			wantErr: "at least one field must be provided for update",
		},
		"empty_card_number": {
			patch:   &models.EmployeePatch{CardNumberID: func() *string { s := ""; return &s }()},
			wantErr: "card_number_id cannot be empty",
		},
		"empty_first_name": {
			patch:   &models.EmployeePatch{FirstName: func() *string { s := ""; return &s }()},
			wantErr: "first_name cannot be empty",
		},
		"empty_last_name": {
			patch:   &models.EmployeePatch{LastName: func() *string { s := ""; return &s }()},
			wantErr: "last_name cannot be empty",
		},
		"zero_warehouse_id": {
			patch:   &models.EmployeePatch{WarehouseID: &zero},
			wantErr: "warehouse_id is required",
		},
		"neg_warehouse_id": {
			patch:   &models.EmployeePatch{WarehouseID: &neg},
			wantErr: "warehouse_id must be positive",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validators.ValidateEmployeePatch(tc.patch)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}
