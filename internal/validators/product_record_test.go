package validators

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

func TestValidateProductRecordCreateRequest(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		req         models.ProductRecordRequest
		wantErr     bool
		wantMessage string
	}{
		{
			name: "valid request => ok",
			req: models.ProductRecordRequest{
				Data: models.ProductRecordCore{
					LastUpdateDate: now,
					ProductID:      1,
					PurchasePrice:  10,
					SalePrice:      15,
				},
			},
			wantErr: false,
		},
		{
			name: "missing last_update_date => validation error",
			req: models.ProductRecordRequest{
				Data: models.ProductRecordCore{
					LastUpdateDate: time.Time{}, // zero
					ProductID:      1,
				},
			},
			wantErr:     true,
			wantMessage: "last_update_date is required",
		},
		{
			name: "non-positive product_id => validation error",
			req: models.ProductRecordRequest{
				Data: models.ProductRecordCore{
					LastUpdateDate: now,
					ProductID:      0,
				},
			},
			wantErr:     true,
			wantMessage: "product_id must be greater than 0",
		},
		{
			name: "both invalid => returns first validation error (last_update_date)",
			req: models.ProductRecordRequest{
				Data: models.ProductRecordCore{
					LastUpdateDate: time.Time{},
					ProductID:      0,
				},
			},
			wantErr:     true,
			wantMessage: "last_update_date is required",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateProductRecordCreateRequest(tt.req)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if tt.wantErr {
				var appErr *apperrors.AppError
				if !errors.As(err, &appErr) {
					t.Fatalf("expected *apperrors.AppError, got %T", err)
				}
				if tt.wantMessage != "" && !strings.Contains(err.Error(), tt.wantMessage) {
					t.Fatalf("expected error message to contain %q, got %q", tt.wantMessage, err.Error())
				}
			}
		})
	}
}

func TestValidateProductRecordBusinessRules(t *testing.T) {
	t.Parallel()

	now := time.Now()
	future := now.Add(2 * time.Hour)

	tests := []struct {
		name        string
		record      models.ProductRecord
		wantErr     bool
		wantMessage string
	}{
		{
			name: "valid record => ok",
			record: models.ProductRecord{
				ID: 1,
				ProductRecordCore: models.ProductRecordCore{
					LastUpdateDate: now,
					PurchasePrice:  0,
					SalePrice:      0,
					ProductID:      10,
				},
			},
			wantErr: false,
		},
		{
			name: "negative purchase price => bad request error",
			record: models.ProductRecord{
				ProductRecordCore: models.ProductRecordCore{
					LastUpdateDate: now,
					PurchasePrice:  -0.01,
					SalePrice:      1.0,
					ProductID:      10,
				},
			},
			wantErr:     true,
			wantMessage: "purchase_price must be greater than or equal to 0",
		},
		{
			name: "negative sale price => bad request error",
			record: models.ProductRecord{
				ProductRecordCore: models.ProductRecordCore{
					LastUpdateDate: now,
					PurchasePrice:  1.0,
					SalePrice:      -0.01,
					ProductID:      10,
				},
			},
			wantErr:     true,
			wantMessage: "sale_price must be greater than or equal to 0",
		},
		{
			name: "last_update_date in the future => bad request error",
			record: models.ProductRecord{
				ProductRecordCore: models.ProductRecordCore{
					LastUpdateDate: future,
					PurchasePrice:  1.0,
					SalePrice:      2.0,
					ProductID:      10,
				},
			},
			wantErr:     true,
			wantMessage: "last_update_date cannot be in the future",
		},
		{
			name: "multiple invalid => returns the first check (purchase_price)",
			record: models.ProductRecord{
				ProductRecordCore: models.ProductRecordCore{
					LastUpdateDate: future,
					PurchasePrice:  -5,
					SalePrice:      -7,
					ProductID:      10,
				},
			},
			wantErr:     true,
			wantMessage: "purchase_price must be greater than or equal to 0",
		},
		{
			name: "last_update_date equal to now => ok",
			record: models.ProductRecord{
				ProductRecordCore: models.ProductRecordCore{
					LastUpdateDate: now, // not after now
					PurchasePrice:  5,
					SalePrice:      7,
					ProductID:      10,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateProductRecordBusinessRules(tt.record)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if tt.wantErr {
				var appErr *apperrors.AppError
				if !errors.As(err, &appErr) {
					t.Fatalf("expected *apperrors.AppError, got %T", err)
				}
				if tt.wantMessage != "" && !strings.Contains(err.Error(), tt.wantMessage) {
					t.Fatalf("expected error message to contain %q, got %q", tt.wantMessage, err.Error())
				}
			}
		})
	}
}
