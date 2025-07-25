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

func TestBuyerRepository_FindAll(t *testing.T) {
	type arrange struct {
		dbMock func(buyers []models.Buyer) (sqlmock.Sqlmock, *sql.DB)
	}
	type output struct {
		buyers []models.Buyer
		err    error
	}
	type testCase struct {
		name    string
		arrange arrange
		output  output
	}

	// Usamos el helper para crear buyers de prueba
	testBuyers := testhelpers.CreateTestBuyers(2)

	testCases := []testCase{
		{
			name: "success - found multiple buyers",
			arrange: arrange{
				dbMock: func(buyers []models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					rows := testhelpers.CreateBuyerRows(mock, buyers)
					mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name FROM buyers").
						WillReturnRows(rows)
					return mock, db
				},
			},
			output: output{
				buyers: testBuyers,
				err:    nil,
			},
		},
		{
			name: "success - no buyers found",
			arrange: arrange{
				dbMock: func(_ []models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name"})
					mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name FROM buyers").
						WillReturnRows(rows)
					return mock, db
				},
			},
			output: output{
				buyers: []models.Buyer{},
				err:    nil,
			},
		},
		{
			name: "error - database query failed",
			arrange: arrange{
				dbMock: func(_ []models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name FROM buyers").
						WillReturnError(errors.New("connection failed"))
					return mock, db
				},
			},
			output: output{
				buyers: nil,
				err:    apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while finding all buyers: connection failed"),
			},
		},
		{
			name: "error - row scanning failed",
			arrange: arrange{
				dbMock: func(_ []models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name"}).
						AddRow("invalid_id", "CARD-001", "John", "Doe") // Valor incorrecto para id

					mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name FROM buyers").
						WillReturnRows(rows)
					return mock, db
				},
			},
			output: output{
				buyers: nil,
				err:    apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while finding all buyers: sql: Scan error on column index 0, name \"id\": converting driver.Value type string (\"invalid_id\") to a int: invalid syntax"),
			},
		},
		{
			name: "error - rows iteration failed",
			arrange: arrange{
				dbMock: func(buyers []models.Buyer) (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockBuyerDB()
					rows := testhelpers.CreateBuyerRows(mock, buyers).
						CloseError(errors.New("rows error"))

					mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name FROM buyers").
						WillReturnRows(rows)
					return mock, db
				},
			},
			output: output{
				buyers: nil,
				err:    apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while finding all buyers: rows error"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mock, db := tc.arrange.dbMock(testBuyers)
			defer db.Close()
			repo := repository.NewBuyerRepository(db)

			// Act
			buyers, err := repo.FindAll(context.Background())

			// Assert
			if tc.output.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.output.err.Error())
			} else {
				require.NoError(t, err)
			}

			if tc.output.buyers == nil {
				require.Nil(t, buyers)
			} else {
				require.Equal(t, tc.output.buyers, buyers)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
