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
)

func TestBuyerRepository_FindById(t *testing.T) {
	type arrange struct {
		dbMock func() (sqlmock.Sqlmock, *sql.DB)
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

	testCases := []testCase{
		{
			name: "success - buyer found",
			arrange: arrange{
				id: 1,
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					db, mock, _ := sqlmock.New()
					row := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name"}).
						AddRow(1, "CARD-001", "John", "Doe")
					mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name FROM buyers WHERE id = \\?").
						WithArgs(1).
						WillReturnRows(row)
					return mock, db
				},
			},
			output: output{
				buyer: &models.Buyer{
					Id:           1,
					CardNumberId: "CARD-001",
					FirstName:    "John",
					LastName:     "Doe",
				},
				err: nil,
			},
		},
		{
			name: "error - buyer not found",
			arrange: arrange{
				id: 1,
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					db, mock, _ := sqlmock.New()
					mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name FROM buyers WHERE id = \\?").
						WithArgs(1).
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
				id: 1,
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					db, mock, _ := sqlmock.New()
					mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name FROM buyers WHERE id = \\?").
						WithArgs(1).
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
			mock, db := tc.arrange.dbMock()
			defer db.Close()
			repo := repository.NewBuyerRepository(db)

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
