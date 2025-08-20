package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestBuyerRepository_FindById(t *testing.T) {
	type arrange struct {
		dbMock func(id int) (sqlmock.Sqlmock, *sql.DB)
		id     int
	}
	type output struct {
		buyer *models.Buyer
		err   error
	}
	type testCase struct {
		name    string
		arrange arrange
		output  output
	}

	const testID = 1

	testCases := []testCase{
		{
			name: "success - buyer found",
			arrange: arrange{
				id: testID,
				dbMock: func(id int) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					expectedBuyer := testhelpers.CreateTestBuyerWithID(id)
					rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name"}).
						AddRow(expectedBuyer.Id, expectedBuyer.CardNumberId, expectedBuyer.FirstName, expectedBuyer.LastName)
					mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name FROM buyers WHERE id = \\?").
						WithArgs(id).
						WillReturnRows(rows)
					return mock, db
				},
			},
			output: output{
				buyer: testhelpers.CreateTestBuyerWithID(testID),
				err:   nil,
			},
		},
		{
			name: "error - buyer not found",
			arrange: arrange{
				id: testID,
				dbMock: func(id int) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name FROM buyers WHERE id = \\?").
						WithArgs(id).
						WillReturnError(sql.ErrNoRows)
					return mock, db
				},
			},
			output: output{
				buyer: nil,
				err:   apperrors.NewAppError(apperrors.CodeNotFound, "The buyer you are looking for does not exist."),
			},
		},
		{
			name: "error - database error",
			arrange: arrange{
				id: testID,
				dbMock: func(id int) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name FROM buyers WHERE id = \\?").
						WithArgs(id).
						WillReturnError(errors.New("database error"))
					return mock, db
				},
			},
			output: output{
				buyer: nil,
				err:   apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the buyer."),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mock, db := tc.arrange.dbMock(tc.arrange.id)
			defer db.Close()
			repo := repository.NewBuyerRepository(db)
			repo.SetLogger(testhelpers.NewTestLogger())

			// Act
			buyer, err := repo.FindById(context.Background(), tc.arrange.id)

			// Assert
			if tc.output.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.output.buyer, buyer)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
