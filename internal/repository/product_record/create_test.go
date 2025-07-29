package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestProductRecordRepository_Create(t *testing.T) {
	t.Parallel()

	type want struct {
		id  int    // expected id on success
		app string // expected AppError.Code on error
	}

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock, models.ProductRecord)
		want      want
	}{
		{
			name: "success",
			mockSetup: func(m sqlmock.Sqlmock, r models.ProductRecord) {
				m.ExpectExec("INSERT INTO product_records").
					WithArgs(sqlmock.AnyArg(), r.PurchasePrice, r.SalePrice, r.ProductID).
					WillReturnResult(sqlmock.NewResult(123, 1))
			},
			want: want{id: 123},
		},
		{
			name: "fk_violation",
			mockSetup: func(m sqlmock.Sqlmock, r models.ProductRecord) {
				m.ExpectExec("INSERT INTO product_records").
					WithArgs(sqlmock.AnyArg(), r.PurchasePrice, r.SalePrice, r.ProductID).
					WillReturnError(&mysql.MySQLError{Number: 1452})
			},
			want: want{app: apperrors.CodeConflict},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, mock, cleanup := testhelpers.NewProductRecordRepoMock(t)
			defer cleanup()

			rec := testhelpers.BuildProductRecord()
			tc.mockSetup(mock, rec)

			got, err := repo.Create(context.Background(), rec)

			if tc.want.app != "" {
				testhelpers.RequireAppErr(t, err, tc.want.app)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want.id, got.ID)
		})
	}
}
