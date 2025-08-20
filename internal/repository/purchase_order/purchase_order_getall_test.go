package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/purchase_order"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestPurchaseOrderRepository_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(mock sqlmock.Sqlmock)
		want           []models.PurchaseOrder
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success - multiple orders",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"}).
					AddRow(1, "PO-001", "2023-01-15 00:00:00", "TRACK001", 101, 201).
					AddRow(2, "PO-002", "2023-02-20 00:00:00", "TRACK002", 102, 202)
				mock.ExpectQuery("SELECT id, order_number, order_date, tracking_code, buyer_id, product_record_id FROM purchase_orders").
					WillReturnRows(rows)
			},
			want: []models.PurchaseOrder{
				testhelpers.PurchaseOrderDummyMap[1],
				testhelpers.PurchaseOrderDummyMap[2],
			},
			wantErr: false,
		},
		{
			name: "success - no orders",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"})
				mock.ExpectQuery("SELECT id, order_number, order_date, tracking_code, buyer_id, product_record_id FROM purchase_orders").
					WillReturnRows(rows)
			},
			want:    []models.PurchaseOrder{},
			wantErr: false,
		},
		{
			name: "error - db failure",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, order_number, order_date, tracking_code, buyer_id, product_record_id FROM purchase_orders").
					WillReturnError(errors.New("db error"))
			},
			wantErr:        true,
			expectedErrMsg: "error querying all purchase orders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setup(mock)
			repo := repository.NewPurchaseOrderRepository(db)
			repo.SetLogger(testhelpers.NewTestLogger())

			got, err := repo.GetAll(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				if tt.expectedErrMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tt.want), len(got))
				for i, po := range tt.want {
					require.Equal(t, po.ID, got[i].ID)
					require.Equal(t, po.OrderNumber, got[i].OrderNumber)
					require.Equal(t, po.TrackingCode, got[i].TrackingCode)
					require.Equal(t, po.BuyerID, got[i].BuyerID)
					require.Equal(t, po.ProductRecordID, got[i].ProductRecordID)
					require.WithinDuration(t, po.OrderDate, got[i].OrderDate, time.Second)
				}
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
