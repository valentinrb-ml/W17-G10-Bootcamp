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

func TestCarryRepository_GetCarriesCountByLocalityID(t *testing.T) {
	type arrange struct {
		dbMock func() (sqlmock.Sqlmock, *sql.DB)
	}
	type input struct {
		localityID string
		context    context.Context
	}
	type output struct {
		report *carry.CarriesReport
		err    error
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
			name: "success - locality found with carries",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					rows := sqlmock.NewRows([]string{"locality_id", "name", "carries_count"}).
						AddRow("1", "Buenos Aires", 5)

					mock.ExpectQuery(`SELECT c\.locality_id, l\.name, COUNT\(\*\) as carries_count FROM carriers c INNER JOIN localities l ON c\.locality_id = l\.id WHERE c\.locality_id = \? GROUP BY c\.locality_id`).
						WithArgs("1").
						WillReturnRows(rows)

					return mock, db
				},
			},
			input: input{
				localityID: "1",
				context:    context.Background(),
			},
			output: output{
				report: testhelpers.CreateTestCarriesReport("1", "Buenos Aires", 5),
				err:    nil,
			},
		},
		{
			name: "error - locality not found",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mock.ExpectQuery(`SELECT c\.locality_id, l\.name, COUNT\(\*\) as carries_count FROM carriers c INNER JOIN localities l ON c\.locality_id = l\.id WHERE c\.locality_id = \? GROUP BY c\.locality_id`).
						WithArgs("999").
						WillReturnError(sql.ErrNoRows)

					return mock, db
				},
			},
			input: input{
				localityID: "999",
				context:    context.Background(),
			},
			output: output{
				report: nil,
				err:    apperrors.Wrap(sql.ErrNoRows, "error getting carries count by locality id"),
			},
		},
		{
			name: "error - database connection error",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					mock.ExpectQuery(`SELECT c\.locality_id, l\.name, COUNT\(\*\) as carries_count FROM carriers c INNER JOIN localities l ON c\.locality_id = l\.id WHERE c\.locality_id = \? GROUP BY c\.locality_id`).
						WithArgs("1").
						WillReturnError(sql.ErrConnDone)

					return mock, db
				},
			},
			input: input{
				localityID: "1",
				context:    context.Background(),
			},
			output: output{
				report: nil,
				err:    apperrors.Wrap(sql.ErrConnDone, "error getting carries count by locality id"),
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
			result, err := repo.GetCarriesCountByLocalityID(tc.input.context, tc.input.localityID)

			// assert
			if tc.output.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tc.output.report.LocalityID, result.LocalityID)
				require.Equal(t, tc.output.report.LocalityName, result.LocalityName)
				require.Equal(t, tc.output.report.CarriesCount, result.CarriesCount)
			}

			// verify all expectations were met
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
