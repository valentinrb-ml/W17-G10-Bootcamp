package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/employee"
)

func TestEmployeeRepository_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		inputID   int
		mockSetup func(mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name:    "delete_ok",
			inputID: 10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM employees").
					WithArgs(10).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectErr: false,
		},
		{
			name:    "delete_db_error",
			inputID: 11,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM employees").
					WithArgs(11).
					WillReturnError(sql.ErrConnDone)
			},
			expectErr: true,
		},
		{
			name:    "no_rows_affected",
			inputID: 12,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM employees").
					WithArgs(12).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectErr: true,
		},
		{
			name:    "rows_affected_fails",
			inputID: 13,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM employees").
					WithArgs(13).
					WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tc.mockSetup(mock)
			repo := repo.NewEmployeeRepository(db)
			ctx := context.Background()
			err = repo.Delete(ctx, tc.inputID)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
