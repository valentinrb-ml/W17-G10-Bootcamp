package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/purchase_order"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestPurchaseOrderRepository_GetCountByBuyer(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(mock sqlmock.Sqlmock)
		argBuyerID     int
		want           []models.BuyerWithPurchaseCount
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "purchase_orders_count"}).
					AddRow(101, "CARD101", "John", "Doe", 3)
				mock.ExpectQuery("SELECT b.id, b.id_card_number, b.first_name, b.last_name, COUNT\\(po.id\\) as purchase_orders_count FROM buyers b LEFT JOIN purchase_orders po ON b.id = po.buyer_id WHERE b.id = \\? GROUP BY b.id").
					WithArgs(101).
					WillReturnRows(rows)
			},
			argBuyerID: 101,
			want:       []models.BuyerWithPurchaseCount{testhelpers.BuyerWithPurchaseCountDummyMap[101]},
			wantErr:    false,
		},
		{
			name: "error - not found",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "purchase_orders_count"})
				mock.ExpectQuery("SELECT b.id, b.id_card_number, b.first_name, b.last_name, COUNT\\(po.id\\) as purchase_orders_count FROM buyers b LEFT JOIN purchase_orders po ON b.id = po.buyer_id WHERE b.id = \\? GROUP BY b.id").
					WithArgs(999).
					WillReturnRows(rows)
			},
			argBuyerID:     999,
			wantErr:        true,
			expectedErrMsg: "buyer not found",
		},
		{
			name: "error - db failure",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT b.id, b.id_card_number, b.first_name, b.last_name, COUNT\\(po.id\\) as purchase_orders_count FROM buyers b LEFT JOIN purchase_orders po ON b.id = po.buyer_id WHERE b.id = \\? GROUP BY b.id").
					WithArgs(101).
					WillReturnError(errors.New("db error"))
			},
			argBuyerID:     101,
			wantErr:        true,
			expectedErrMsg: "error querying purchase count by buyer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setup(mock)
			repo := repository.NewPurchaseOrderRepository(db)

			got, err := repo.GetCountByBuyer(context.Background(), tt.argBuyerID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.expectedErrMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tt.want), len(got))
				for i, buyer := range tt.want {
					require.Equal(t, buyer.ID, got[i].ID)
					require.Equal(t, buyer.CardNumberID, got[i].CardNumberID)
					require.Equal(t, buyer.FirstName, got[i].FirstName)
					require.Equal(t, buyer.LastName, got[i].LastName)
					require.Equal(t, buyer.PurchaseOrdersCount, got[i].PurchaseOrdersCount)
				}
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
