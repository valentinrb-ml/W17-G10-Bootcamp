package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/purchase_order"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestPurchaseOrderRepository_ExistsOrderNumber(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(mock sqlmock.Sqlmock)
		argNumber  string
		wantExists bool
	}{
		{
			name: "exists",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM purchase_orders WHERE order_number = \\?\\)").
					WithArgs("EXISTS").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			argNumber:  "EXISTS",
			wantExists: true,
		},
		{
			name: "does not exist",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM purchase_orders WHERE order_number = \\?\\)").
					WithArgs("NOT-EXISTS").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			argNumber:  "NOT-EXISTS",
			wantExists: false,
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

			exists := repo.ExistsOrderNumber(context.Background(), tt.argNumber)
			require.Equal(t, tt.wantExists, exists)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
