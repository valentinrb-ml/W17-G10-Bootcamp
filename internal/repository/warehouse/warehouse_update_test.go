package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func TestWarehouseMySQL_Update(t *testing.T) {
	type arrange struct {
		dbMock func() (sqlmock.Sqlmock, *sql.DB)
	}
	type input struct {
		id        int
		warehouse warehouse.Warehouse
		context   context.Context
	}
	type output struct {
		warehouse *warehouse.Warehouse
		err       error
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
			name: "success - warehouse updated",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mock.ExpectExec("UPDATE warehouse SET (.+) WHERE id = ?").
						WithArgs("WH001", "123 Main St", 10.5, 1000, "555-1234", "LOC001", 1).
						WillReturnResult(sqlmock.NewResult(0, 1))

					return mock, db
				},
			},
			input: input{
				id:        1,
				warehouse: testhelpers.CreateTestWarehouse(),
				context:   context.Background(),
			},
			output: output{
				warehouse: testhelpers.CreateExpectedWarehouse(1),
				err:       nil,
			},
		},
		{
			name: "error - warehouse not found",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mock.ExpectExec("UPDATE warehouse SET (.+) WHERE id = ?").
						WithArgs("WH001", "123 Main St", 10.5, 1000, "555-1234", "LOC001", 999).
						WillReturnResult(sqlmock.NewResult(0, 0))

					return mock, db
				},
			},
			input: input{
				id:        999,
				warehouse: testhelpers.CreateTestWarehouse(),
				context:   context.Background(),
			},
			output: output{
				warehouse: nil,
				err:       apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found"),
			},
		},
		{
			name: "error - warehouse_code already exists",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mysqlErr := &mysql.MySQLError{
						Number:  1062,
						Message: "Duplicate entry 'WH001' for key 'warehouse_code'",
					}

					mock.ExpectExec("UPDATE warehouse SET (.+) WHERE id = ?").
						WithArgs("WH001", "123 Main St", 10.5, 1000, "555-1234", "LOC001", 1).
						WillReturnError(mysqlErr)

					return mock, db
				},
			},
			input: input{
				id:        1,
				warehouse: testhelpers.CreateTestWarehouse(),
				context:   context.Background(),
			},
			output: output{
				warehouse: nil,
				err:       apperrors.NewAppError(apperrors.CodeConflict, "warehouse_code already exists"),
			},
		},
		{
			name: "error - database error",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mock.ExpectExec("UPDATE warehouse SET (.+) WHERE id = ?").
						WithArgs("WH001", "123 Main St", 10.5, 1000, "555-1234", "LOC001", 1).
						WillReturnError(sql.ErrConnDone)

					return mock, db
				},
			},
			input: input{
				id:        1,
				warehouse: testhelpers.CreateTestWarehouse(),
				context:   context.Background(),
			},
			output: output{
				warehouse: nil,
				err:       apperrors.Wrap(sql.ErrConnDone, "error updating warehouse"),
			},
		},
		{
			name: "error - rows affected check failed",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					result := sqlmock.NewErrorResult(sql.ErrTxDone)
					mock.ExpectExec("UPDATE warehouse SET (.+) WHERE id = ?").
						WithArgs("WH001", "123 Main St", 10.5, 1000, "555-1234", "LOC001", 1).
						WillReturnResult(result)

					return mock, db
				},
			},
			input: input{
				id:        1,
				warehouse: testhelpers.CreateTestWarehouse(),
				context:   context.Background(),
			},
			output: output{
				warehouse: nil,
				err:       apperrors.Wrap(sql.ErrTxDone, "error updating warehouse"),
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

			// act
			result, err := repo.Update(tc.input.context, tc.input.id, tc.input.warehouse)

			// assert
			if tc.output.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tc.output.warehouse.Id, result.Id)
				require.Equal(t, tc.output.warehouse.WarehouseCode, result.WarehouseCode)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
