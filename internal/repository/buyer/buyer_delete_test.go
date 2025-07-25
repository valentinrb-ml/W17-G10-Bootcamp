package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
)

func TestBuyerRepository_Delete(t *testing.T) {
	type arrange struct {
		dbMock func() (sqlmock.Sqlmock, *sql.DB)
	}
	type input struct {
		id  int
		ctx context.Context
	}
	type output struct {
		err error
	}
	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}

	testCases := []testCase{
		{
			name: "success - buyer deleted successfully",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					db, mock, _ := sqlmock.New()
					mock.ExpectExec("DELETE FROM buyers WHERE id = ?").
						WithArgs(1).
						WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
					return mock, db
				},
			},
			input: input{
				id:  1,
				ctx: context.Background(),
			},
			output: output{
				err: nil,
			},
		},
		{
			name: "error - buyer not found",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					db, mock, _ := sqlmock.New()
					mock.ExpectExec("DELETE FROM buyers WHERE id = ?").
						WithArgs(999).
						WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
					return mock, db
				},
			},
			input: input{
				id:  999,
				ctx: context.Background(),
			},
			output: output{
				err: apperrors.NewAppError(apperrors.CodeNotFound, "buyer not found"),
			},
		},
		{
			name: "error - buyer has purchase orders (FK constraint)",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					db, mock, _ := sqlmock.New()
					mysqlErr := &mysql.MySQLError{
						Number:  1451,
						Message: "Cannot delete or update a parent row: a foreign key constraint fails",
					}
					mock.ExpectExec("DELETE FROM buyers WHERE id = ?").
						WithArgs(1).
						WillReturnError(mysqlErr)
					return mock, db
				},
			},
			input: input{
				id:  1,
				ctx: context.Background(),
			},
			output: output{
				err: apperrors.NewAppError(apperrors.CodeConflict, "cannot delete buyer: there are purchase orders associated"),
			},
		},
		{
			name: "error - database connection failed",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					db, mock, _ := sqlmock.New()
					mock.ExpectExec("DELETE FROM buyers WHERE id = ?").
						WithArgs(1).
						WillReturnError(errors.New("connection failed"))
					return mock, db
				},
			},
			input: input{
				id:  1,
				ctx: context.Background(),
			},
			output: output{
				err: apperrors.NewAppError(apperrors.CodeInternal, "an internal server error occurred while deleting the buyer: connection failed"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mock, db := tc.arrange.dbMock()
			defer db.Close()
			repo := repository.NewBuyerRepository(db)

			// Act
			err := repo.Delete(tc.input.ctx, tc.input.id)

			// Assert
			if tc.output.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
