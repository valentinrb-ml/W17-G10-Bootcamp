package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestWarehouseMySQL_Create(t *testing.T) {
	type arrange struct {
		dbMock func() (sqlmock.Sqlmock, *sql.DB)
	}
	type input struct {
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
			name: "success - warehouse created successfully",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()
					mock.ExpectExec("INSERT INTO warehouse").
						WithArgs("WH001", "123 Main St", 10.5, 1000, "5551234567", "LOC001").
						WillReturnResult(sqlmock.NewResult(1, 1))
					return mock, db
				},
			},
			input: input{
				warehouse: testhelpers.CreateTestWarehouse(),
				context:   context.Background(),
			},
			output: output{
				warehouse: testhelpers.CreateExpectedWarehouse(1),
				err:       nil,
			},
		},
		{
			name: "error - warehouse_code already exists",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()
					//Simulo error de duplicidad de clave
					mysqlErr := &mysql.MySQLError{
						Number:  1062,
						Message: "Duplicate entry 'WH001' for key 'warehouse_code'",
					}

					mock.ExpectExec("INSERT INTO warehouse").
						WithArgs("WH001", "123 Main St", 10.5, 1000, "5551234567", "LOC001").
						WillReturnError(mysqlErr)
					return mock, db
				},
			},
			input: input{
				warehouse: testhelpers.CreateTestWarehouse(),
				context:   context.Background(),
			},
			output: output{
				warehouse: nil,
				err:       apperrors.NewAppError(apperrors.CodeConflict, "warehouse_code already exists"),
			},
		},
		{
			name: "error - database generic error",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mock.ExpectExec("INSERT INTO warehouse").
						WithArgs("WH001", "123 Main St", 10.5, 1000, "5551234567", "LOC001").
						WillReturnError(sql.ErrConnDone)
					return mock, db
				},
			},
			input: input{
				warehouse: testhelpers.CreateTestWarehouse(),
				context:   context.Background(),
			},
			output: output{
				warehouse: nil,
				err:       apperrors.Wrap(sql.ErrConnDone, "error creating warehouse"),
			},
		},
		{
			name: "error - failed to get last insert id",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					//Simulo error al obtener ID
					result := sqlmock.NewErrorResult(sql.ErrNoRows)
					mock.ExpectExec("INSERT INTO warehouse").
						WithArgs("WH001", "123 Main St", 10.5, 1000, "5551234567", "LOC001").
						WillReturnResult(result)

					return mock, db
				},
			},
			input: input{
				warehouse: testhelpers.CreateTestWarehouse(),
				context:   context.Background(),
			},
			output: output{
				warehouse: nil,
				err:       apperrors.Wrap(sql.ErrNoRows, "error creating warehouse"),
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
			result, err := repo.Create(tc.input.context, tc.input.warehouse)

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
				require.Equal(t, tc.output.warehouse.Address, result.Address)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWarehouseRepository_Create_WithNilLogger(t *testing.T) {
	// arrange
	mock, db := testhelpers.CreateMockDB()
	defer db.Close()

	expectedWarehouse := testhelpers.CreateExpectedWarehouse(1)

	mock.ExpectExec("INSERT INTO warehouse").
		WithArgs("WH001", "123 Main St", 10.5, 1000, "5551234567", "LOC001").
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewWarehouseRepository(db)
	// Don't set logger, so it remains nil

	warehouse := testhelpers.CreateTestWarehouse()

	// act
	result, err := repo.Create(context.Background(), warehouse)

	// assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, expectedWarehouse.Id, result.Id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_Create_LastInsertIdError(t *testing.T) {
	// arrange
	mock, db := testhelpers.CreateMockDB()
	defer db.Close()

	// Mock ExecContext to return a result that will fail on LastInsertId
	mock.ExpectExec("INSERT INTO warehouse").
		WithArgs("WH001", "123 Main St", 10.5, 1000, "5551234567", "LOC001").
		WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

	repo := repository.NewWarehouseRepository(db)
	repo.SetLogger(&SimpleTestLogger{})

	warehouse := testhelpers.CreateTestWarehouse()

	// act
	result, err := repo.Create(context.Background(), warehouse)

	// assert
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "error creating warehouse")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_Create_Success_WithLogger(t *testing.T) {
	// arrange - success case with logger
	mock, db := testhelpers.CreateMockDB()
	defer db.Close()

	expectedWarehouse := testhelpers.CreateExpectedWarehouse(1)

	mock.ExpectExec("INSERT INTO warehouse").
		WithArgs("WH001", "123 Main St", 10.5, 1000, "5551234567", "LOC001").
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewWarehouseRepository(db)
	repo.SetLogger(&SimpleTestLogger{})

	warehouse := testhelpers.CreateTestWarehouse()

	// act
	result, err := repo.Create(context.Background(), warehouse)

	// assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, expectedWarehouse.Id, result.Id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_Create_DuplicateKey_WithLogger(t *testing.T) {
	// arrange - duplicate key with logger
	mock, db := testhelpers.CreateMockDB()
	defer db.Close()

	mysqlErr := &mysql.MySQLError{
		Number:  1062,
		Message: "Duplicate entry 'WH001' for key 'warehouse_code'",
	}

	mock.ExpectExec("INSERT INTO warehouse").
		WithArgs("WH001", "123 Main St", 10.5, 1000, "5551234567", "LOC001").
		WillReturnError(mysqlErr)

	repo := repository.NewWarehouseRepository(db)
	repo.SetLogger(&SimpleTestLogger{})

	warehouse := testhelpers.CreateTestWarehouse()

	// act
	result, err := repo.Create(context.Background(), warehouse)

	// assert
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "warehouse_code already exists")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_Create_GenericError_WithLogger(t *testing.T) {
	// arrange - generic error with logger
	mock, db := testhelpers.CreateMockDB()
	defer db.Close()

	mock.ExpectExec("INSERT INTO warehouse").
		WithArgs("WH001", "123 Main St", 10.5, 1000, "5551234567", "LOC001").
		WillReturnError(sql.ErrConnDone)

	repo := repository.NewWarehouseRepository(db)
	repo.SetLogger(&SimpleTestLogger{})

	warehouse := testhelpers.CreateTestWarehouse()

	// act
	result, err := repo.Create(context.Background(), warehouse)

	// assert
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "error creating warehouse")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_Create_LastInsertIdError_WithLogger(t *testing.T) {
	// arrange - last insert id error with logger
	mock, db := testhelpers.CreateMockDB()
	defer db.Close()

	mock.ExpectExec("INSERT INTO warehouse").
		WithArgs("WH001", "123 Main St", 10.5, 1000, "5551234567", "LOC001").
		WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

	repo := repository.NewWarehouseRepository(db)
	repo.SetLogger(&SimpleTestLogger{})

	warehouse := testhelpers.CreateTestWarehouse()

	// act
	result, err := repo.Create(context.Background(), warehouse)

	// assert
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "error creating warehouse")
	require.NoError(t, mock.ExpectationsWereMet())
}
