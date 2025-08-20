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

func TestPurchaseOrderRepository_Create(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(mock sqlmock.Sqlmock)
		arg            models.PurchaseOrder
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				// Mock para verificar existencia de buyer
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM buyers WHERE id = \\?\\)").
					WithArgs(101).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				// Mock para verificar existencia de product record
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM product_records WHERE id = \\?\\)").
					WithArgs(201).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				// Mock para verificar si existe order number
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM purchase_orders WHERE order_number = \\?\\)").
					WithArgs("PO-001").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

				// Mock para la creación
				mock.ExpectExec("INSERT INTO purchase_orders").
					WithArgs("PO-001", sqlmock.AnyArg(), "TRACK001", 101, 201).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			arg:     testhelpers.PurchaseOrderDummyMap[1],
			wantErr: false,
		},
		{
			name: "error - buyer does not exist",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM buyers WHERE id = \\?\\)").
					WithArgs(100).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			arg: models.PurchaseOrder{
				OrderNumber:     "PO-999",
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK999",
				BuyerID:         100,
				ProductRecordID: 200,
			},
			wantErr:        true,
			expectedErrMsg: "buyer with id 100 does not exist",
		},
		{
			name: "error - product record does not exist",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM buyers WHERE id = \\?\\)").
					WithArgs(101).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM product_records WHERE id = \\?\\)").
					WithArgs(999).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			arg: models.PurchaseOrder{
				OrderNumber:     "PO-999",
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK999",
				BuyerID:         101,
				ProductRecordID: 999,
			},
			wantErr:        true,
			expectedErrMsg: "product record with id 999 does not exist",
		},
		{
			name: "error - order number already exists",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM buyers WHERE id = \\?\\)").
					WithArgs(101).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM product_records WHERE id = \\?\\)").
					WithArgs(201).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM purchase_orders WHERE order_number = \\?\\)").
					WithArgs("DUPLICATE").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			arg: models.PurchaseOrder{
				OrderNumber:     "DUPLICATE",
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK999",
				BuyerID:         101,
				ProductRecordID: 201,
			},
			wantErr:        true,
			expectedErrMsg: "order_number already exists",
		},
		{
			name: "error - db failure on create",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM buyers WHERE id = \\?\\)").
					WithArgs(101).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM product_records WHERE id = \\?\\)").
					WithArgs(201).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM purchase_orders WHERE order_number = \\?\\)").
					WithArgs("PO-001").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

				mock.ExpectExec("INSERT INTO purchase_orders").
					WithArgs("PO-001", sqlmock.AnyArg(), "TRACK001", 101, 201).
					WillReturnError(errors.New("db error"))
			},
			arg:            testhelpers.PurchaseOrderDummyMap[1],
			wantErr:        true,
			expectedErrMsg: "error creating purchase order",
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

			got, err := repo.Create(context.Background(), tt.arg)

			if tt.wantErr {
				require.Error(t, err)
				if tt.expectedErrMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tt.arg.OrderNumber, got.OrderNumber)
				require.Equal(t, tt.arg.BuyerID, got.BuyerID)
				require.Equal(t, tt.arg.ProductRecordID, got.ProductRecordID)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
