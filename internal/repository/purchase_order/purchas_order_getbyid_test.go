package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/purchase_order"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestPurchaseOrderRepository_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(mock sqlmock.Sqlmock)
		argID          int
		want           *models.PurchaseOrder
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				row := sqlmock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"}).
					AddRow(1, "PO-001", "2023-01-15 00:00:00", "TRACK001", 101, 201)
				mock.ExpectQuery("SELECT id, order_number, order_date, tracking_code, buyer_id, product_record_id FROM purchase_orders WHERE id = \\?").
					WithArgs(1).
					WillReturnRows(row)
			},
			argID: 1,
			// ANTES:
			// want:    &testhelpers.PurchaseOrderDummyMap[1],
			// AHORA:
			want: func() *models.PurchaseOrder {
				v := testhelpers.PurchaseOrderDummyMap[1]
				return &v
			}(),
			wantErr: false,
		},
		{
			name: "error - not found",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, order_number, order_date, tracking_code, buyer_id, product_record_id FROM purchase_orders WHERE id = \\?").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			argID:          999,
			wantErr:        true,
			expectedErrMsg: "purchase order not found",
		},
		{
			name: "error - db failure",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, order_number, order_date, tracking_code, buyer_id, product_record_id FROM purchase_orders WHERE id = \\?").
					WithArgs(1).
					WillReturnError(errors.New("db error"))
			},
			argID:          1,
			wantErr:        true,
			expectedErrMsg: "error querying purchase order by id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setup(mock)
			repo := repository.NewPurchaseOrderRepository(db)

			got, err := repo.GetByID(context.Background(), tt.argID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.expectedErrMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tt.want.ID, got.ID)
				require.Equal(t, tt.want.OrderNumber, got.OrderNumber)
				require.Equal(t, tt.want.TrackingCode, got.TrackingCode)
				require.Equal(t, tt.want.BuyerID, got.BuyerID)
				require.Equal(t, tt.want.ProductRecordID, got.ProductRecordID)
				require.WithinDuration(t, tt.want.OrderDate, got.OrderDate, time.Second)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
