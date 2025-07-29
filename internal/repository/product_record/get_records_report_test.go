package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

var cols = []string{"product_id", "description", "records_count"}

/*
Scenarios:
  - all_products_success  – happy path, no id param
  - product_not_found     – query by id returns empty slice → NOT_FOUND
*/
func TestProductRecordRepository_GetRecordsReport(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		productID int
		mockSetup func(sqlmock.Sqlmock)
		wantLen   int    // expected slice length when no error
		wantApp   string // expected AppError.Code ("" == success)
	}{
		{
			name: "all_products",
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols).AddRow(1, "p1", 3).AddRow(2, "p2", 0)
				m.ExpectQuery("FROM products p.*LEFT JOIN product_records").
					WillReturnRows(rows)
			},
			wantLen: 2,
		},
		{
			name:      "product_not_found",
			productID: 999,
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols) // empty result
				m.ExpectQuery("WHERE p.id").
					WithArgs(999).
					WillReturnRows(rows)
			},
			wantApp: apperrors.CodeNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, mock, cleanup := testhelpers.NewProductRecordRepoMock(t)
			defer cleanup()

			tc.mockSetup(mock)

			reports, err := repo.GetRecordsReport(context.Background(), tc.productID)

			if tc.wantApp != "" {
				testhelpers.RequireAppErr(t, err, tc.wantApp)
				require.Nil(t, reports)
				return
			}

			require.NoError(t, err)
			require.Len(t, reports, tc.wantLen)
		})
	}
}
