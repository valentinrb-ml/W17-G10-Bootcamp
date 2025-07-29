package repository

// Tests for GetByID(): success, not-found, generic error.

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	productmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestProductRepository_GetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      int
		setup   func(sqlmock.Sqlmock, *sqlmock.ExpectedPrepare)
		wantErr bool
		appCode string
	}{
		{
			name: "success",
			id:   7,
			setup: func(m sqlmock.Sqlmock, p *sqlmock.ExpectedPrepare) {
				rows := sqlmock.NewRows(allColumns()).
					AddRow(7, "A", "B", 1, 1, 1, 1, 1, 1, 1, 10, 20)
				p.ExpectQuery().WithArgs(7).WillReturnRows(rows)
			},
		},
		{
			name:    "not found",
			id:      99,
			wantErr: true,
			appCode: apperrors.CodeNotFound,
			setup: func(m sqlmock.Sqlmock, p *sqlmock.ExpectedPrepare) {
				p.ExpectQuery().WithArgs(99).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:    "generic db error",
			id:      1,
			wantErr: true,
			appCode: apperrors.CodeInternal,
			setup: func(m sqlmock.Sqlmock, p *sqlmock.ExpectedPrepare) {
				p.ExpectQuery().WithArgs(1).WillReturnError(sql.ErrConnDone)
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
			_, err := repo.GetByID(context.Background(), tc.id)

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
