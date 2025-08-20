package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/carry"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestCarryRepository_GetCarriesCountByAllLocalities(t *testing.T) {
	type arrange struct {
		dbMock func() (sqlmock.Sqlmock, *sql.DB)
	}
	type input struct {
		context context.Context
	}
	type output struct {
		reports []carry.CarriesReport
		err     error
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
			name: "success - multiple localities with carries",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					rows := sqlmock.NewRows([]string{"locality_id", "name", "carries_count"}).
						AddRow("1", "Buenos Aires", 5).
						AddRow("2", "Córdoba", 3).
						AddRow("3", "Rosario", 2)

					mock.ExpectQuery(`SELECT c\.locality_id, l\.name, COUNT\(\*\) as carries_count FROM carriers c INNER JOIN localities l ON c\.locality_id = l\.id GROUP BY c\.locality_id`).
						WillReturnRows(rows)

					return mock, db
				},
			},
			input: input{
				context: context.Background(),
			},
			output: output{
				reports: []carry.CarriesReport{
					*testhelpers.CreateTestCarriesReport("1", "Buenos Aires", 5),
					*testhelpers.CreateTestCarriesReport("2", "Córdoba", 3),
					*testhelpers.CreateTestCarriesReport("3", "Rosario", 2),
				},
				err: nil,
			},
		},
		{
			name: "success - empty result",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					rows := sqlmock.NewRows([]string{"locality_id", "name", "carries_count"})

					mock.ExpectQuery(`SELECT c\.locality_id, l\.name, COUNT\(\*\) as carries_count FROM carriers c INNER JOIN localities l ON c\.locality_id = l\.id GROUP BY c\.locality_id`).
						WillReturnRows(rows)

					return mock, db
				},
			},
			input: input{
				context: context.Background(),
			},
			output: output{
				reports: []carry.CarriesReport{},
				err:     nil,
			},
		},
		{
			name: "error - database query fails",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mock.ExpectQuery(`SELECT c\.locality_id, l\.name, COUNT\(\*\) as carries_count FROM carriers c INNER JOIN localities l ON c\.locality_id = l\.id GROUP BY c\.locality_id`).
						WillReturnError(sql.ErrConnDone)

					return mock, db
				},
			},
			input: input{
				context: context.Background(),
			},
			output: output{
				reports: nil,
				err:     sql.ErrConnDone,
			},
		},
		{
			name: "error - rows iteration error",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					rows := sqlmock.NewRows([]string{"locality_id", "name", "carries_count"}).
						AddRow("1", "Buenos Aires", 5).
						CloseError(sql.ErrConnDone)

					mock.ExpectQuery(`SELECT c\.locality_id, l\.name, COUNT\(\*\) as carries_count FROM carriers c INNER JOIN localities l ON c\.locality_id = l\.id GROUP BY c\.locality_id`).
						WillReturnRows(rows)

					return mock, db
				},
			},
			input: input{
				context: context.Background(),
			},
			output: output{
				reports: nil,
				err:     apperrors.Wrap(sql.ErrConnDone, "error getting carries count by all localities"),
			},
		},
		{
			name: "partial success - scan error on some rows",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					// Simulate rows where some have scan errors (wrong data types)
					rows := sqlmock.NewRows([]string{"locality_id", "name", "carries_count"}).
						AddRow("1", "Buenos Aires", 5).
						AddRow("invalid", "Córdoba", "invalid_count").
						AddRow("3", "Rosario", 2)

					mock.ExpectQuery(`SELECT c\.locality_id, l\.name, COUNT\(\*\) as carries_count FROM carriers c INNER JOIN localities l ON c\.locality_id = l\.id GROUP BY c\.locality_id`).
						WillReturnRows(rows)

					return mock, db
				},
			},
			input: input{
				context: context.Background(),
			},
			output: output{
				reports: []carry.CarriesReport{
					*testhelpers.CreateTestCarriesReport("1", "Buenos Aires", 5),
					*testhelpers.CreateTestCarriesReport("3", "Rosario", 2),
				},
				err: nil,
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
			result, err := repo.GetCarriesCountByAllLocalities(tc.input.context)

			// assert
			if tc.output.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tc.output.reports), len(result))

				for i, expectedReport := range tc.output.reports {
					require.Equal(t, expectedReport.LocalityID, result[i].LocalityID)
					require.Equal(t, expectedReport.LocalityName, result[i].LocalityName)
					require.Equal(t, expectedReport.CarriesCount, result[i].CarriesCount)
				}
			}

			// verify all expectations were met
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
