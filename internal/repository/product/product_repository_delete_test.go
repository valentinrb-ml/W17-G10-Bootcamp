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
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// Tests Delete(): OK, not-found, FK-constraint error.
func TestProductRepository_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rows    int64
		dbErr   error
		wantErr bool
		appCode string
	}{
		{"delete ok", 1, nil, false, ""},
		{"not found", 0, nil, true, apperrors.CodeNotFound},
		{"fk error", 0, &mysql.MySQLError{Number: 1451, Message: "fk"}, true, apperrors.CodeConflict},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, mock, cancel := productmock.NewMockDB(t)
			defer cancel()

			expectPrepareProductRepository(mock)
			exec := mock.ExpectExec(regexp.QuoteMeta(deleteProduct)).WithArgs(10)
			if tc.dbErr != nil {
				exec.WillReturnError(tc.dbErr)
			} else {
				exec.WillReturnResult(sqlmock.NewResult(0, tc.rows))
			}

			repo, _ := NewProductRepository(db)
			err := repo.Delete(context.Background(), 10)

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
