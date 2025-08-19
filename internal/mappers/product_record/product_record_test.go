package mappers

import (
	"testing"
	"time"

	productrecord "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

func TestProductRecordRequestToDomain(t *testing.T) {
	t.Parallel()

	// Use fixed times (no monotonic component) to avoid equality issues
	t1 := time.Date(2024, 7, 10, 15, 30, 0, 0, time.UTC)
	tZero := time.Time{}

	tests := []struct {
		name string
		in   productrecord.ProductRecordRequest
		want productrecord.ProductRecord
	}{
		{
			name: "maps all fields from request to domain and sets ID=0",
			in: productrecord.ProductRecordRequest{
				Data: productrecord.ProductRecordCore{
					LastUpdateDate: t1,
					PurchasePrice:  10.5,
					SalePrice:      12.3,
					ProductID:      42,
				},
			},
			want: productrecord.ProductRecord{
				ID: 0,
				ProductRecordCore: productrecord.ProductRecordCore{
					LastUpdateDate: t1,
					PurchasePrice:  10.5,
					SalePrice:      12.3,
					ProductID:      42,
				},
			},
		},
		{
			name: "maps zero values and keeps ID=0",
			in: productrecord.ProductRecordRequest{
				Data: productrecord.ProductRecordCore{
					LastUpdateDate: tZero,
					PurchasePrice:  0,
					SalePrice:      0,
					ProductID:      0,
				},
			},
			want: productrecord.ProductRecord{
				ID: 0,
				ProductRecordCore: productrecord.ProductRecordCore{
					LastUpdateDate: tZero,
					PurchasePrice:  0,
					SalePrice:      0,
					ProductID:      0,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ProductRecordRequestToDomain(tt.in)

			// ID must always be 0 per mapper contract
			if got.ID != 0 {
				t.Fatalf("ID mismatch: got %d, want 0", got.ID)
			}

			assertEqualProductRecordCore(t, got.ProductRecordCore, tt.want.ProductRecordCore)
		})
	}
}

// Helpers

func assertEqualProductRecordCore(t *testing.T, got, want productrecord.ProductRecordCore) {
	t.Helper()

	// Compare time with Equal to avoid monotonic differences
	if (got.LastUpdateDate.IsZero() && !want.LastUpdateDate.IsZero()) ||
		(!got.LastUpdateDate.IsZero() && want.LastUpdateDate.IsZero()) ||
		(!got.LastUpdateDate.IsZero() && !want.LastUpdateDate.IsZero() && !got.LastUpdateDate.Equal(want.LastUpdateDate)) {
		t.Fatalf("LastUpdateDate mismatch: got %v, want %v", got.LastUpdateDate, want.LastUpdateDate)
	}

	if got.PurchasePrice != want.PurchasePrice ||
		got.SalePrice != want.SalePrice ||
		got.ProductID != want.ProductID {
		t.Fatalf("ProductRecordCore mismatch.\n got: %+v\nwant: %+v", got, want)
	}
}
