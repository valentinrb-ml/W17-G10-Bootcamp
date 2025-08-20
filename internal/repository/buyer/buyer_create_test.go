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
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestBuyerRepository_Create(t *testing.T) {
	type arrange struct {
		dbMock func(b models.Buyer) (sqlmock.Sqlmock, *sql.DB)
	}
	type input struct {
		buyer models.Buyer
		ctx   context.Context
	}
	type output struct {
		buyer *models.Buyer
		err   error
	}
	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}

	// Test cases
	testCases := []testCase{
		{
			name: "success - buyer created successfully",
			arrange: arrange{
				dbMock: func(b models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mock.ExpectExec("INSERT INTO buyers").
						WithArgs(b.CardNumberId, b.FirstName, b.LastName).
						WillReturnResult(sqlmock.NewResult(1, 1))
					return mock, db
				},
			},
			input: input{
				buyer: testhelpers.NewBuyerBuilder().Build(),
				ctx:   context.Background(),
			},
			output: output{
				buyer: testhelpers.CreateTestBuyerWithID(1),
				err:   nil,
			},
		},
		{
			name: "error - id_card_number duplicate",
			arrange: arrange{
				dbMock: func(b models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mysqlErr := &mysql.MySQLError{
						Number:  1062,
						Message: "Duplicate entry 'CARD-001' for key 'id_card_number'",
					}
					mock.ExpectExec("INSERT INTO buyers").
						WithArgs(b.CardNumberId, b.FirstName, b.LastName).
						WillReturnError(mysqlErr)
					return mock, db
				},
			},
			input: input{
				buyer: testhelpers.NewBuyerBuilder().Build(),
				ctx:   context.Background(),
			},
			output: output{
				buyer: nil,
				err:   apperrors.NewAppError(apperrors.CodeConflict, "Could not create buyer: card number already exists."),
			},
		},
		{
			name: "error - database connection failed",
			arrange: arrange{
				dbMock: func(b models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mock.ExpectExec("INSERT INTO buyers").
						WithArgs(b.CardNumberId, b.FirstName, b.LastName).
						WillReturnError(sql.ErrConnDone)
					return mock, db
				},
			},
			input: input{
				buyer: testhelpers.NewBuyerBuilder().Build(),
				ctx:   context.Background(),
			},
			output: output{
				buyer: nil,
				err:   apperrors.Wrap(sql.ErrConnDone, "An internal server error occurred while creating a buyer: sql: connection is already closed"),
			},
		},
		{
			name: "error - failed to get last insert id",
			arrange: arrange{
				dbMock: func(b models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					result := sqlmock.NewErrorResult(errors.New("error getting ID"))
					mock.ExpectExec("INSERT INTO buyers").
						WithArgs(b.CardNumberId, b.FirstName, b.LastName).
						WillReturnResult(result)
					return mock, db
				},
			},
			input: input{
				buyer: testhelpers.NewBuyerBuilder().Build(),
				ctx:   context.Background(),
			},
			output: output{
				buyer: nil,
				err:   apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while creating a buyer: error getting ID"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mock, db := tc.arrange.dbMock(tc.input.buyer)
			defer db.Close()
			repo := repository.NewBuyerRepository(db)
			repo.SetLogger(testhelpers.NewTestLogger())

			// Act
			result, err := repo.Create(tc.input.ctx, tc.input.buyer)

			// Assert
			if tc.output.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.output.buyer, result)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
