package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	productmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// Patch covers: no changes, single update, multiple updates, rows=0, and errors
// MySQL 1048 (null) to reach case 1048 of handleDBError.
func TestProductRepository_Patch(t *testing.T) {
	t.Parallel()

	// helpers to build pointers quickly
	str := func(s string) *string { return &s }
	f64 := func(v float64) *float64 { return &v }

	code := "NEWCODE"
	reqUpdate := models.ProductPatchRequest{ProductCode: &code}

	multiReq := models.ProductPatchRequest{
		Description: str("x"),
		Width:       f64(2.5),
		Height:      f64(3.5),
	}

	tests := []struct {
		name    string
		id      int
		req     models.ProductPatchRequest
		setup   func(sqlmock.Sqlmock, *sqlmock.ExpectedPrepare)
		wantErr bool
		appCode string
	}{
		{
			name: "nothing to update",
			id:   1, req: models.ProductPatchRequest{},
			setup: func(m sqlmock.Sqlmock, p *sqlmock.ExpectedPrepare) {
				rows := sqlmock.NewRows(allColumns()).
					AddRow(1, "A", "B", 1, 1, 1, 1, 1, 1, 1, 10, 20)
				p.ExpectQuery().WithArgs(1).WillReturnRows(rows)
			},
		},
		{
			name: "patch product_code",
			id:   9, req: reqUpdate,
			setup: func(m sqlmock.Sqlmock, p *sqlmock.ExpectedPrepare) {
				m.ExpectExec(regexp.QuoteMeta(
					"UPDATE products SET product_code = ? WHERE id = ?")).
					WithArgs(code, 9).
					WillReturnResult(sqlmock.NewResult(0, 1))

				rows := sqlmock.NewRows(allColumns()).
					AddRow(9, code, "B", 1, 1, 1, 1, 1, 1, 1, 10, 20)
				p.ExpectQuery().WithArgs(9).WillReturnRows(rows)
			},
		},
		{
			name: "patch multiple fields",
			id:   12, req: multiReq,
			setup: func(m sqlmock.Sqlmock, p *sqlmock.ExpectedPrepare) {
				m.ExpectExec(regexp.QuoteMeta(
					"UPDATE products SET description = ?, width = ?, height = ? WHERE id = ?")).
					WithArgs("x", 2.5, 3.5, 12).
					WillReturnResult(sqlmock.NewResult(0, 1))

				rows := sqlmock.NewRows(allColumns()).
					AddRow(12, "C", "x", 2.5, 3.5, 1, 1, 1, 1, 1, 10, 20)
				p.ExpectQuery().WithArgs(12).WillReturnRows(rows)
			},
		},
		{
			name:    "not found (rows=0)",
			id:      5,
			req:     reqUpdate,
			wantErr: true,
			appCode: apperrors.CodeNotFound,
			setup: func(m sqlmock.Sqlmock, _ *sqlmock.ExpectedPrepare) {
				m.ExpectExec(regexp.QuoteMeta(
					"UPDATE products SET product_code = ? WHERE id = ?")).
					WithArgs(code, 5).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name:    "null field 1048",
			id:      7,
			req:     reqUpdate,
			wantErr: true,
			appCode: apperrors.CodeBadRequest,
			setup: func(m sqlmock.Sqlmock, _ *sqlmock.ExpectedPrepare) {
				mysqlErr := &mysql.MySQLError{Number: 1048, Message: "null"}
				m.ExpectExec(regexp.QuoteMeta(
					"UPDATE products SET product_code = ? WHERE id = ?")).
					WithArgs(code, 7).
					WillReturnError(mysqlErr)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, mock, cancel := productmock.NewMockDB(t)
			defer cancel()

			prepSel := expectPrepareProductRepository(mock)
			tc.setup(mock, prepSel)

			repo, _ := NewProductRepository(db)
			repo.SetLogger(testhelpers.NewTestLogger())
			_, err := repo.Patch(context.Background(), tc.id, tc.req)

			if tc.wantErr {
				require.Error(t, err)
				testhelpers.RequireAppErr(t, err, tc.appCode)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
