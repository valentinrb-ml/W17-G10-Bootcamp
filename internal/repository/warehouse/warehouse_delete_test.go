package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestWarehouseMySQL_Delete(t *testing.T) {
	type arrange struct {
		dbMock func() (sqlmock.Sqlmock, *sql.DB)
	}
	type input struct {
		id      int
		context context.Context
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

	// test cases
	testCases := []testCase{
		{
			name: "success - warehouse deleted",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mock.ExpectExec("DELETE FROM warehouse WHERE id = ?").
						WithArgs(1).
						WillReturnResult(sqlmock.NewResult(0, 1))

					return mock, db
				},
			},
			input: input{
				id:      1,
				context: context.Background(),
			},
			output: output{
				err: nil,
			},
		},
		{
			name: "error - warehouse not found",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mock.ExpectExec("DELETE FROM warehouse WHERE id = ?").
						WithArgs(99).
						WillReturnResult(sqlmock.NewResult(0, 0))

					return mock, db
				},
			},
			input: input{
				id:      99,
				context: context.Background(),
			},
			output: output{
				err: apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found"),
			},
		},
		{
			name: "error - database error",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mock.ExpectExec("DELETE FROM warehouse WHERE id = ?").
						WithArgs(1).
						WillReturnError(sql.ErrConnDone)

					return mock, db
				},
			},
			input: input{
				id:      1,
				context: context.Background(),
			},
			output: output{
				err: apperrors.Wrap(sql.ErrConnDone, "error deleting warehouse"),
			},
		},
		{
			name: "error - rows affected check failed",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					result := sqlmock.NewErrorResult(sql.ErrTxDone)
					mock.ExpectExec("DELETE FROM warehouse WHERE id = ?").
						WithArgs(1).
						WillReturnResult(result)

					return mock, db
				},
			},
			input: input{
				id:      1,
				context: context.Background(),
			},
			output: output{
				err: apperrors.Wrap(sql.ErrTxDone, "error deleting warehouse"),
			},
		},
		{
			name: "error - foreign key constraint violation",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mysqlErr := &mysql.MySQLError{
						Number:  1451,
						Message: "Cannot delete or update a parent row: a foreign key constraint fails",
					}
					mock.ExpectExec("DELETE FROM warehouse WHERE id = ?").
						WithArgs(1).
						WillReturnError(mysqlErr)

					return mock, db
				},
			},
			input: input{
				id:      1,
				context: context.Background(),
			},
			output: output{
				err: apperrors.NewAppError(apperrors.CodeConflict, "cannot delete warehouse: it is being referenced by other records"),
			},
		},
	}

	// run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			mock, db := tc.arrange.dbMock()
			defer db.Close()
			repo := repository.NewWarehouseRepository(db)
			repo.SetLogger(testhelpers.NewTestLogger())

			// act
			err := repo.Delete(tc.input.context, tc.input.id)

			// assert
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
