package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/product_batch"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestProductBatchesRepository_GetReportProduct(t *testing.T) {
	type arrange struct {
		dbMock func(sqlmock.Sqlmock)
	}
	type output struct {
		expected      []models.ReportProduct
		expectedError bool
		err           error
	}
	type testCase struct {
		name    string
		arrange arrange
		output  output
	}

	expected := testhelpers.DummyReportProductsList()

	const query = `SELECT s.id, s.section_number, SUM\(p.current_quantity\) FROM product_batches p INNER JOIN sections s on p.section_id = s.id\s+GROUP BY p.section_id`

	testCases := []testCase{
		{
			name: "success - returns products report",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{
						"id", "section_number", "sum",
					}).
						AddRow(expected[0].SectionId, expected[0].SectionNumber, expected[0].ProductsCount).
						AddRow(expected[1].SectionId, expected[1].SectionNumber, expected[1].ProductsCount)
					m.ExpectQuery(query).WillReturnRows(rows)
				},
			},
			output: output{
				expected:      expected,
				expectedError: false,
				err:           nil,
			},
		},
		{
			name: "success - empty report",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id", "section_number", "sum"})
					m.ExpectQuery(query).WillReturnRows(rows)
				},
			},
			output: output{
				expected:      []models.ReportProduct{},
				expectedError: false,
				err:           nil,
			},
		},
		{
			name: "error in query (db error)",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectQuery(query).WillReturnError(errors.New("internal error"))
				},
			},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the products report."),
			},
		},
		{
			name: "error in scan (malformed row)",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					// Provoke scan error: put wrong type in one column
					rows := sqlmock.NewRows([]string{"id", "section_number", "sum"}).
						AddRow("bad", 7, 42)
					m.ExpectQuery(query).WillReturnRows(rows)
				},
			},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the product report."),
			},
		},
		{
			name: "error in rows.Err() after loop",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id", "section_number", "sum"}).
						AddRow(expected[0].SectionId, expected[0].SectionNumber, expected[0].ProductsCount)
					rows.RowError(0, sql.ErrConnDone)
					m.ExpectQuery(query).WillReturnRows(rows)
				},
			},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the product report."),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repository := repo.NewProductBatchesRepository(db)
			repository.SetLogger(testhelpers.NewTestLogger())

			tc.arrange.dbMock(mock)

			result, err := repository.GetReportProduct(context.Background())
			if tc.output.expectedError {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.output.expected, result)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
