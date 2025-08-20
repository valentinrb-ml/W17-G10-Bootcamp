package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	productmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// Tests for GetAll(): happy-path and DB error.
func TestProductRepository_GetAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		wantErr bool
		appCode string
	}{
		{
			name: "success",
			setup: func(m sqlmock.Sqlmock) {
				expectPrepareProductRepository(m)
				rows := sqlmock.NewRows(allColumns()).
					AddRow(1, "CODE-1", "desc", 1.1, 2.2, 3.3, 10, 5, 3.5, 10, 100, 200)
				m.ExpectQuery(regexp.QuoteMeta(baseSelect + " ORDER BY id")).
					WillReturnRows(rows)
			},
		},
		{
			name:    "db error",
			wantErr: true,
			appCode: apperrors.CodeInternal,
			setup: func(m sqlmock.Sqlmock) {
				expectPrepareProductRepository(m)
				m.ExpectQuery(regexp.QuoteMeta(baseSelect + " ORDER BY id")).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, mock, cancel := productmock.NewMockDB(t)
			defer cancel()

			tc.setup(mock)
			repo, _ := NewProductRepository(db)
			repo.SetLogger(testhelpers.NewTestLogger())

			_, err := repo.GetAll(context.Background())
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
