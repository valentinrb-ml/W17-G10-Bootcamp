package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/carry"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestCarryRepository_Create(t *testing.T) {
	type arrange struct {
		dbMock func() (sqlmock.Sqlmock, *sql.DB)
	}
	type input struct {
		carry   carry.Carry
		context context.Context
	}
	type output struct {
		carry *carry.Carry
		err   error
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
			name: "success - carry created successfully",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()
					mock.ExpectExec(`INSERT INTO carriers \(cid, company_name, address, telephone, locality_id\) VALUES \(\?, \?, \?, \?, \?\)`).
						WithArgs("CAR000", "Test Company 0", "Test Address 0", "5551234567", "1").
						WillReturnResult(sqlmock.NewResult(1, 1))
					return mock, db
				},
			},
			input: input{
				carry:   *testhelpers.CreateTestCarry(0),
				context: context.Background(),
			},
			output: output{
				carry: func() *carry.Carry {
					expected := testhelpers.CreateTestCarry(0)
					expected.Id = 1
					return expected
				}(),
				err: nil,
			},
		},
		{
			name: "error - duplicate CID",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()
					mysqlErr := &mysql.MySQLError{
						Number:  1062,
						Message: "Duplicate entry 'CAR000' for key 'cid'",
					}
					mock.ExpectExec(`INSERT INTO carriers \(cid, company_name, address, telephone, locality_id\) VALUES \(\?, \?, \?, \?, \?\)`).
						WithArgs("CAR000", "Test Company 0", "Test Address 0", "5551234567", "1").
						WillReturnError(mysqlErr)
					return mock, db
				},
			},
			input: input{
				carry:   *testhelpers.CreateTestCarry(0),
				context: context.Background(),
			},
			output: output{
				carry: nil,
				err:   apperrors.NewAppError(apperrors.CodeConflict, "cid already exists"),
			},
		},
		{
			name: "error - invalid locality_id",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()
					mysqlErr := &mysql.MySQLError{
						Number:  1452,
						Message: "Cannot add or update a child row: a foreign key constraint fails",
					}
					mock.ExpectExec(`INSERT INTO carriers \(cid, company_name, address, telephone, locality_id\) VALUES \(\?, \?, \?, \?, \?\)`).
						WithArgs("CAR000", "Test Company 0", "Test Address 0", "5551234567", "1").
						WillReturnError(mysqlErr)
					return mock, db
				},
			},
			input: input{
				carry:   *testhelpers.CreateTestCarry(0),
				context: context.Background(),
			},
			output: output{
				carry: nil,
				err:   apperrors.NewAppError(apperrors.CodeConflict, "locality_id does not exist"),
			},
		},
		{
			name: "error - LastInsertId fails",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()
					mock.ExpectExec(`INSERT INTO carriers \(cid, company_name, address, telephone, locality_id\) VALUES \(\?, \?, \?, \?, \?\)`).
						WithArgs("CAR000", "Test Company 0", "Test Address 0", "5551234567", "1").
						WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))
					return mock, db
				},
			},
			input: input{
				carry:   *testhelpers.CreateTestCarry(0),
				context: context.Background(),
			},
			output: output{
				carry: nil,
				err:   apperrors.Wrap(sql.ErrNoRows, "error creating carry"),
			},
		},
		{
			name: "error - database connection error",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()
					mock.ExpectExec(`INSERT INTO carriers \(cid, company_name, address, telephone, locality_id\) VALUES \(\?, \?, \?, \?, \?\)`).
						WithArgs("CAR000", "Test Company 0", "Test Address 0", "5551234567", "1").
						WillReturnError(sql.ErrConnDone)
					return mock, db
				},
			},
			input: input{
				carry:   *testhelpers.CreateTestCarry(0),
				context: context.Background(),
			},
			output: output{
				carry: nil,
				err:   apperrors.Wrap(sql.ErrConnDone, "error creating carry"),
			},
		},
	}

	// run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			mock, db := tc.arrange.dbMock()
			defer db.Close()
			repo := repository.NewCarryRepository(db)
			repo.SetLogger(testhelpers.NewTestLogger())

			// act
			result, err := repo.Create(tc.input.context, tc.input.carry)

			// assert
			if tc.output.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tc.output.carry.Id, result.Id)
				require.Equal(t, tc.output.carry.Cid, result.Cid)
				require.Equal(t, tc.output.carry.CompanyName, result.CompanyName)
				require.Equal(t, tc.output.carry.Address, result.Address)
				require.Equal(t, tc.output.carry.Telephone, result.Telephone)
				require.Equal(t, tc.output.carry.LocalityId, result.LocalityId)
			}

			// verify all expectations were met
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
