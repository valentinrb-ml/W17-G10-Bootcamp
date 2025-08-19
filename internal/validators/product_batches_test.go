package validators_test

import (
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
	"testing"
	"time"
)

func intPtr(v int) *int             { return &v }
func float64Ptr(v float64) *float64 { return &v }

func validPostProductBatch() models.PostProductBatches {
	return models.PostProductBatches{
		BatchNumber:        123,
		CurrentQuantity:    intPtr(10),
		CurrentTemperature: float64Ptr(5.5),
		DueDate:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		InitialQuantity:    intPtr(20),
		ManufacturingDate:  time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		ManufacturingHour:  12,
		MinimumTemperature: float64Ptr(1.2),
		ProductId:          1,
		SectionId:          2,
	}
}

func TestValidateProductBatchPost(t *testing.T) {
	testCases := []struct {
		name    string
		input   models.PostProductBatches
		wantMsg string
	}{
		{
			name:    "valid product batch returns nil",
			input:   validPostProductBatch(),
			wantMsg: "",
		},
		{
			name: "fails if a required int field is 0",
			input: func() models.PostProductBatches {
				pb := validPostProductBatch()
				pb.BatchNumber = 0
				return pb
			}(),
			wantMsg: "All fields are required",
		},
		{
			name: "fails if a required pointer field is nil",
			input: func() models.PostProductBatches {
				pb := validPostProductBatch()
				pb.CurrentQuantity = nil
				return pb
			}(),
			wantMsg: "All fields are required",
		},
		{
			name: "fails if date is zero (due date)",
			input: func() models.PostProductBatches {
				pb := validPostProductBatch()
				pb.DueDate = time.Time{}
				return pb
			}(),
			wantMsg: "All fields are required",
		},
		{
			name: "fails if current quantity is negative",
			input: func() models.PostProductBatches {
				pb := validPostProductBatch()
				pb.CurrentQuantity = intPtr(-1)
				return pb
			}(),
			wantMsg: "Quantity values cannot be negative.",
		},
		{
			name: "fails if initial quantity is negative",
			input: func() models.PostProductBatches {
				pb := validPostProductBatch()
				pb.InitialQuantity = intPtr(-2)
				return pb
			}(),
			wantMsg: "Quantity values cannot be negative.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ValidateProductBatchPost(tc.input)
			if tc.wantMsg == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				appErr, ok := err.(*apperrors.AppError)
				require.True(t, ok)
				require.Contains(t, appErr.Message, tc.wantMsg)
			}
		})
	}
}
