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
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestBuyerRepository_Delete(t *testing.T) {
	type arrange struct {
		dbMock func(id int) (sqlmock.Sqlmock, *sql.DB)
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

	// Test ID base para reutilización
	const testID = 1
	const notFoundID = 999

	testCases := []testCase{
		{
			name: "success - buyer deleted successfully",
			arrange: arrange{
				dbMock: func(id int) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mock.ExpectExec("DELETE FROM buyers WHERE id = ?").
						WithArgs(id).
						WillReturnResult(sqlmock.NewResult(0, 1))
					return mock, db
				},
			},
			input: input{
				id:  testID,
				ctx: context.Background(),
			},
			output: output{
				err: nil,
			},
		},
		{
			name: "error - buyer not found",
			arrange: arrange{
				dbMock: func(id int) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mock.ExpectExec("DELETE FROM buyers WHERE id = ?").
						WithArgs(id).
						WillReturnResult(sqlmock.NewResult(0, 0))
					return mock, db
				},
			},
			input: input{
				id:  notFoundID,
				ctx: context.Background(),
			},
			output: output{
				err: apperrors.NewAppError(apperrors.CodeNotFound, "buyer not found"),
			},
		},
		{
			name: "error - buyer has purchase orders (FK constraint)",
			arrange: arrange{
				dbMock: func(id int) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mysqlErr := &mysql.MySQLError{
						Number:  1451,
						Message: "Cannot delete or update a parent row: a foreign key constraint fails",
					}
					mock.ExpectExec("DELETE FROM buyers WHERE id = ?").
						WithArgs(id).
						WillReturnError(mysqlErr)
					return mock, db
				},
			},
			input: input{
				id:  testID,
				ctx: context.Background(),
			},
			output: output{
				err: apperrors.NewAppError(apperrors.CodeConflict, "cannot delete buyer: there are purchase orders associated"),
			},
		},
		{
			name: "error - database connection failed",
			arrange: arrange{
				dbMock: func(id int) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mock.ExpectExec("DELETE FROM buyers WHERE id = ?").
						WithArgs(id).
						WillReturnError(errors.New("connection failed"))
					return mock, db
				},
			},
			input: input{
				id:  testID,
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
			mock, db := tc.arrange.dbMock(tc.input.id)
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
