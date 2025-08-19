package validators_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

func TestValidateInboundOrder(t *testing.T) {
	testCases := map[string]struct {
		input   *models.InboundOrder
		wantErr string
	}{
		"ok": {
			input: &models.InboundOrder{
				OrderDate:      "2024-06-10",
				OrderNumber:    "INB001",
				EmployeeID:     1,
				ProductBatchID: 5,
				WarehouseID:    42,
			},
			wantErr: "",
		},
		"nil_inbound_order": {
			input:   nil,
			wantErr: "inbound order cannot be nil",
		},
		"empty_order_date": {
			input: &models.InboundOrder{
				OrderDate:      "",
				OrderNumber:    "X",
				EmployeeID:     2,
				ProductBatchID: 3,
				WarehouseID:    4,
			},
			wantErr: "order_date is required",
		},
		"spaces_order_date": {
			input: &models.InboundOrder{
				OrderDate:      "   ",
				OrderNumber:    "X",
				EmployeeID:     2,
				ProductBatchID: 3,
				WarehouseID:    4,
			},
			wantErr: "order_date is required",
		},
		"empty_order_number": {
			input: &models.InboundOrder{
				OrderDate:      "2024-01-01",
				OrderNumber:    "",
				EmployeeID:     2,
				ProductBatchID: 3,
				WarehouseID:    4,
			},
			wantErr: "order_number is required",
		},
		"spaces_order_number": {
			input: &models.InboundOrder{
				OrderDate:      "2024-01-01",
				OrderNumber:    "   ",
				EmployeeID:     2,
				ProductBatchID: 3,
				WarehouseID:    4,
			},
			wantErr: "order_number is required",
		},
		"employee_id_zero": {
			input: &models.InboundOrder{
				OrderDate:      "2024-01-01",
				OrderNumber:    "X",
				EmployeeID:     0,
				ProductBatchID: 1,
				WarehouseID:    2,
			},
			wantErr: "employee_id is required and must be positive",
		},
		"employee_id_negative": {
			input: &models.InboundOrder{
				OrderDate:      "2024-01-01",
				OrderNumber:    "X",
				EmployeeID:     -1,
				ProductBatchID: 1,
				WarehouseID:    2,
			},
			wantErr: "employee_id is required and must be positive",
		},
		"product_batch_id_zero": {
			input: &models.InboundOrder{
				OrderDate:      "2024-01-01",
				OrderNumber:    "X",
				EmployeeID:     1,
				ProductBatchID: 0,
				WarehouseID:    2,
			},
			wantErr: "product_batch_id is required and must be positive",
		},
		"product_batch_id_negative": {
			input: &models.InboundOrder{
				OrderDate:      "2024-01-01",
				OrderNumber:    "X",
				EmployeeID:     1,
				ProductBatchID: -2,
				WarehouseID:    2,
			},
			wantErr: "product_batch_id is required and must be positive",
		},
		"warehouse_id_zero": {
			input: &models.InboundOrder{
				OrderDate:      "2024-01-01",
				OrderNumber:    "X",
				EmployeeID:     1,
				ProductBatchID: 2,
				WarehouseID:    0,
			},
			wantErr: "warehouse_id is required and must be positive",
		},
		"warehouse_id_negative": {
			input: &models.InboundOrder{
				OrderDate:      "2024-01-01",
				OrderNumber:    "X",
				EmployeeID:     1,
				ProductBatchID: 2,
				WarehouseID:    -3,
			},
			wantErr: "warehouse_id is required and must be positive",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validators.ValidateInboundOrder(tc.input)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}
