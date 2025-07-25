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

func TestBuyerRepository_Update(t *testing.T) {
	type arrange struct {
		dbMock func(id int, b models.Buyer) (sqlmock.Sqlmock, *sql.DB)
	}
	type input struct {
		id    int
		buyer models.Buyer
		ctx   context.Context
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

	// const para IDs de prueba
	const (
		testID     = 1
		notFoundID = 999
	)

	// builder para crear buyer de prueba
	testBuyer := testhelpers.NewBuyerBuilder().Build()

	testCases := []testCase{
		{
			name: "success - buyer updated successfully",
			arrange: arrange{
				dbMock: func(id int, b models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mock.ExpectExec("UPDATE buyers SET id_card_number = \\?, first_name = \\?, last_name = \\? WHERE id = \\?").
						WithArgs(b.CardNumberId, b.FirstName, b.LastName, id).
						WillReturnResult(sqlmock.NewResult(0, 1))
					return mock, db
				},
			},
			input: input{
				id:    testID,
				buyer: testBuyer,
				ctx:   context.Background(),
			},
			output: output{
				err: nil,
			},
		},
		{
			name: "error - id_card_number duplicate",
			arrange: arrange{
				dbMock: func(id int, b models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mysqlErr := &mysql.MySQLError{
						Number:  1062,
						Message: "Duplicate entry 'CARD-001' for key 'id_card_number'",
					}
					mock.ExpectExec("UPDATE buyers SET id_card_number = \\?, first_name = \\?, last_name = \\? WHERE id = \\?").
						WithArgs(b.CardNumberId, b.FirstName, b.LastName, id).
						WillReturnError(mysqlErr)
					return mock, db
				},
			},
			input: input{
				id:    testID,
				buyer: testBuyer,
				ctx:   context.Background(),
			},
			output: output{
				err: apperrors.NewAppError(apperrors.CodeConflict, "Could not update buyer: card number already exists."),
			},
		},
		{
			name: "error - buyer not found",
			arrange: arrange{
				dbMock: func(id int, b models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mock.ExpectExec("UPDATE buyers SET id_card_number = \\?, first_name = \\?, last_name = \\? WHERE id = \\?").
						WithArgs(b.CardNumberId, b.FirstName, b.LastName, id).
						WillReturnResult(sqlmock.NewResult(0, 0))
					return mock, db
				},
			},
			input: input{
				id:    notFoundID,
				buyer: testBuyer,
				ctx:   context.Background(),
			},
			output: output{
				err: nil,
			},
		},
		{
			name: "error - database connection failed",
			arrange: arrange{
				dbMock: func(id int, b models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mock.ExpectExec("UPDATE buyers SET id_card_number = \\?, first_name = \\?, last_name = \\? WHERE id = \\?").
						WithArgs(b.CardNumberId, b.FirstName, b.LastName, id).
						WillReturnError(errors.New("connection failed"))
					return mock, db
				},
			},
			input: input{
				id:    testID,
				buyer: testBuyer,
				ctx:   context.Background(),
			},
			output: output{
				err: apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while updating a buyer: connection failed"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mock, db := tc.arrange.dbMock(tc.input.id, tc.input.buyer)
			defer db.Close()
			repo := repository.NewBuyerRepository(db)

			// Act
			err := repo.Update(tc.input.ctx, tc.input.id, tc.input.buyer)

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
